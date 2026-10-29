#!/usr/bin/env bash
# DesignCode-AI 服务器发布入口（唯一入口：本地改这个文件，服务器就跟着变）
#   bash deploy/deploy.sh              拉代码 + 装依赖 + 构建 + 重启 + 自检
#   bash deploy/deploy.sh --no-pull    跳过 git，只构建并重启
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

step "0/6 版本对齐"
install_node || exit 1
export PATH="$NODE_HOME/bin:$GO_HOME/bin:$PATH"
node -v; npm -v; go version || { echo "!!! go 未安装"; exit 1; }

if [ "${1:-}" != "--no-pull" ]; then
  step "1/6 拉取代码（服务器不做任何手工改动）"
  git fetch origin main && git reset --hard origin/main || { echo "!!! git 失败"; exit 1; }
  git log --oneline -1
else
  step "1/6 跳过 git（--no-pull）"
fi

step "2/6 同步 systemd 单元（软链到仓库）"
sync_units

step "3/6 编译后端"
cd "$ROOT/frontend_api" || exit 1
go build -trimpath -ldflags="-s -w" -o server.new main.go || { echo "!!! 后端编译失败"; exit 1; }
mv -f server.new server
ls -lh server | awk '{print "  server", $5}'

step "4/6 安装前端依赖"
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

step "5/6 构建前端（期间停 MySQL 腾内存）"
systemctl stop mysqld 2>/dev/null || true
npm run build
rc=$?
systemctl start mysqld 2>/dev/null || true
sleep 3
ss -tlnp | grep -q ':3306 ' && echo "  MySQL 3306 正常" || echo "  !!! MySQL 未监听 3306"
[ $rc -eq 0 ] || { echo "!!! 前端构建失败 rc=$rc"; exit $rc; }
ls -lh .output/server/index.mjs

step "6/6 重启并自检"
systemctl restart designcode-api designcode-frontend
sleep 4
for u in designcode-api designcode-frontend; do
  printf '  %-22s %s\n' "$u" "$(systemctl is-active "$u")"
done
curl -s -o /dev/null -w '  backend  :%{http_code}\n' http://127.0.0.1:8888/api/v1/config/options
curl -s -o /dev/null -w '  frontend :%{http_code}\n' http://127.0.0.1:3000
curl -s -o /dev/null -w '  gateway  :%{http_code}\n' http://127.0.0.1/
curl -s -o /dev/null -w '  api80    :%{http_code}\n' http://127.0.0.1/api/v1/config/options
echo
echo "发布结束：四个数字都是 200 才算成功"
