import monacoEditorPlugin from 'vite-plugin-monaco-editor'

export default defineNuxtConfig({
  ssr: true,
  devtools: { enabled: false },

  future: {
    compatibilityVersion: 4,
  },

  // ========== 自动导入配置 ==========
  imports: {
    dirs: ['composables/**'],
  },

  modules: [
    '@pinia/nuxt',
    '@nuxt/image',
    '@nuxtjs/tailwindcss',
    '@nuxtjs/color-mode',
  ],

  css: [
    '~/assets/css/main.css',
    '@fortawesome/fontawesome-free/css/all.min.css',
  ],


  build: {
    transpile: ['monaco-editor'],
  },

  // ========== Vite 配置 ==========
  vite: {
    optimizeDeps: {
      include: [
        'element-plus',
        'element-plus/dist/locale/zh-cn.mjs',
        // 'swiper',
        // 'swiper/vue',
        // 'video.js',
        // 'vue-color',
        // 'vue-video-player',
        // 'vue3-seamless-scroll',
      ],
      // exclude: ['@vueuse/core'],
    },
    plugins: [
      monacoEditorPlugin({
        languageWorkers: ['editorWorkerService', 'typescript', 'json', 'css', 'html'],
      }),
    ],
    build: {
      sourcemap: false,
      chunkSizeWarningLimit: 2 * 1024,
      rollupOptions: {
        onwarn(warning, warn) {
          if (warning.code === 'INVALID_ANNOTATION' && warning.message.includes('__PURE__')) return
          if (warning.message?.includes('module-preload-polyfill') && warning.message?.includes('Sourcemap')) return
          if (warning.code === 'FILE_SIZE_EXCEEDS_LIMIT') return
          warn(warning)
        },
        output: {
          manualChunks(id) {
            if (id.includes('node_modules')) {
              // // ECharts 核心库
              // if (id.includes('echarts') || id.includes('zrender')) {
              //   if (id.includes('vue-echarts')) return
              //   return 'echarts'
              // }
              // // Swiper 核心库
              // if (id.includes('swiper')) {
              //   if (id.includes('vue-awesome-swiper')) return
              //   return 'swiper'
              // }
              // 工具库（独立、无 Vue 依赖）
              if (id.includes('crypto-js')) {
                return 'utils'
              }
              // // 后续安装 moment / jszip / clipboard 时可取消注释
              // if (id.includes('moment') || id.includes('jszip') || id.includes('clipboard')) {
              //   return 'utils'
              // }
              // // 动画库
              // if (id.includes('gsap')) return 'gsap'
            }
          },
        },
      },
    },
    // css: {
    //   preprocessorOptions: {
    //     scss: {
    //       additionalData: `@use "~/assets/scss/variables.scss" as *;@use "~/assets/scss/config.scss" as *;`
    //     }
    //   }
    // },
    server: {
      warmup: {
        clientFiles: ['./pages/**/*.vue', './layouts/**/*.vue'],
      },
    },
    logLevel: 'warn',
  },

  // ========== Nitro 服务端配置 ==========
  nitro: {
    sourceMap: false,
  },

  compatibilityDate: '2026-03-31',
})