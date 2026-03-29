export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },  // 关键：关闭 devtools
  ssr: false,  // 关闭 SSR 测试
  modules: [],  // 清空所有模块
  css: [
    '~/assets/css/main.css',  // 这里配置全局样式
     '@fortawesome/fontawesome-free/css/all.min.css'  // 添加 Font Awesome
  ],
   app: {
    head: {
      link: [
        { rel: 'stylesheet', href: 'https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css' }
      ]
    }
  },
  typescript: {
    typeCheck: false
  },
   build: {
    transpile: ['monaco-editor']
  },
  
  vite: {
    optimizeDeps: {
      include: ['monaco-editor']
    }
  }
})

// // nuxt.config.ts
// export default defineNuxtConfig({
//   compatibilityDate: '2025-07-15',
//   // devtools: { enabled: true },
  
//   // 模块配置
//   // modules: [
//   //   '@nuxtjs/tailwindcss',
//   //   '@nuxt/image',
//   //   '@pinia/nuxt',
//   //   '@nuxtjs/color-mode'
//   // ],
  
//   // ⭐ 关键：全局加载 CSS
//   css: [
//     '~/assets/css/main.css'  // 这里配置全局样式
//   ],
  
//   // 应用配置
//   app: {
//     head: {
//       title: 'Design2Code AI',
//       meta: [
//         { charset: 'utf-8' },
//         { name: 'viewport', content: 'width=device-width, initial-scale=1' },
//         { name: 'description', content: 'AI 驱动的设计稿转代码工具' }
//       ],
//       link: [
//         { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
//         { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
//         { href: 'https://fonts.googleapis.com/css2?family=Inter:opsz,wght@14..32,300;14..32,400;14..32,500;14..32,600;14..32,700;14..32,800&display=swap', rel: 'stylesheet' }
//       ]
//     }
//   },
  
//   // // TailwindCSS 配置
//   // tailwindcss: {
//   //   cssPath: '~/assets/css/main.css',  // 如果使用 Tailwind，可以在这里配置
//   //   configPath: 'tailwind.config.ts',
//   //   viewer: true
//   // },
  
//   // runtimeConfig: {
//   //   public: {
//   //     apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api'
//   //   }
//   // },
  
//   // typescript: {
//   //   typeCheck: false,
//   //   strict: false
//   // },
  
//   // devServer: {
//   //   port: 3000,
//   //   host: '0.0.0.0'
//   // }
// })