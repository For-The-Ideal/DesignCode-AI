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