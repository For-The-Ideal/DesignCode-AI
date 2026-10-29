#!/usr/bin/env bash
# DesignCode-AI 服务器发布入口（唯一入口：本地改这个文件，服务器就跟着变）
#   bash deploy/deploy.sh                全量：拉代码 + 编译后端 + 构建前端 + 重启 + 自检
#   bash deploy/deploy.sh --frontend     只改前端：拉代码 + 依赖 + 构建 + 重启前端
#   bash deploy/deploy.sh --backend      只改后端：拉代码 + 编译 + 重启后端
#   bash deploy/deploy.sh --no-pull      不拉代码，只按上面的方式重新构建
# 可组合：bash deploy/deploy.sh --frontend --no-pull
set -uo pipefail

# 先把自己拷到 /tmp 再执行：git reset 会覆盖仓库里的脚本，不能边改边跑
if [ -z "${DEPLOY_REEXEC:-}" ]; then
  export DEPLOY_REEXEC=1
  export DEPLOY_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
  cp "${BASH_SOURCE[0]}" /tmp/designcode-deploy.sh
  exec bash /tmp/designcode-deploy.sh "$@"
fi

ROOT="${DEPLOY_ROOT:-/root/DesignCode-AI}"
cd "$ROOT" || { echo "!!! 目录不存在: $ROOT"; exit 1; }
# shellcheck disable=SC1091
source "$ROOT/deploy/versions.env"

MODE=all
NOPULL=0
for a in "$@"; do
  case "$a" in
    --frontend) MODE=frontend ;;
    --backend)  MODE=backend ;;
    --no-pull)  NOPULL=1 ;;
    -h|--help)  sed -n '2,7p' "$ROOT/deploy/deploy.sh"; exit 0 ;;
    *) echo "未知参数: $a（可用: --frontend --backend --no-pull）"; exit 1 ;;
  esac
done

step() { echo; echo "===== $* ====="; }

install_node() {
  local want="v$NODE_VERSION"
  if [ -x "$NODE_HOME/bin/node" ] && [ "$("$NODE_HOME/bin/node" -v 2>/dev/null)" = "$want" ]; then
    echo "Node $want 已就绪 ($NODE_HOME)"
    return 0
  fi
  echo ">>> 安装 Node $want -> $NODE_HOME"
  local tgz="node-$want-linux-x64.tar.gz" base
  rm -f "/tmp/$tgz"
  for base in https://cdn.npmmirror.com/binaries/node https://nodejs.org/dist; do
    curl -fsSL --retry 3 -o "/tmp/$tgz" "$base/$want/$tgz" && break
    echo "  下载失败: $base"
  done
  [ -s "/tmp/$tgz" ] || { echo "!!! Node 下载失败"; return 1; }
  rm -rf "$NODE_HOME" && mkdir -p "$NODE_HOME"
  tar -xzf "/tmp/$tgz" -C "$NODE_HOME" --strip-components=1 || return 1
  echo 'export PATH=/usr/local/node/bin:/usr/local/go/bin:$PATH' > /etc/profile.d/designcode.sh
  "$NODE_HOME/bin/node" -v || return 1
}

sync_units() {
  local n
  for n in designcode-api designcode-frontend; do
    ln -sfn "$ROOT/deploy/$n.service" "/etc/systemd/system/$n.service"
  done
  systemctl daemon-reload
  systemctl enable designcode-api designcode-frontend >/dev/null 2>&1 || true
  ls -l /etc/systemd/system/designcode-*.service
}

# npm 有时跳过 optional 依赖，缺平台包时直接解压安装
ensure_oxc() {
  local pkg="$1" ver="$2" dir="node_modules/$1" tmp
  [ -d "$dir" ] && { echo "  ok  $pkg"; return 0; }
  tmp=$(mktemp -d)
  ( cd "$tmp" && npm pack "$pkg@$ver" >/dev/null 2>&1 && tar -xzf ./*.tgz ) \
    || { echo "  !!! 下载失败 $pkg@$ver"; rm -rf "$tmp"; return 1; }
  mkdir -p "$(dirname "$dir")"
  rm -rf "$dir"
  mv "$tmp/package" "$dir"
  rm -rf "$tmp"
  echo "  ok  $pkg (解压安装)"
}

step "0 版本对齐"
install_node || exit 1
export PATH="$NODE_HOME/bin:$GO_HOME/bin:$PATH"
node -v; npm -v; go version || { echo "!!! go 未安装"; exit 1; }

if [ "$NOPULL" = "1" ]; then
  step "1 跳过 git（--no-pull）"
else
  step "1 拉取代码（服务器不做任何手工改动）"
  git fetch origin main && git reset --hard origin/main || { echo "!!! git 失败"; exit 1; }
  git log --oneline -1
fi

step "2 同步 systemd 单元（软链到仓库）"
sync_units

if [ "$MODE" != "frontend" ]; then
  step "3 编译后端"
  cd "$ROOT/frontend_api" || exit 1
  go build -trimpath -ldflags="-s -w" -o server.new main.go || { echo "!!! 后端编译失败"; exit 1; }
  mv -f server.new server
  ls -lh server | awk '{print "  server", $5}'
else
  step "3 跳过后端编译（--frontend）"
fi

if [ "$MODE" != "backend" ]; then
  step "4 安装前端依赖"
  cd "$ROOT/frontend" || exit 1
  rm -rf node_modules .nuxt
  npm ci --include=optional --no-audit --no-fund \
    || echo ">>> npm ci 退出码非 0（postinstall 可能失败），继续检查依赖树"

  OXV=$(node -p "require('./package.json').optionalDependencies['@oxc-parser/binding-linux-x64-gnu']" 2>/dev/null || echo "0.156.0")
  echo "  oxc 平台包版本: $OXV"
  ensure_oxc "@oxc-parser/binding-linux-x64-gnu" "$OXV"
  ensure_oxc "@oxc-transform/binding-linux-x64-gnu" "$OXV"
  ensure_oxc "@oxc-minify/binding-linux-x64-gnu" "$OXV"
  node --input-type=module -e "await import('oxc-parser');console.log('  oxc native ok')" \
    || { echo "!!! 原生绑定仍缺失"; exit 1; }

  step "5 构建前端（期间停 MySQL 腾内存）"
  systemctl stop mysqld 2>/dev/null || true
  npm run build
  rc=$?
  systemctl start mysqld 2>/dev/null || true
  sleep 3
  ss -tlnp | grep -q ':3306 ' && echo "  MySQL 3306 正常" || echo "  !!! MySQL 未监听 3306"
  [ $rc -eq 0 ] || { echo "!!! 前端构建失败 rc=$rc"; exit $rc; }
  ls -lh .output/server/index.mjs
else
  step "4-5 跳过前端构建（--backend）"
fi

step "6 重启并自检（模式: $MODE）"
case "$MODE" in
  backend)  UNITS="designcode-api" ;;
  frontend) UNITS="designcode-frontend" ;;
  *)        UNITS="designcode-api designcode-frontend" ;;
esac
systemctl restart $UNITS
sleep 4

OK=1
for u in designcode-api designcode-frontend; do
  printf '  %-22s %s\n' "$u" "$(systemctl is-active "$u")"
done
for u in $UNITS; do
  [ "$(systemctl is-active "$u")" = "active" ] || { echo "  !!! $u 不是 active"; OK=0; }
done

if [ "$MODE" != "frontend" ]; then
  c=$(curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:8888/api/v1/config/options)
  echo "  backend  :$c"
  [ "$c" = "200" ] || OK=0
fi
if [ "$MODE" != "backend" ]; then
  c=$(curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:3000)
  echo "  frontend :$c"
  [ "$c" = "200" ] || OK=0
fi
c=$(curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1/)
echo "  gateway  :$c"
[ "$c" = "200" ] || OK=0

if nginx -T 2>/dev/null | grep -q 'proxy_pass http://127.0.0.1:8888'; then
  echo '  api 路由 :已配置 -> 8888'
else
  echo '  api 路由 :!! 未配置'
  OK=0
fi

case "$MODE" in
  frontend) NAME="前端" ;;
  backend)  NAME="后端" ;;
  *)        NAME="全量" ;;
esac

echo
echo "=================================================="
if [ "$OK" = "1" ]; then
  echo "  $NAME部署成功"
  echo "=================================================="
else
  echo "  $NAME部署失败 —— 检查上面非 200 或标 !!! 的项"
  echo "=================================================="
  exit 1
fi
