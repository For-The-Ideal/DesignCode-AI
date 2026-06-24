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

  // ========== ✅ 正确的运行时配置 ==========
  runtimeConfig: {
    // 私有配置 (仅服务端可访问)
    apiSecret: '',
    // 公开配置 (客户端和服务端都可访问)
    public: {
      apiBase: '',
      wsBase: '',
      debug: false,
      cryptoKey: '',
      cryptoIv: '',
    },
  },

  // ========== Vite 配置 ==========
  vite: {
    optimizeDeps: {
      include: [
        'element-plus',
        'element-plus/dist/locale/zh-cn.mjs',
      ],
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
              if (id.includes('crypto-js')) {
                return 'utils'
              }
            }
          },
        },
      },
    },
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