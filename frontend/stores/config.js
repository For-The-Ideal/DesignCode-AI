import { defineStore } from 'pinia'
import { ref } from 'vue'
import { configApi } from '~/api/config'

// 生成配置字典：由后端 /api/v1/config/options 提供
// 跨页面（生成面板 / 任务筛选）共享，只拉取一次并缓存
export const useConfigStore = defineStore('config', () => {
  const loaded = ref(false)
  const loading = ref(false)

  const frameworks = ref([])
  const platforms = ref([])
  const options = ref([])
  const advanced = ref([])
  const statuses = ref([])
  const componentLibs = ref({}) // { Flutter: [...], React: [...], Vue: [...] }

  // 拉取字典；已加载或正在加载则跳过（force 可强制刷新）
  async function load(force = false) {
    if (loaded.value && !force) return
    if (loading.value) return
    loading.value = true
    try {
      const res = await configApi.getGenConfig()
      if (res && res.code === 200 && res.data) {
        const d = res.data
        frameworks.value = d.frameworks || []
        platforms.value = d.platforms || []
        options.value = d.options || []
        advanced.value = d.advanced || []
        statuses.value = d.statuses || []
        componentLibs.value = d.component_libs || {}
        loaded.value = true
      }
    } finally {
      loading.value = false
    }
  }

  // 指定框架 + 目标平台下适用的组件库（隐藏不适用库）
  function libsFor(framework, platform) {
    const list = componentLibs.value[framework] || []
    if (!platform) return list
    return list.filter(l => l.platform === 'both' || l.platform === platform)
  }

  // ── 字典回显：label / icon / 组合名 均以此处为唯一来源 ──
  const findLabel = (list, val) => (list.find(x => x.value === val) || {}).label || val
  const platformLabel  = v => findLabel(platforms.value, v)
  const frameworkLabel = v => findLabel(frameworks.value, v)
  const optionLabel    = v => findLabel(options.value, v)
  const advancedLabel  = v => findLabel(advanced.value, v)

  // 状态 label：优先字典，字典无值时回退前端生命周期默认（避免后端未 seed 时露英文码）
  const STATUS_LABELS = { success: '成功', running: '生成中', pending: '排队中', failed: '失败' }
  const statusLabel = (v) => {
    const hit = statuses.value.find(x => x.value === v)
    return hit ? hit.label : (STATUS_LABELS[v] || v)
  }
  const platformIcon   = v => (platforms.value.find(x => x.value === v) || {}).icon || 'fas fa-question-circle'

  // 任务显示名：优先图片标题，否则回退“平台 · 框架”
  const taskName = (t) =>
    t?.images?.[0]?.title || `${platformLabel(t?.platform) || '未知'} · ${t?.framework || '未知'}`

  return {
    loaded, loading,
    frameworks, platforms, options, advanced, statuses, componentLibs,
    load, libsFor,
    platformLabel, statusLabel, frameworkLabel, optionLabel, advancedLabel, platformIcon, taskName,
  }
})
