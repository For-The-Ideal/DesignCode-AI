# Flutter 真机预览(腾讯云 TKE Serverless)落地手册

> 目标：生成 Flutter 代码后，用腾讯云**弹性容器(TKE Serverless)**按需跑 `flutter build web`，
> 产物传 **COS 对象存储**，容器**回调主后端**写入 `iframe_url`，前端详情页用 `<iframe>` 显示真机预览。
> 原则：**绝不在现有低配服务器上编译**(会 OOM)；编译全部在独立、按秒计费的弹性容器里完成。

本手册分 6 个阶段，按顺序做。标 🔒 的步骤完成后打勾，全勾即链路跑通。

---

## 0. 架构与数据流(先看懂再动手)

```
用户点生成 → 后端 POST /api/v1/task → 入队
   → worker 跑工作流：Vision → FlutterGenerateSkill(出 Dart 代码) → 存 Code/Preview(HTML 兜底)
   → 【新增】触发构建：后端用 client-go 在 TKE Serverless 建一个 Job
        Job(Pod, 2C4G, 按秒计费)：
           1) 取代码 → 写入 Flutter Web 脚手架 lib/main.dart
           2) flutter pub get && flutter build web --release
           3) coscli 上传 build/web → COS://preview/{task_id}/
           4) curl 回调 主后端 /api/v1/internal/preview-callback {task_id, iframe_url, status}
        Job 结束即销毁
   → 后端收到回调：写 Result.iframe_url + SSE 推 preview_ready
   → 前端详情页监听 preview_ready → 把预览从 HTML 兜底切成 <iframe src=iframe_url>
```

关键约定(定死后各段都按它来)：
- 回调地址：`http://111.229.28.247/api/v1/internal/preview-callback`(走现有 nginx `/api/` 规则，不用改 nginx)。
- 回调鉴权：请求头 `X-Build-Token: <你自己定的一串长随机密钥>`，后端与容器两侧一致即可。
- COS 产物 URL：`https://<bucket>.cos.<region>.myqcloud.com/preview/<task_id>/index.html`，`preview/` 前缀**公有读**。
- 编译超时：8 分钟硬熔断；失败重试 1 次；仍失败 → 前端保留 HTML 预览(降级)。
- 并发上限：默认 2~3 个 Job 同时跑(防费用失控)。

---

## 1. 阶段一：账号与腾讯云 COS 准备(10 分钟)

1. 登录 https://console.cloud.tencent.com/ ，确认可用「对象存储 COS」。
2. 你已有桶：`upload-image-1301638170`(区域 `ap-guangzhou`)。
3. 设置访问权限：COS 控制台 → 该桶 → **权限管理** → 访问权限 → 改为 **公有读私有写**。
   > 若不放心整桶公有读，可只对 `preview/` 前缀配「存储桶权限」里的子账号/公有读策略；最简单先用整桶公有读跑通。
4. 记下这对密钥(已在你后端 `config.dev.yaml`/`config.prod.yaml` 的 `cos` 段，容器也用同一对)：
   - `secret_id` / `secret_key`(COS 控制台 → 访问管理 API 密钥可查/新建，建议用子账号最小权限密钥)。
5. 🔒 COS 公有读已开、密钥可用。

---

## 2. 阶段二：制作 Flutter 构建镜像并推到 TCR(核心，30~60 分钟)

编译在容器里跑，需要一个**装好 Flutter + Chrome + coscli** 的镜像。

### 2.1 开通容器镜像服务 TCR
- 控制台 → **容器镜像服务** → 开通 **个人版(免费)** → 建一个命名空间(如 `designcode`)。
- 设置仓库访问密码(或获取临时凭证)。记下：
  - registry 地址(形如 `ccr.ccs.tencentyun.com`)
  - 用户名(一般是你的 UIN/账号)
  - 密码

### 2.2 在本地(或某台能联网的机器)准备构建目录

新建文件夹 `flutter-builder/`，内含 4 个文件：

**① `Dockerfile`**
```dockerfile
FROM ubuntu:22.04
ENV DEBIAN_FRONTEND=noninteractive

# 基础依赖 + Chrome(flutter web 构建/渲染需要) + git + curl
RUN apt-get update && apt-get install -y --no-install-recommends \
    git curl unzip xz-utils ca-certificates \
    wget gnupg2 \
    && rm -rf /var/lib/apt/lists/*

# 安装 coscli（COS 命令行 v5）
RUN curl -L -o /usr/local/bin/coscli \
    https://github.com/tencentyun/coscli/releases/latest/download/coscli-linux \
    && chmod +x /usr/local/bin/coscli

# 安装 Flutter stable（走国内镜像加速 pub）
ENV FLUTTER_HOME=/opt/flutter
RUN git clone -b stable https://github.com/flutter/flutter.git $FLUTTER_HOME \
    && echo 'export PATH=$PATH:$FLUTTER_HOME/bin:$FLUTTER_HOME/bin/cache/dart-sdk/bin' >> /etc/profile
ENV PATH="$PATH:/opt/flutter/bin:/opt/flutter/bin/cache/dart-sdk/bin"
ENV PUB_HOSTED_URL=https://pub.flutter-io.cn
ENV FLUTTER_STORAGE_BASE_URL=https://storage.flutter-io.cn

# 预热：创建模板工程并预编译一次(把 pub 依赖 & web 引擎缓存进镜像，省每次下载时间)
COPY template/ /opt/template/
RUN cd /opt/template && flutter config --no-analytics && flutter pub get || true

COPY build.sh /usr/local/bin/build.sh
RUN chmod +x /usr/local/bin/build.sh
ENTRYPOINT ["/usr/local/bin/build.sh"]
```

**② `build.sh`**(容器被拉起时执行的脚本 = 整条链路的编译端)
```bash
#!/usr/bin/env bash
set -euo pipefail

# ── 由后端 Job 注入的环境变量 ──
: "${TASK_ID:?}"          # 任务 ID
: "${CODE_B64:?}"         # 生成的 Dart 代码，base64（避免多行/引号问题）
: "${CALLBACK_URL:?}"     # http://111.229.28.247/api/v1/internal/preview-callback
: "${BUILD_TOKEN:?}"      # 共享密钥，回调时带上
: "${COS_BUCKET:?}"       # upload-image-1301638170
: "${COS_REGION:?}"       # ap-guangzhou
: "${COS_SECRET_ID:?}"
: "${COS_SECRET_KEY:?}"

IFRAME_PATH="preview/${TASK_ID}"
INDEX_URL="https://${COS_BUCKET}.cos.${COS_REGION}.myqcloud.com/${IFRAME_PATH}/index.html"

post_callback() { # $1=status $2=url $3=error
  curl -sS -X POST "$CALLBACK_URL" \
    -H "Content-Type: application/json" \
    -H "X-Build-Token: ${BUILD_TOKEN}" \
    -d "{\"task_id\":\"${TASK_ID}\",\"status\":\"$1\",\"iframe_url\":\"$2\",\"error\":\"$3\"}" || true
}

# 出错即回调降级（必须放在 post_callback 定义之后、编译步骤之前才能捕获失败）
trap 'post_callback "failed" "" "build error on line $LINENO"' ERR

cleanup() {
  # coscli 凭证
  cat > /root/.cos.yaml <<EOF
cos:
  base_url: ${COS_REGION}.cos.myqcloud.com
  bucket: ${COS_BUCKET}
  secretid: ${COS_SECRET_ID}
  secretkey: ${COS_SECRET_KEY}
  region: ${COS_REGION}
EOF
}

echo "[build] task=$TASK_ID 解码代码…"
mkdir -p /work && cp -r /opt/template/. /work/
echo "$CODE_B64" | base64 -d > /work/lib/generated_page.dart

echo "[build] flutter pub get…"
cd /work && flutter pub get

echo "[build] flutter build web --release…"
flutter build web --release --web-renderer canvaskit

echo "[build] 上传产物到 COS…"
cleanup
coscli cp -r /work/build/web "${COS_BUCKET}/${IFRAME_PATH}" --only-show-error

echo "[build] 回调成功 url=$INDEX_URL"
post_callback "success" "$INDEX_URL" ""
```

**③ `template/` 目录**：一个最小 Flutter Web 脚手架，把 AI 生成的 widget 当页面加载。
```
template/
├── pubspec.yaml
└── lib/
    ├── main.dart            # 固定入口，引用 generated_page.dart
    └── generated_page.dart  # 构建时被覆盖成 AI 生成的代码
```

`template/pubspec.yaml`
```yaml
name: preview_app
description: DesignCode flutter web preview
publish_to: 'none'
version: 0.0.1
environment:
  sdk: '>=3.0.0 <4.0.0'
dependencies:
  flutter:
    sdk: flutter
flutter:
  uses-material-design: true
```

`template/lib/main.dart`（固定不变，约定生成的代码里必须有一个 `class GeneratedWidget extends StatelessWidget`）
```dart
import 'package:flutter/material.dart';
import 'generated_page.dart';

void main() => runApp(const PreviewApp());

class PreviewApp extends StatelessWidget {
  const PreviewApp({super.key});
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      home: Scaffold(body: SafeArea(child: GeneratedWidget())),
    );
  }
}
```
> ⚠️ 契约：后端 Flutter 提示词要让模型输出**带 `class GeneratedWidget extends StatelessWidget`** 的整份 Dart（顶部含 import）。若现有提示词输出的是别的类名，需要在 `build.sh` 里做类名归一，或调整提示词。这一步和后端改动要对齐。

### 2.3 构建并推送镜像
在能联网、装了 Docker 的机器上（用你本地或临时机器都行）：
```bash
cd flutter-builder
docker build -t ccr.ccs.tencentyun.com/<命名空间>/flutter-builder:v1 .
docker login ccr.ccs.tencentyun.com -u <用户名>
docker push ccr.ccs.tencentyun.com/<命名空间>/flutter-builder:v1
```
🔒 镜像已推到 TCR，能在控制台看到。

> 镜像首次构建较大（含 Flutter SDK 约数 GB），属正常。之后每次任务复用镜像 + 复用镜像里的 pub 缓存，只编译当次代码。

---

## 3. 阶段三：开通 TKE Serverless 集群 + 授权后端拉 Job(20 分钟)

1. 控制台 → **容器服务 TKE** → 创建集群 → 集群类型选 **Serverless 集群**（弹性容器 ECI）→ 地域 `ap-guangzhou`。
2. 创建后进入集群 → **网络**：确认挂了 **公网 NAT 网关**（容器要联网拉 pub/推 COS/回调主后端）。没有就新建一个 NAT。
3. **凭据管理**：获取集群 **kubeconfig**（外网访问），后面给后端用。
4. **最小权限 ServiceAccount**（比整份 kubeconfig 更安全，推荐）：
   在集群里创建一个只允许建 Job 的账号，把它的 token 给后端：
   ```yaml
   # rbac.yaml  → kubectl apply -f rbac.yaml
   apiVersion: v1
   kind: ServiceAccount
   metadata: { name: build-trigger, namespace: default }
   ---
   apiVersion: rbac.authorization.k8s.io/v1
   kind: Role
   metadata: { name: job-creator, namespace: default }
   rules:
     - apiGroups: ["batch"]   resources: ["jobs"]       verbs: ["create","get","list","delete"]
     - apiGroups: [""]        resources: ["pods","pods/log"] verbs: ["get","list"]
   ---
   apiVersion: rbac.authorization.k8s.io/v1
   kind: RoleBinding
   metadata: { name: build-trigger-bind, namespace: default }
   subjects: [{ kind: ServiceAccount, name: build-trigger, namespace: default }]
   roleRef:  { kind: Role, name: job-creator, apiGroup: rbac.authorization.k8s.io }
   ```
5. 🔒 集群可用、kubeconfig/SA token 已拿到、NAT 公网可出。

### 3.1 被拉起的 Job 长这样（后端用 client-go 按此模板动态创建）
```yaml
apiVersion: batch/v1
kind: Job
metadata: { name: "flutter-build-<task_id>", namespace: default, labels: { app: flutter-build } }
spec:
  backoffLimit: 1              # 失败重试 1 次
  activeDeadlineSeconds: 480   # 8 分钟硬超时熔断
  ttlSecondsAfterFinished: 60  # 完成 1 分钟后清理，按秒计费省钱
  template:
    metadata: { labels: { app: flutter-build } }
    spec:
      restartPolicy: Never
      containers:
        - name: builder
          image: ccr.ccs.tencentyun.com/<命名空间>/flutter-builder:v1
          resources:
            requests: { cpu: "2", memory: "4Gi" }
            limits:   { cpu: "2", memory: "4Gi" }
          env:
            - { name: TASK_ID,        value: "<task_id>" }
            - { name: CODE_B64,       value: "<后端 base64 注入>" }
            - { name: CALLBACK_URL,   value: "http://111.229.28.247/api/v1/internal/preview-callback" }
            - { name: BUILD_TOKEN,    value: "<共享密钥>" }
            - { name: COS_BUCKET,     value: "upload-image-1301638170" }
            - { name: COS_REGION,     value: "ap-guangzhou" }
            - { name: COS_SECRET_ID,  value: "<COS SecretId>" }
            - { name: COS_SECRET_KEY, value: "<COS SecretKey>" }
```
> 费用：ECI 2C4G 按秒计费（≈¥0.158/时），镜像预热后单次 `flutter build web` 约 1~3 分钟，**单次约 ¥0.005~0.02（几分钱）**；1000 次/月仅计算约 ¥5~20。
> ⚠️ 真正的大头往往是**公网 NAT 网关**（若单独常开约 ¥0.2~0.3/时，一月几十到上百）与 **COS 外网下行流量**，不是计算本身——NAT 尽量复用、产物配 CDN、`preview/` 设生命周期定期删。
> 用 `labels.app=flutter-build` + 后端维护并发计数，限制同时跑的 Job 数（默认 2~3），防费用失控。

---

## 4. 阶段四：后端代码改动(这部分我来写，你审阅)

> 下面列的是我要动的点，你不用改代码，但需要给我阶段 5 的配置值。

1. **配置新增**（`config.dev.yaml` / `config.prod.yaml` + `config/config.go`）：
   ```yaml
   build:
     enabled: true
     callback_url: "http://111.229.28.247/api/v1/internal/preview-callback"
     build_token: "<共享密钥>"
     cos_bucket: "upload-image-1301638170"
     cos_region: "ap-guangzhou"
     kubeconfig_path: "/root/DesignCode-AI/kubeconfig.yaml"   # 或 SA token
     image: "ccr.ccs.tencentyun.com/<命名空间>/flutter-builder:v1"
     namespace: "default"
     max_concurrency: 3
     timeout_seconds: 480
   ```
2. **Flutter 触发构建桩**（`internal/workflow/generate_ui_workflow.go`）：Flutter 代码生成后，不阻塞主流程，调用 `PreviewBuilder.Submit(taskID, code)` 建 Job（用 client-go）。
3. **回调路由/handler**（新增 `routes/internal_route.go` + handler）：`POST /api/v1/internal/preview-callback`
   - 校验头 `X-Build-Token`；
   - `status=success` → 按 `task_id` 更新 `Result.iframe_url`，推 SSE `preview_ready{task_id,url}`；
   - `status=failed` → 不写 iframe_url（前端保留 HTML 兜底），可选告警。
4. **SSE 事件**：新增 `preview_ready` 事件类型（复用现有 `sseManager`）。
5. **依赖**：`go.mod` 引入 `k8s.io/client-go`（建 Job 用）。

> 可先做「不依赖云」的一半（回调路由 + SSE + 前端监听），用手动 `curl` 假回调即可端到端验证「存库 + 显示 iframe」；等阶段 2/3 的镜像和集群就绪，再把真实建 Job 接上。

---

## 5. 阶段五：前端改动(我来写)

`frontend/pages/detail/[id].vue`：
- 现有已支持：`Flutter && taskDetail.iframeUrl` → `<iframe>`；GET 详情已能读 `r.iframe_url`。
- 新增：监听 SSE `preview_ready` → 收到即 `taskDetail.iframeUrl = url` → 预览自动从 HTML 兜底切成真机 iframe。

---

## 6. 阶段六：联调验证(按顺序，每步确认再进行下一步)

1. **后端先起、路由通**：部署含回调接口的后端。本机自测：
   ```bash
   curl -X POST http://111.229.28.247/api/v1/internal/preview-callback \
     -H "Content-Type: application/json" -H "X-Build-Token: <共享密钥>" \
     -d '{"task_id":"<某个已存在的task_id>","status":"success","iframe_url":"https://upload-image-1301638170.cos.ap-guangzhou.myqcloud.com/preview/demo/index.html","error":""}'
   ```
   → 查库 `result.iframe_url` 写入；打开该任务详情页 → 出现 iframe（先塞个能公开访问的 html 验证渲染）。
2. **COS 手动放产物验证显示**：把任意 `build/web` 产物传到 `preview/demo/`，公有读能开 `index.html`，详情页 iframe 正常。
3. **单跑一次真实构建**：先绕过后端，本地 `docker run` 手动注入 env 跑一次 `build.sh`，验证：pub get→build web→传 COS→回调全绿。
   ```bash
   docker run --rm -e TASK_ID=docker-test -e CODE_B64=$(base64 -w0 some.dart) \
     -e CALLBACK_URL=... -e BUILD_TOKEN=... -e COS_BUCKET=... -e COS_REGION=... \
     -e COS_SECRET_ID=... -e COS_SECRET_KEY=... \
     ccr.ccs.tencentyun.com/<命名空间>/flutter-builder:v1
   ```
4. **端到端**：网页点生成 Flutter 任务 → 后端建 Job → TKE 拉起容器 → COS 出现 `preview/{task_id}/` → 回调 → SSE → 前端 iframe 显示。
5. **异常路径**：故意用会编译失败的代码 → 验证回调 failed + 前端降级到 HTML 预览；验证超时熔断与重试。

---

## 7. 上线前加固(建议但不阻塞)
- 并发闸门：后端维护运行中 Job 计数，超 `max_concurrency` 排队，避免费用/资源失控。
- 产物 GC：`preview/{task_id}/` 设 COS 生命周期规则（如 30 天自动删）或删任务时连带删。
- 密钥安全：COS/kubeconfig/build_token 放 `config.prod.yaml`(已 gitignore) 或环境变量，不入库不进前端。
- 回调安全：只信 `X-Build-Token`；回调接口幂等（同 task_id 多次回调以最后一次 success 为准）。
- 成本监控：腾讯云费用告警设阈值；ECI 全走按秒。

---

## 8. 我需要你提供/确认的信息(填了我就开始写后端+前端代码)

| 项 | 值 | 说明 |
|---|---|---|
| ① COS 已设公有读私有写 | 是/否 | 阶段一 |
| ② TCR 命名空间 + registry 账号 | ____ | 阶段二推镜像用 |
| ③ 镜像名 | ccr.ccs.tencentyun.com/<ns>/flutter-builder:v1 | 确认命名 |
| ④ TKE Serverless 集群已建 + NAT 公网 | 是/否 | 阶段三 |
| ⑤ kubeconfig 或 build-trigger SA token | 文件/串 | 后端拉 Job 用 |
| ⑥ 回调公网地址可用 | `http://111.229.28.247/api/v1/internal/preview-callback` | 确认能访问 |
| ⑦ 共享密钥 build_token | 你定一个长随机串 | 后端与容器两侧一致 |
| ⑧ Flutter 生成代码类名契约 | 是否统一为 `GeneratedWidget` | 决定改提示词还是改 build.sh |
| ⑨ 并发上限 / 单次超时 | 默认 3 / 480s | 可调 |

---

## 附：常见坑
- **主服务器 OOM**：一定别在这台 2C 机 `flutter build`，全部丢给 ECI。
- **容器拉不到 pub/COS**：NAT 没配或安全组限制出网。
- **iframe 打开空白**：COS 没公有读，或路径不是 `.../preview/{id}/index.html`，或 flutter web 用了 `--web-renderer html` 与资源 base href 不对（用 canvaskit 一般更稳）。
- **回调进不来**：nginx `/api/` 已代理到 Go，路径务必带 `/api/v1/`；`X-Build-Token` 两侧不一致会被拒。
- **生成类名不匹配**：脚手架靠 `GeneratedWidget` 类；模型输出不稳定时，在提示词里强约束或 build.sh 里做归一。
