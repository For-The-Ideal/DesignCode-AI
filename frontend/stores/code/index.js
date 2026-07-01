import { defineStore } from 'pinia'
import { reactive, computed, toRefs } from 'vue'
import { createUploadActions } from './upload'
import { createTaskActions } from './task'
import { useUserStore } from '~/stores/user'

export const useCodeStore = defineStore('code', () => {
  const state = reactive({
    // ── 上传 ──
    images: [],
    maxImages: 3,

    // ── 配置 ──
    config: {
      framework: 'Flutter',
      platform: 'mobile',
      options: ['responsive'],
      advanced: [],
      componentLib: '',
    },

    // ── 任务 ──
    // tasks: { [taskId]: { id, target, status, progress, currentStep, result, images, ... } }
    tasks: {},

    // ── 任务统计 ──
    taskCounts: { pending: 0, running: 0, success: 0, failed: 0 },

    // ── 并发 ──
    maxConcurrent: 5,

    // ── 提交 ──
    isSubmitting: false,
  })

  // ═══ Getters ═══

  // 服务端实际并发数（所有会话，非仅当前页面创建的）
  const runningCount = computed(() =>state.taskCounts.pending + state.taskCounts.running)

  // 是否并发满载 = 服务端排队 + 运行 >= 上限
  const isConcurrencyFull = computed(() => runningCount.value >= state.maxConcurrent)

  // 是否已达上传上限 = 图片数量 >= 最大图片数
  const isMaxReached = computed(() => state.images.length >= state.maxImages)

  // 是否所有图片已上传 = 图片数量 > 0 且 每张图片的 cosUrl 都存在
  const allUploaded = computed(() => state.images.length > 0 && state.images.every(img => !!img.cosUrl))

  const uploadHint = computed(() => {
    const userStore = useUserStore()
    if (!userStore.isLogin) return { title: '请先登录', desc: '登录后即可上传设计稿' }
    if (isConcurrencyFull.value) return { title: '已达并发上限', desc: `已有 ${state.maxConcurrent} 个任务排队/执行中，请等待完成后再提交` }
    if (state.isSubmitting) return { title: '正在生成中', desc: '暂不支持操作图片' }
    if (isMaxReached.value) return { title: '已达上传上限', desc: `最多上传 ${state.maxImages} 张图片，请先移除再上传` }
    return { title: '拖拽图片到这里，或点击上传', desc: `支持 PNG、JPG、JPEG，最多 ${state.maxImages} 张，单张不超过 10MB` }
  })

  // ═══ Actions ═══
  const ctx = { state, runningCount, isConcurrencyFull, maxConcurrent: state.maxConcurrent }

  const uploadModule = createUploadActions(ctx)
  const taskModule = createTaskActions(ctx)

  const actions = {
    ...uploadModule,
    ...taskModule,
  }

  // 清除所有
  const clearAll = () => {
    if (state.isSubmitting) return
    taskModule.stopAllPolling()
    state.images = []
    state.tasks = {}
    state.isSubmitting = false
  }

  return {
    ...toRefs(state),
    runningCount,
    isConcurrencyFull,
    isMaxReached,
    allUploaded,
    uploadHint,
    clearAll,
    ...actions,
  }
})
