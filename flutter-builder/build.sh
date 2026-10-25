#!/usr/bin/env bash
set -euo pipefail

# ── 由后端 Job 注入的环境变量 ──
: "${TASK_ID:?}"
: "${CODE_B64:?}"
: "${CALLBACK_URL:?}"
: "${BUILD_TOKEN:?}"
: "${COS_BUCKET:?}"
: "${COS_REGION:?}"
: "${COS_SECRET_ID:?}"
: "${COS_SECRET_KEY:?}"

IFRAME_PATH="preview/${TASK_ID}"
INDEX_URL="https://${COS_BUCKET}.cos.${COS_REGION}.myqcloud.com/${IFRAME_PATH}/index.html"

post_callback() {
  curl -sS -X POST "$CALLBACK_URL" \
    -H "Content-Type: application/json" \
    -H "X-Build-Token: ${BUILD_TOKEN}" \
    -d "{\"task_id\":\"${TASK_ID}\",\"status\":\"$1\",\"iframe_url\":\"$2\",\"error\":\"$3\"}" || true
}

trap 'post_callback "failed" "" "build error on line $LINENO"' ERR

echo "[build] task=$TASK_ID 解码代码…"
mkdir -p /work && cp -r /opt/template/. /work/
echo "$CODE_B64" | base64 -d > /work/lib/generated_page.dart

echo "[build] flutter pub get…"
cd /work && flutter pub get

echo "[build] flutter build web --release…"
flutter build web --release

echo "[build] 上传产物到 COS…"

# 🌟 终极必杀：伪造它生成的配置格式，并把 stdin 重定向到 /dev/null，彻底杜绝它弹交互模式
mkdir -p /root
cat > /root/.cos.yaml <<EOF
cos:
  base_url: https://${COS_REGION}.cos.myqcloud.com
  secret_id: ${COS_SECRET_ID}
  secret_key: ${COS_SECRET_KEY}
  protocol: https
  buckets:
    - name: ${COS_BUCKET}
      alias: ${COS_BUCKET}
      region: ${COS_REGION}
      endpoint: ${COS_REGION}.cos.myqcloud.com
EOF

# 把 < /dev/null 加上，它就永远别想弹出 Input Your Mode
coscli cp -r /work/build/web "cos://${COS_BUCKET}/${IFRAME_PATH}" < /dev/null

echo "[build] 回调成功 url=$INDEX_URL"
post_callback "success" "$INDEX_URL" ""