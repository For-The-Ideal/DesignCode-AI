<template>
  <div class="tasks-page">
    <AppSidebar
      v-model="sidebarOpen"
      brand="任务中心"
      expandedWidth="225px"
      :navItems="navItems"
      @navClick="handleNavClick"
    >
      <TasksStatusOverview :statusList="statusList" />
    </AppSidebar>

    <main class="task-main">
      <!-- ═══ 页面头 ═══ -->
      <div class="page-top">
        <div class="top-left">
          <h2 class="page-title">任务列表</h2>
        </div>
        <el-button type="primary" :icon="Plus" round @click="handleCreate">
          新建任务
        </el-button>
      </div>

      <!-- ═══ 筛选表单 ═══ -->
      <TasksFilter v-model="filterForm" />

      <!-- ═══ 任务卡片网格 ═══ -->
      <TasksGrid
        :tasks="filteredTasks"
        :isEmptyFiltered="!!filterForm.search || !!filterForm.framework || !!filterForm.platform || !!filterForm.status"
        v-model:page="page"
        @detail="goToDetail"
        @delete="handleDelete"
      />
    </main>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import AppSidebar from '~/components/layout/AppSidebar.vue'
import TasksFilter from '~/components/tasks/TasksFilter.vue'
import TasksGrid from '~/components/tasks/TasksGrid.vue'
import TasksStatusOverview from '~/components/tasks/TasksStatusOverview.vue'
import { useMenuListStore } from '~/stores/menuList'
import { useConfigStore } from '~/stores/config'
import { storeToRefs } from 'pinia'
import { useUserStore } from '~/stores/user'
import { taskApi } from '~/api/task'
const router = useRouter()
const route = useRoute()
const sidebarOpen = ref(true)
const userStore = useUserStore()
const { taskCounts, taskProgress } = storeToRefs(userStore)
const configStore = useConfigStore()
onMounted(() => configStore.load())

// ═══ 状态列表 ═══
// 状态枚举键为固定领域常量（与后端 Status* / 概览 dot--{key} 对齐），行不依赖字典；
// label 仍唯一来自配置字典（未加载/后端未重启则回退键值），count 来自 taskCounts
const STATUS_ORDER = ['success', 'running', 'pending', 'failed']
const statusList = computed(() =>
  STATUS_ORDER.map(k => ({ key: k, label: configStore.statusLabel(k), count: taskCounts.value[k] ?? 0 }))
)
// ═══ 导航 ═══
const menuStore = useMenuListStore()
// 当前路径在 setup 里取，交给 store getter 算 active（getter 里不能 useRouter）
const navItems = computed(() => menuStore.navItemsWithActive(route.path))
const handleNavClick = (item) => {
  if (item.active) return
  router.push(item.to)
}

// ═══ 筛选 / 分页（参数从 URL 初始化，刷新与分享链接保持状态）═══
const filterForm = ref({
  search: String(route.query.search || ''),
  framework: String(route.query.framework || ''),
  platform: String(route.query.platform || ''),
  status: String(route.query.status || ''),
})
const page = ref(Math.max(1, Number(route.query.page) || 1))
const loading = ref(false)

// 筛选条件变化时回到第一页（翻页本身不受影响）
watch(filterForm, () => { page.value = 1 }, { deep: true })

// URL 同步：搜索输入 / 筛选项 / 翻页变化时写回 query；默认值（第 1 页、空筛选）不写入，保持 URL 干净
watch(
  [filterForm, page],
  ([f, p]) => {
    const query = { ...route.query }
    const setOrRemove = (key, val) => {
      if (val) query[key] = String(val)
      else delete query[key]
    }
    setOrRemove('search', f.search)
    setOrRemove('framework', f.framework)
    setOrRemove('platform', f.platform)
    setOrRemove('status', f.status)
    setOrRemove('page', p > 1 ? p : '')
    router.replace({ query })
  },
  { deep: true },
)

// ═══ 操作 ═══
const goToDetail = row => router.push(`/detail/${row.id}`)
const handleCreate = () => router.push('/code')

const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定删除任务"${configStore.taskName(row)}"吗？此操作不可恢复。`,
    '删除确认',
    { confirmButtonText: '确定删除', cancelButtonText: '取消', type: 'warning' }
  ).then(async () => {
    const res = await taskApi.deleteTask(row.id)
    if (res.code === 200) {
      const idx = tasks.value.findIndex(t => t.id === row.id)
      if (idx !== -1) tasks.value.splice(idx, 1)
      ElMessage.success('任务已删除')
    }
  }).catch(() => {})
}

// ═══ 数据获取 ═══
const tasks = ref([])

const fetchTasks = async () => {
  loading.value = true
  try {
    const res = await taskApi.getTaskList()
    if (res.code === 200 && res.data) {
      tasks.value = res.data
    }
  } catch {
    // 静默
  } finally {
    loading.value = false
  }
}


// ═══ 筛选计算 ═══
const filteredTasks = computed(() => {
  let r = tasks.value
  const f = filterForm.value
  if (f.search) {
    const q = f.search.toLowerCase()
    r = r.filter(t =>
      configStore.taskName(t).toLowerCase().includes(q) ||
      t.id.toLowerCase().includes(q) ||
      (t.framework || '').toLowerCase().includes(q)
    )
  }
  if (f.framework) r = r.filter(t => t.framework === f.framework)
  if (f.platform)  r = r.filter(t => t.platform === f.platform)
  if (f.status)    r = r.filter(t => t.status === f.status)
  return r
})

fetchTasks()

// ═══ 实时进度：将用户级 SSE 推来的 taskProgress 写回对应行，进度条自行推进（无轮询）═══
watch(taskProgress, (map) => {
  if (!map) return
  for (const t of tasks.value) {
    const live = map[t.id]
    if (!live) continue
    if (typeof live.progress === 'number') t.progress = live.progress
    if (live.status) t.status = live.status
    if (live.current_step) t.current_step = live.current_step
  }
}, { immediate: true })
</script>

<style scoped>
/* ═══ 页面布局 ═══ */
.tasks-page {
  display: flex;
  height: 100%;
  background: #0a0a0f;
}

.task-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 24px 28px;
  min-height: 0;
  overflow-y: auto;
  animation: fadeIn 0.35s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(6px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* ═══ 页面头 ═══ */
.page-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
  flex-shrink: 0;
}
.top-left { display: flex; align-items: center; gap: 12px; }
.page-title {
  font-size: 22px;
  font-weight: 700;
  color: #e2e8f0;
  margin: 0;
  letter-spacing: -0.3px;
}
.count-tag { font-size: 11px; opacity: 0.7; }

/* ═══ 响应式 ═══ */
@media (max-width: 1100px) {
  .task-main { padding: 16px; }
}
</style>

<style>
/* ═══ 非 scoped：覆盖 teleported 下拉面板 ═══ */
.task-popper {
  background: #0f1a1f !important;
  border: 1px solid rgba(0, 255, 255, 0.2) !important;
  border-radius: 12px !important;
  box-shadow: 0 8px 32px rgba(0,0,0,0.7) !important;
  --el-select-option-hover-bg: rgba(18, 26, 36, 0.7);
  --el-select-option-hover-bg-color: rgba(18, 26, 36, 0.7);
  --el-select-option-selected-bg: rgba(15, 22, 32, 0.95);
  --el-select-option-selected-bg-color: rgba(15, 22, 32, 0.95);
  --el-select-option-selected-text-color: #e2e8f0;
  --el-fill-color-blank: transparent;
}
.task-popper .el-select-dropdown__item {
  color: #94a3b8 !important;
  font-size: 13px;
  padding: 9px 18px;
  transition: all 0.15s;
}
.task-popper .el-select-dropdown__item.hover,
.task-popper .el-select-dropdown__item.is-hovering,
.task-popper .el-select-dropdown__item:hover {
  background: rgba(18, 26, 36, 0.7) !important;
  color: #e2e8f0 !important;
}
.task-popper .el-select-dropdown__item.selected {
  color: #22d3ee !important;
  font-weight: 700;
  background: rgba(34, 211, 238, 0.10) !important;
}
.task-popper .el-select-dropdown__item.selected::after { display: none; }
.task-popper .opt-dot {
  display: inline-block; width: 8px; height: 8px; border-radius: 50%;
  margin-right: 8px; vertical-align: middle;
}
.task-popper .el-popper__arrow::before {
  background: #0f1a1f !important;
  border-color: rgba(0, 255, 255, 0.2) !important;
}
</style>
