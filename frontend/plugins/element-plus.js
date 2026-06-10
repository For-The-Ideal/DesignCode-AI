/**
 * Element Plus 全栈插件
 * 
 * 在 SSR 和客户端都注册 Element Plus 组件
 * 并提供 SSR 必需的 ID 和 ZIndex injection
 * 
 * 注意：ID_INJECTION_KEY 使用固定的 current 值，
 * 确保 SSR 和客户端水合时生成的 ID 一致，避免 hydration mismatch
 */
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import { ID_INJECTION_KEY, ZINDEX_INJECTION_KEY } from 'element-plus'

export default defineNuxtPlugin((nuxtApp) => {
  const isServer = import.meta.server
  
  nuxtApp.vueApp.use(ElementPlus, {
    locale: zhCn,
    // SSR 时禁用 teleport，避免水合不匹配
    zIndex: 2000,
  })
  
  // 提供 SSR 必需的 injection
  // 使用固定的 current 值确保 SSR 和客户端 ID 一致
  // 每次 SSR 都从 0 开始，保证 ID 生成的确定性
  nuxtApp.vueApp.provide(ID_INJECTION_KEY, {
    prefix: 1024,
    current: 0,
  })
  
  nuxtApp.vueApp.provide(ZINDEX_INJECTION_KEY, {
    current: 0,
  })
  
  // 在客户端水合后，确保 popper 容器存在
  if (!isServer) {
    nuxtApp.hook('app:mounted', () => {
      // 确保 el-popper-container 存在
      let popperContainer = document.getElementById('el-popper-container')
      if (!popperContainer) {
        popperContainer = document.createElement('div')
        popperContainer.id = 'el-popper-container'
        document.body.appendChild(popperContainer)
      }
    })
  }
})
