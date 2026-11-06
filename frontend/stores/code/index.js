import { defineStore } from 'pinia'
import { reactive, computed, toRefs } from 'vue'
import { createUploadActions } from './upload'
import { createTaskActions } from './task'
import { useUserStore } from '~/stores/user'

// 各模式主图上限：单图 = 1，多图 = 8（页序上限）
export const MODE_MAX_IMAGES = { single: 1, multi: 8 }

export const useCodeStore = defineStore('code', () => {
  const route = useRoute()
  const state = reactive({
    // ── 上传 ──（单图/多图各自独立数据桶：切换模式互不影响，不串图）
    modes: {
      single: { images: [], assets: [] },
      multi: { images: [], assets: [] },
    },
    mode: route.query.mode || 'single', // single 单主图 | multi 多主图（页序）
    maxImages: MODE_MAX_IMAGES.single,

    // ── 配置 ──
    config: {
      framework: 'Flutter',
      platform: 'mobile',
      options: ['responsive'],
      advanced: [],
      componentLib: '',
    },

    // ── 并发 ──
    // 单用户活跃任务上限（pending + running）：与后端 config 的
    // queue.max_active_per_user 保持一致（默认 2），前端先拦，后端兜底返回 429
    maxConcurrent: 2,

    // ── 提交 ──
    isSubmitting: false,
  })

  // ═══ Getters ═══

  // 当前模式的主图 / 附图：组件与动作统一读写这两份，切换模式即切换数据源
  const images = computed(() => state.modes[state.mode].images)
  const assets = computed(() => state.modes[state.mode].assets)

  // 服务端实际并发数（从 userStore 读取，跨页面共享）
  const runningCount = computed(() => {
    const userStore = useUserStore()
    const tc = userStore.taskCounts
    return tc.pending + tc.running
  })

  // 是否并发满载
  const isConcurrencyFull = computed(() => runningCount.value >= state.maxConcurrent)

  // 是否已达上传上限
  const isMaxReached = computed(() => images.value.length >= state.maxImages)

  // 是否所有图片已上传
  const allUploaded = computed(() => images.value.length > 0 && images.value.every(img => !!img.cosUrl))

  const uploadHint = computed(() => {
    const userStore = useUserStore()
    if (!userStore.isLogin) return { title: '请先登录', desc: '登录后即可上传设计稿' }
    if (isConcurrencyFull.value) return { title: '已达并发上限', desc: `已有 ${state.maxConcurrent} 个任务排队/执行中，请等待完成后再提交` }
    if (state.isSubmitting) return { title: '正在生成中', desc: '暂不支持操作图片' }
    if (isMaxReached.value) return { title: '已达上传上限', desc: `最多上传 ${state.maxImages} 张图片，请先移除再上传` }
    return { title: '拖拽图片到这里，或点击上传', desc: `支持 PNG、JPG、JPEG，最多 ${state.maxImages} 张，单张不超过 10MB` }
  })

  // ═══ Actions ═══

  const ctx = { state, images, assets, runningCount, isConcurrencyFull, maxConcurrent: state.maxConcurrent }

  const uploadModule = createUploadActions(ctx)
  const taskModule = createTaskActions(ctx)

  // 模式切换：单图 ↔ 多图。两套数据独立保存，切换只改数据源，不动已上传内容
  const setMode = (mode) => {
    if (mode === state.mode || !MODE_MAX_IMAGES[mode] || state.isSubmitting) return
    state.mode = mode
    state.maxImages = MODE_MAX_IMAGES[mode]
  }

  const actions = {
    ...uploadModule,
    ...taskModule,
    setMode,
  }

  // 仅清空当前模式的主图（保留附图素材）：更换/删除主图时用，附件不随主图丢失
  const clearImages = () => {
    if (state.isSubmitting) return
    const bucket = state.modes[state.mode]
    bucket.images.forEach((img) => { if (img.uploadId) uploadModule.cancelUpload(img.uploadId) })
    bucket.images.splice(0)
  }

  // 清除当前模式的主图 + 附图素材（另一模式的数据不受影响）
  const clearAll = () => {
    if (state.isSubmitting) return
    const bucket = state.modes[state.mode]
    bucket.images.forEach((img) => { if (img.uploadId) uploadModule.cancelUpload(img.uploadId) })
    bucket.assets.forEach((a) => { if (a.assetId) uploadModule.cancelAsset(a.assetId) })
    bucket.images.splice(0)
    bucket.assets.splice(0)
    state.isSubmitting = false
  }

  return {
    ...toRefs(state),
    images,
    assets,
    runningCount,
    isConcurrencyFull,
    isMaxReached,
    allUploaded,
    uploadHint,
    setMode,
    clearImages,
    clearAll,
    ...actions,
  }
})
