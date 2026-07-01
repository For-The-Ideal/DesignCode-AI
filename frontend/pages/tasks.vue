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
        :isLoading="loading"
        :isEmptyFiltered="!!filterForm.search || !!filterForm.framework || !!filterForm.platform || !!filterForm.status"
        @detail="goToDetail"
        @delete="handleDelete"
      />
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import AppSidebar from '~/components/layout/AppSidebar.vue'
import TasksFilter from '~/components/tasks/TasksFilter.vue'
import TasksTable from '~/components/tasks/TasksTable.vue'
import TasksStatusOverview from '~/components/tasks/TasksStatusOverview.vue'
import { useMenuListStore } from '~/stores/menuList'
import { taskDisplayName } from '~/utils/taskHelpers'
import { storeToRefs } from 'pinia'
import { useUserStore } from '~/stores/user'
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
const route = useRoute()
const menuStore = useMenuListStore()
const navItems = menuStore.navItems.map(item => ({
  ...item,
  active: route.path === item.to,
}))
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
  ).then(() => {
    const idx = tasks.value.findIndex(t => t.id === row.id)
    if (idx !== -1) tasks.value.splice(idx, 1)
    ElMessage.success('任务已删除')
  }).catch(() => {})
}

// ═══ Mock 数据 ═══
const tasks = ref([
  {
    id: 'T-20240601-001', images: [{ desc: '电商首页', url: '' }],
    framework: 'Flutter', platform: 'mobile', status: 'success', progress: 100,
    duration: '2分38秒', createdAt: '2024-06-01 15:36',
    options: ['响应式', '组件化'],
    code: `import 'package:flutter/material.dart';\n\nvoid main() => runApp(const ECommerceApp());\n\nclass ECommerceApp extends StatelessWidget {\n  const ECommerceApp({super.key});\n  @override\n  Widget build(BuildContext context) {\n    return MaterialApp(\n      title: 'ShopNow',\n      theme: ThemeData(\n        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF6C63FF)),\n        useMaterial3: true,\n      ),\n      home: const HomePage(),\n    );\n  }\n}`,
    steps: [
      { title: '上传设计稿', time: '15:36:20', completed: true },
      { title: 'AI 视觉分析', time: '15:36:28', completed: true },
      { title: '生成设计结构', time: '15:36:47', completed: true },
      { title: '生成代码', time: '15:37:13', completed: true },
      { title: '渲染预览', time: '15:38:58', completed: true },
      { title: '任务完成', time: '15:39:02', completed: true },
    ],
  },
  {
    id: 'T-20240601-002', images: [{ desc: '登录注册页', url: '' }],
    framework: 'React', platform: 'desktop', status: 'running', progress: 68,
    duration: '—', createdAt: '2024-06-01 15:40',
    options: ['响应式'], code: '',
    steps: [
      { title: '上传设计稿', time: '15:40:02', completed: true },
      { title: 'AI 视觉分析', time: '15:40:18', completed: true },
      { title: '生成设计结构', time: '15:40:35', completed: true },
      { title: '生成代码', time: '15:41:10', completed: false },
      { title: '渲染预览', time: '—', completed: false },
      { title: '任务完成', time: '—', completed: false },
    ],
  },
  {
    id: 'T-20240601-003', images: [{ desc: '商品详情', url: '' }],
    framework: 'Vue3', platform: 'mobile', status: 'failed', progress: 45,
    duration: '1分20秒', createdAt: '2024-06-01 15:38',
    error: 'AI 视觉分析超时，请重试', code: '',
    steps: [
      { title: '上传设计稿', time: '15:38:02', completed: true },
      { title: 'AI 视觉分析', time: '15:38:15', completed: true },
      { title: '生成设计结构', time: '15:38:40', completed: false },
      { title: '生成代码', time: '—', completed: false },
      { title: '渲染预览', time: '—', completed: false },
      { title: '任务完成', time: '—', completed: false },
    ],
  },
  {
    id: 'T-20240602-004', images: [{ desc: '个人中心', url: '' }],
    framework: 'Flutter', platform: 'mobile', status: 'success', progress: 100,
    duration: '3分12秒', createdAt: '2024-06-02 09:15',
    options: ['响应式', 'Material 3'],
    code: `class ProfilePage extends StatelessWidget {\n  const ProfilePage({super.key});\n  @override\n  Widget build(BuildContext context) {\n    return Scaffold(\n      backgroundColor: const Color(0xFFF5F5F7),\n      appBar: AppBar(title: const Text('个人中心')),\n      body: Column(children: [\n        Container(padding: EdgeInsets.all(32), child: Column(children: [\n          CircleAvatar(radius: 40, child: Icon(Icons.person, size: 40)),\n          SizedBox(height: 12),\n          Text('用户昵称', style: TextStyle(fontSize: 18)),\n        ])),\n      ]),\n    );\n  }\n}`,
    steps: [
      { title: '上传设计稿', time: '09:15:10', completed: true },
      { title: 'AI 视觉分析', time: '09:15:22', completed: true },
      { title: '生成设计结构', time: '09:15:44', completed: true },
      { title: '生成代码', time: '09:16:30', completed: true },
      { title: '渲染预览', time: '09:18:15', completed: true },
      { title: '任务完成', time: '09:18:27', completed: true },
    ],
  },
  {
    id: 'T-20240602-005', images: [{ desc: '数据仪表盘', url: '' }],
    framework: 'React', platform: 'desktop', status: 'success', progress: 100,
    duration: '4分05秒', createdAt: '2024-06-02 10:22',
    options: ['图表', '深色模式'],
    code: `export default function Dashboard() {\n  const [activeTab, setActiveTab] = useState('overview');\n  return (\n    <div className="dashboard">\n      <header className="dash-header">\n        <h1>数据仪表盘</h1>\n        <div className="header-stats">\n          <div className="stat-card"><span>¥128,430</span><span>总营收</span></div>\n          <div className="stat-card"><span>2,847</span><span>总订单</span></div>\n        </div>\n      </header>\n    </div>\n  );\n}`,
    steps: [
      { title: '上传设计稿', time: '10:22:05', completed: true },
      { title: 'AI 视觉分析', time: '10:22:18', completed: true },
      { title: '生成设计结构', time: '10:22:40', completed: true },
      { title: '生成代码', time: '10:23:50', completed: true },
      { title: '渲染预览', time: '10:26:05', completed: true },
      { title: '任务完成', time: '10:26:27', completed: true },
    ],
  },
  {
    id: 'T-20240602-006', images: [{ desc: '设置页面', url: '' }],
    framework: 'Vue3', platform: 'tablet', status: 'pending', progress: 0,
    duration: '—', createdAt: '2024-06-02 11:08', code: '',
    steps: [
      { title: '上传设计稿', time: '11:08:10', completed: true },
      { title: 'AI 视觉分析', time: '—', completed: false },
      { title: '生成设计结构', time: '—', completed: false },
      { title: '生成代码', time: '—', completed: false },
      { title: '渲染预览', time: '—', completed: false },
      { title: '任务完成', time: '—', completed: false },
    ],
  },
  {
    id: 'T-20240603-007', images: [{ desc: '聊天界面', url: '' }],
    framework: 'Flutter', platform: 'mobile', status: 'running', progress: 32,
    duration: '—', createdAt: '2024-06-03 14:30', code: '',
    steps: [
      { title: '上传设计稿', time: '14:30:02', completed: true },
      { title: 'AI 视觉分析', time: '14:30:15', completed: true },
      { title: '生成设计结构', time: '14:30:40', completed: false },
      { title: '生成代码', time: '—', completed: false },
      { title: '渲染预览', time: '—', completed: false },
      { title: '任务完成', time: '—', completed: false },
    ],
  },
  {
    id: 'T-20240603-008', images: [{ desc: '订单列表', url: '' }],
    framework: 'React', platform: 'mobile', status: 'success', progress: 100,
    duration: '2分55秒', createdAt: '2024-06-03 15:10',
    options: ['筛选', '状态标签'],
    code: `export default function OrderList() {\n  const [filter, setFilter] = useState('all');\n  const orders = [\n    { id: 'ORD001', product: '无线蓝牙耳机', price: '¥299', status: '已发货' },\n    { id: 'ORD002', product: '机械键盘 K8', price: '¥599', status: '配送中' },\n  ];\n  return (<div className="order-page">...</div>);\n}`,
    steps: [
      { title: '上传设计稿', time: '15:10:05', completed: true },
      { title: 'AI 视觉分析', time: '15:10:16', completed: true },
      { title: '生成设计结构', time: '15:10:38', completed: true },
      { title: '生成代码', time: '15:11:50', completed: true },
      { title: '渲染预览', time: '15:13:00', completed: true },
      { title: '任务完成', time: '15:13:05', completed: true },
    ],
  },
])



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
