export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },  // 关键：关闭 devtools
  ssr: false,  // 关闭 SSR 测试
  css: [
    '~/assets/css/main.css',  // 这里配置全局样式
     '@fortawesome/fontawesome-free/css/all.min.css'  // 添加 Font Awesome
  ],
   modules: [
    '@pinia/nuxt',
  ],
  runtimeConfig: {
    // 只能在服务端访问的私有变量
    apiSecret: '', 
    public: {
      // 暴露给客户端的变量
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
      wsBase: process.env.NUXT_PUBLIC_WS_BASE || 'ws://localhost:8080/ws',
      siteTitle: process.env.NUXT_PUBLIC_SITE_TITLE || 'Design2Code AI',
      debug: process.env.NUXT_PUBLIC_DEBUG === 'true',
      cryptoSecret: process.env.NUXT_PUBLIC_CRYPTO_SECRET || 'design2code-ai-secret-key-2025'
    }
  },
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
