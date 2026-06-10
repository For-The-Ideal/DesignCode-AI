/**
 * 注册全局自定义指令
 * Nuxt 3 会自动加载 plugins/ 目录下的文件
 */
import directives from '~/directives/index.js'

export default defineNuxtPlugin((nuxtApp) => {
  Object.entries(directives).forEach(([name, directive]) => {
    nuxtApp.vueApp.directive(name, directive)
  })
})
