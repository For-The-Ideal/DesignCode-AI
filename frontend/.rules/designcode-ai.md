---
trigger: always_on
alwaysApply: true
scope: workspace
---

# DesignCode-AI 项目编码规范

> **执行规则**：每次生成代码前，必须逐章自查本规范，确认生成的代码符合所有条款。

## 技术栈

- **框架**: Nuxt 4 (compatibilityVersion 4) / Vue 3.5 + Composition API
- **UI 库**: Element Plus 2.14 (自动导入)
- **状态管理**: Pinia 3 + storeToRefs
- **样式**: Tailwind CSS + SCSS
- **后端**: Go (Gin) 运行在 :8888 端口
- **路径别名**: `~/` → 项目根目录

## 一、文件命名规范

1. **Vue 组件文件**: 使用 PascalCase，如 `AuthModal.vue`、`UploaderImage.vue`
2. **JS/TS 工具文件**: 使用 camelCase，如 `utils.js`、`useGeneration.js`
3. **页面文件**: 使用 camelCase，如 `index.vue`（页面目录名使用 camelCase）
4. **API 文件**: 使用 camelCase，如 `ai.js`、`common.js`
5. **Composable 文件**: `useXxx.js`，如 `useEnv.js`、`useSSE.js`

## 二、Vue 组件规范

- 所有组件**必须使用 `<script setup>`** 语法
- Props/Emits 使用驼峰命名
- 样式使用 `<style scoped>`（SCSS 需加 `lang="scss"`）

## 三、请求架构规范

### 3.1 双通道请求

```
浏览器 ──POST /api/ajaxData──▶ ajaxData.js ──▶ Go :8888
  body: { url, method, params }        req.client=true → 明文直取
  header: client:true

SSR ──POST /api/ajaxData──▶ ajaxData.js ──▶ Go :8888
  body: { aes: setEncrypt(...) }       req.client=false → getDecrypt()
```

### 3.2 请求规则

1. **客户端请求**: `utils/request.js` 统一封装，`body` 明文 `{ url, method, params }`，header 带 `client: true`
2. **SSR 请求**: 同一 `utils/request.js`，`body` 加密为 `{ aes: setEncrypt({ url, method, params }) }`
3. **服务端代理**: `server/api/ajaxData.js` 根据 `req.client` 自动解密或透传
4. **API 模块**: `api/*.js` 封装具体的接口调用，内部调用 `utils/request.js`

### 3.3 环境适配

1. **本地开发**: `.env` → `NUXT_PUBLIC_API_BASE=http://localhost:8888`
2. **生产部署**: `.env.production` → `NUXT_PUBLIC_API_BASE=http://IP:8888`
3. **Nuxt 4 composable**: `useRuntimeConfig()` 在 `.vue` 的 `<script setup>` 中自动可用；在 `.js` composable 中同样可直接调用；`process.env.NUXT_PUBLIC_API_BASE` 在服务端工具文件中可用

## 四、请求工具规范

1. 所有 HTTP 请求通过 `utils/request.js` 的 `httpRequest` 实例发起
2. API 文件放在 `api/` 目录下，按功能模块拆分
3. **禁止**在组件中直接调用 `$fetch` 或 `axios`

## 五、Composable 规范

1. Composable 放在 `composables/` 下，文件名 `useXxx.js`
2. 使用命名导出 `export const useXxx = () => { ... }`
3. 业务逻辑集中处理，组件仅负责渲染和事件触发

## 六、环境变量规范

1. **公开变量**: 使用 `NUXT_PUBLIC_` 前缀，定义在 `nuxt.config.ts` 的 `runtimeConfig.public` 中
2. **私有变量**: 无前缀的服务端变量，仅在 SSR 端可访问
3. **禁止**硬编码 API 地址、域名等环境相关值
4. **开发环境**: 运行 `npm run dev --dotenv .env.dev`

## 七、目录结构约定

```
frontend/
├── api/              # API 接口封装
├── assets/           # 静态资源
├── components/       # Vue 组件
├── composables/      # Composable 函数
├── config/           # 配置文件（localUrl 等）
├── pages/            # 页面（Nuxt 文件路由）
├── server/           # Nitro 服务端
│   ├── api/          # SSR API 代理
│   └── utils/        # 服务端工具（加密解密、请求转发）
├── stores/           # Pinia stores
├── utils/            # 客户端工具函数
├── .env              # 开发环境变量
└── .env.production   # 生产环境变量
```