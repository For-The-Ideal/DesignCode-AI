<template>
  <div class="tasks-page">
    <AppSidebar
      v-model="sidebarOpen"
      brand="任务中心"
      expandedWidth="225px"
      :navItems="menuStore.navItemsWithActive"
      @navClick="handleNavClick"
    >
      <TasksStatusOverview :statusList="statusList" />
    </AppSidebar>

    <main class="task-main">
      <!-- ═══ 页面头 ═══ -->
      <div class="page-top">
        <div class="top-left">
          <h2 class="page-title">任务列表</h2>
          <el-tag type="info" size="small" round class="count-tag">
            {{ filteredTasks.length }} 个任务
          </el-tag>
        </div>
        <el-button type="primary" :icon="Plus" round @click="handleCreate">
          新建任务
        </el-button>
      </div>

      <!-- ═══ 筛选表单 ═══ -->
      <TasksFilter v-model="filterForm" />

      <!-- ═══ 任务表格 ═══ -->
      <TasksTable
        :tasks="filteredTasks"
        :isEmptyFiltered="!!filterForm.search || !!filterForm.framework || !!filterForm.platform || !!filterForm.status"
        @detail="goToDetail"
        @delete="handleDelete"
      />
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import AppSidebar from '~/components/layout/AppSidebar.vue'
import TasksFilter from '~/components/tasks/TasksFilter.vue'
import TasksTable from '~/components/tasks/TasksTable.vue'
import TasksStatusOverview from '~/components/tasks/TasksStatusOverview.vue'
import { useMenuListStore } from '~/stores/menuList'
import { taskDisplayName } from '~/helpers/index'
import { storeToRefs } from 'pinia'
import { useUserStore } from '~/stores/user'
import { taskApi } from '~/api/task'
const router = useRouter()
const sidebarOpen = ref(true)
const userStore = useUserStore()
const { taskCounts } = storeToRefs(userStore)

// ═══ 状态列表（含 count 字段） ═══
const statusList = computed(() => [
  { key: 'success', label: '已完成' },
  { key: 'running', label: '生成中' },
  { key: 'pending', label: '排队中' },
  { key: 'failed', label: '失败' },
].map(s => ({ ...s, count: taskCounts.value[s.key] ?? 0 }))
)
// ═══ 导航 ═══
const menuStore = useMenuListStore()
const handleNavClick = (item) => {
  if (item.active) return
  router.push(item.to)
}

// ═══ 筛选 ═══
const filterForm = ref({ search: '', framework: '', platform: '', status: '' })
const loading = ref(false)

// ═══ 操作 ═══
const goToDetail = row => router.push(`/detail/${row.id}`)
const handleCreate = () => router.push('/code')

const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定删除任务"${taskDisplayName(row)}"吗？此操作不可恢复。`,
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

onMounted(() => {
  fetchTasks()
})



// ═══ 筛选计算 ═══
const filteredTasks = computed(() => {
  let r = tasks.value
  const f = filterForm.value
  if (f.search) {
    const q = f.search.toLowerCase()
    r = r.filter(t =>
      taskDisplayName(t).toLowerCase().includes(q) ||
      t.id.toLowerCase().includes(q) ||
      (t.framework || '').toLowerCase().includes(q)
    )
  }
  if (f.framework) r = r.filter(t => t.framework === f.framework)
  if (f.platform)  r = r.filter(t => t.platform === f.platform)
  if (f.status)    r = r.filter(t => t.status === f.status)
  return r
})
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
  color: #e2e8f0 !important;
  font-weight: 700;
  background: rgba(15, 22, 32, 0.95) !important;
}
.task-popper .el-select-dropdown__item.selected::after { display: none; }
.task-popper .el-popper__arrow::before {
  background: #0f1a1f !important;
  border-color: rgba(0, 255, 255, 0.2) !important;
}
</style>
