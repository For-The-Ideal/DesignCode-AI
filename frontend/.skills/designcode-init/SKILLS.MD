---
name: designcode-init
description: 初始化 Design2Code AI 完整项目结构
type: skill
trigger: manual
---

## 任务目标
在当前项目根目录下，创建完整的 Design2Code AI 项目结构，包含 Nuxt3 前端和 Golang 后端。

## 执行步骤

### 1. 创建项目根目录结构

```bash
# 创建前端目录
mkdir -p frontend/pages
mkdir -p frontend/components/layout
mkdir -p frontend/components/upload
mkdir -p frontend/components/code
mkdir -p frontend/components/evaluation
mkdir -p frontend/components/comparison
mkdir -p frontend/components/common
mkdir -p frontend/components/auth
mkdir -p frontend/composables
mkdir -p frontend/stores
mkdir -p frontend/utils
mkdir -p frontend/assets/css
mkdir -p frontend/assets/images
mkdir -p frontend/plugins
mkdir -p frontend/layouts
mkdir -p frontend/server/api
mkdir -p frontend/server/middleware

