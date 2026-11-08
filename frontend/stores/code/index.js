import { defineStore } from 'pinia'
import { reactive, computed, toRefs } from 'vue'
import { createUploadActions } from './upload'
import { createTaskActions } from './task'
import { createRelationActions } from './relation'
import { useUserStore } from '~/stores/user'

// 各模式主图上限：单图 = 1，多图 = 5（页序上限；与任务级并发 2 无关，
// 一次多图生成是「一个任务内逐页串行」，不会占用多个并发名额）
export const MODE_MAX_IMAGES = { single: 1, multi: 5 }

// 模式归一化：URL/参数只接受 single / multi，其他值一律回退单图
// （state.modes 只有这两个桶，非法值会让 images/assets 取到 undefined）
export const normalizeMode = (mode) => (MODE_MAX_IMAGES[mode] ? mode : 'single')

export const useCodeStore = defineStore('code', () => {
  const route = useRoute()
  // 初始模式来自 URL，上限必须跟着初始化，否则多图入口进来会被卡在 1 张
  const initialMode = normalizeMode(route.query.mode)
  const state = reactive({
    // ── 上传 ──（单图/多图各自独立数据桶：切换模式互不影响，不串图）
    modes: {
      single: { images: [], assets: [], pageLinks: [] },
      multi: { images: [], assets: [], pageLinks: [] },
    },
    mode: initialMode, // single 单主图 | multi 多主图（页序），来源：URL ?mode=
    maxImages: MODE_MAX_IMAGES[initialMode],
    // 多图模式当前编辑页（pageId）；单图忽略。切页只换数据源，不丢任何一页的标注
    activePageId: '',

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
  // 页面关系边表（多图模式），与主图同桶，切模式互不影响
  const pageLinks = computed(() => state.modes[state.mode].pageLinks)
  // 当前编辑的主图：多图按 activePageId 取页（找不到回退第一页），单图恒为唯一主图
  const currentImage = computed(() => {
    if (!images.value.length) return null
    return images.value.find((i) => i.pageId === state.activePageId) || images.value[0]
  })

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

  const ctx = { state, images, assets, pageLinks, runningCount, isConcurrencyFull, maxConcurrent: state.maxConcurrent }

  const uploadModule = createUploadActions(ctx)
  const relationModule = createRelationActions(ctx)
  // 关系动作挂到 ctx，供 task 提交时先迁移旧 parentPageId 数据
  Object.assign(ctx, relationModule)
  const taskModule = createTaskActions(ctx)

  // 模式切换：单图 ↔ 多图。两套数据独立保存，切换只改数据源，不动已上传内容
  const setMode = (mode) => {
    if (mode === state.mode || !MODE_MAX_IMAGES[mode] || state.isSubmitting) return
    state.mode = mode
    state.maxImages = MODE_MAX_IMAGES[mode]
    // 当前页跟随目标模式的第一页，画布/树随之恢复（由工作台会话 watch 完成）
    state.activePageId = state.modes[mode].images[0]?.pageId || ''
  }

  // 切到多图的某一页（单图忽略）；传入非法 pageId 不生效
  const setActivePage = (pageId) => {
    if (images.value.some((i) => i.pageId === pageId)) state.activePageId = pageId
  }

  // ═══ 多图画板：页面摆放与页序 ═══
  // board = { x, y } 仅用于画布排版（不提交给生成）；null = 只在左侧列表，尚未拖入画布
  const placeOnBoard = (pageId, x, y) => {
    const img = images.value.find((i) => i.pageId === pageId)
    if (!img) return
    img.board = { x: Math.round(x), y: Math.round(y) }
  }

  // 拖动已在画布上的页（未放置的页不响应移动）
  const moveOnBoard = (pageId, x, y) => {
    const img = images.value.find((i) => i.pageId === pageId)
    if (!img || !img.board) return
    img.board = { x: Math.round(x), y: Math.round(y) }
  }

  // 默认把第一页放在画布左上（进入多图/首次上传时）；已有页在画布上则不动
  // 落点靠近左上角：单张图时卡片连同标题栏完整落在视口内，不产生多余滚动
  const ensureFirstOnBoard = () => {
    if (state.mode !== 'multi' || !images.value.length) return
    if (images.value.some((i) => i.board)) return
    images.value[0].board = { x: 24, y: 24 }
  }

  // 页面列表上下拖动排序：只改页序（P1..Pn 与提交 sort_order），pageId 不变
  const moveImage = (fromIdx, toIdx) => {
    const list = images.value
    if (fromIdx === toIdx) return
    if (fromIdx < 0 || fromIdx >= list.length || toIdx < 0 || toIdx >= list.length) return
    const [item] = list.splice(fromIdx, 1)
    list.splice(toIdx, 0, item)
  }

  const actions = {
    ...uploadModule,
    ...relationModule,
    ...taskModule,
    setMode,
  }

  // 仅清空当前模式的主图（保留附图素材）：更换/删除主图时用，附件不随主图丢失
  const clearImages = () => {
    if (state.isSubmitting) return
    const bucket = state.modes[state.mode]
    bucket.images.forEach((img) => { if (img.uploadId) uploadModule.cancelUpload(img.uploadId) })
    bucket.images.splice(0)
    bucket.pageLinks.splice(0)
    state.activePageId = ''
  }

  // 清除当前模式的主图 + 附图素材（另一模式的数据不受影响）
  const clearAll = () => {
    if (state.isSubmitting) return
    const bucket = state.modes[state.mode]
    bucket.images.forEach((img) => { if (img.uploadId) uploadModule.cancelUpload(img.uploadId) })
    bucket.assets.forEach((a) => { if (a.assetId) uploadModule.cancelAsset(a.assetId) })
    bucket.images.splice(0)
    bucket.assets.splice(0)
    bucket.pageLinks.splice(0)
    state.activePageId = ''
    state.isSubmitting = false
  }

  return {
    ...toRefs(state),
    images,
    currentImage,
    assets,
    pageLinks,
    runningCount,
    isConcurrencyFull,
    isMaxReached,
    allUploaded,
    uploadHint,
    setMode,
    setActivePage,
    placeOnBoard,
    moveOnBoard,
    ensureFirstOnBoard,
    moveImage,
    clearImages,
    clearAll,
    ...actions,
  }
})
