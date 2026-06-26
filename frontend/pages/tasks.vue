<template>
  <div class="tasks-page">
    <AppSidebar
      v-model="sidebarOpen"
      brand="任务中心"
      :navItems="navItems"
      @nav-click="handleNavClick"
    >
      <div class="sidebar-slot-content">
        <div class="slot-section-title">
          <i class="fas fa-chart-pie"></i>
          <span>任务概览</span>
        </div>
        <div class="slot-stats">
          <div
            v-for="s in statusList"
            :key="s.key"
            class="slot-stat-row"
          >
            <div class="slot-stat-left">
              <span class="slot-dot" :class="'sdot--' + s.key"></span>
              <span class="slot-stat-label">{{ s.label }}</span>
            </div>
            <span class="slot-stat-num">{{ statusCounts[s.key] }}</span>
          </div>
        </div>
      </div>
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
      <el-form :model="filterForm" inline class="filter-form">
        <el-form-item label="搜索">
          <el-input
            v-model="filterForm.search"
            placeholder="任务名称、框架..."
            :prefixIcon="Search"
            clearable
            class="f-search"
          />
        </el-form-item>
        <el-form-item label="框架">
          <el-select
            v-model="filterForm.framework"
            placeholder="全部"
            clearable
            class="f-sel"
            popper-class="task-popper"
          >
            <el-option label="Flutter" value="Flutter" />
            <el-option label="React" value="React" />
            <el-option label="Vue3" value="Vue3" />
          </el-select>
        </el-form-item>
        <el-form-item label="平台">
          <el-select
            v-model="filterForm.platform"
            placeholder="全部"
            clearable
            class="f-sel"
            popper-class="task-popper"
          >
            <el-option label="📱 手机" value="mobile" />
            <el-option label="🖥️ 桌面" value="desktop" />
            <el-option label="📋 平板" value="tablet" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="filterForm.status"
            placeholder="全部"
            clearable
            class="f-sel"
            popper-class="task-popper"
          >
            <el-option label="✅ 已完成" value="success" />
            <el-option label="🔄 生成中" value="running" />
            <el-option label="⏳ 排队中" value="pending" />
            <el-option label="❌ 失败" value="failed" />
          </el-select>
        </el-form-item>
      </el-form>

      <!-- ═══ 任务表格 ═══ -->
      <div class="table-wrap">
        <el-table
          :data="paginatedTasks"
          row-class-name="task-row"
          highlight-current-row
          style="width:100%"
        >
          <el-table-column label="项目名称" align="left">
            <template #default="{ row }">
              <div class="cell-proj">
                <div class="proj-icon">
                  <el-image
                    v-if="row.images?.[0]?.url"
                    :src="row.images[0].url"
                    fit="cover"
                    class="icon-img"
                  >
                    <template #error><i class="fas fa-image"></i></template>
                  </el-image>
                  <i v-else class="fas fa-image"></i>
                </div>
                <span class="proj-name">{{ taskDisplayName(row) }}</span>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="框架" align="center">
            <template #default="{ row }">
              <el-tag :type="fwTagType(row.framework)" size="small" effect="dark" round>
                {{ row.framework }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="平台" align="center">
            <template #default="{ row }">
              <span class="cell-plat">
                <i :class="platformIcon(row.platform)"></i>
                {{ platformLabel(row.platform) }}
              </span>
            </template>
          </el-table-column>

          <el-table-column label="状态" align="center">
            <template #default="{ row }">
              <el-tag :type="statusTagType(row.status)" size="small" effect="dark" round>
                {{ statusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="进度" align="center">
            <template #default="{ row }">
              <div class="cell-prog">
                <el-progress
                  :percentage="row.progress"
                  :status="row.status === 'success' ? 'success' : row.status === 'failed' ? 'exception' : ''"
                  :strokeWidth="6"
                  :showText="false"
                  :color="progressColor(row.status)"
                  style="flex:1"
                />
                <span class="prog-num">{{ row.progress }}%</span>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="创建时间" align="center">
            <template #default="{ row }">
              <span class="cell-time">{{ row.createdAt }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="120" align="center" fixed="right">
            <template #default="{ row }">
              <div class="cell-acts" @click.stop>
                <el-button link type="primary" size="small" @click="goToDetail(row)">
                  <span>详情</span>
                </el-button>
                <el-button link type="danger" size="small" @click="handleDelete(row)">
                  <span>删除</span>
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 空状态 -->
      <div class="empty-place" v-if="filteredTasks.length === 0">
        <el-empty
          :imageSize="80"
          :description="filterForm.search ? '未找到匹配的任务' : '暂无任务，快去创建一个吧'"
        />
      </div>

      <!-- 分页 -->
      <div class="pager-wrap" v-if="filteredTasks.length > pageSize">
        <el-pagination
          v-model:currentPage="currentPage"
          :pageSize="pageSize"
          :total="filteredTasks.length"
          layout="prev, pager, next"
          background
          small
        />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import AppSidebar from '~/components/layout/AppSidebar.vue'

const router = useRouter()
const sidebarOpen = ref(true)

// ═══ 任务状态统计 ═══
const statusList = [
  { key: 'success', label: '已完成' },
  { key: 'running', label: '生成中' },
  { key: 'pending', label: '排队中' },
  { key: 'failed', label: '失败' },
]
const statusCounts = computed(() => {
  const counts = { success: 0, running: 0, pending: 0, failed: 0 }
  tasks.value.forEach(t => {
    if (counts[t.status] !== undefined) counts[t.status]++
  })
  return counts
})

// ═══ 筛选 ═══
const filterForm = ref({ search: '', framework: '', platform: '', status: '' })

// ═══ 分页 ═══
const currentPage = ref(1)
const pageSize = ref(10)

// ═══ 导航 ═══
const navItems = [
  { icon: 'fa-solid fa-code',       label: '代码生成', active: false, to: '/code' },
  { icon: 'fa-regular fa-copy',     label: '模板市场', active: false, to: '/templates' },
  { icon: 'fa-regular fa-folder',   label: '任务列表', active: true,  to: '/tasks' },
  { icon: 'fa-regular fa-file',     label: '我的项目', active: false, to: '/projects' },
]

const handleNavClick = (item) => {
  if (item.active) return
  router.push(item.to)
}

// ═══ 工具函数 ═══
const PLAT = { mobile: '手机', desktop: '桌面', tablet: '平板' }
const PLAT_ICON = { mobile: 'fas fa-mobile-alt', desktop: 'fas fa-desktop', tablet: 'fas fa-tablet-alt' }
const STATUS_MAP = { pending: '排队中', running: '生成中', success: '已完成', failed: '失败' }

const platformLabel = p => PLAT[p] || p || '—'
const platformIcon = p => PLAT_ICON[p] || 'fas fa-question-circle'
const statusLabel  = s => STATUS_MAP[s] || s

const statusTagType = s =>
  ({ success: 'success', running: 'warning', pending: 'info', failed: 'danger' }[s] || 'info')

const fwTagType = fw =>
  ({ Flutter: 'primary', React: '', Vue3: 'success' }[fw] || 'info')

const progressColor = s =>
  ({ success: '#34d399', failed: '#f87171', running: '#facc15' }[s])

const taskDisplayName = t =>
  t.images?.[0]?.desc || `${t.platform || '未知'} · ${t.framework || '未知'}`


// ═══ 数据 ═══
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

const paginatedTasks = computed(() => {
  const s = (currentPage.value - 1) * pageSize.value
  return filteredTasks.value.slice(s, s + pageSize.value)
})

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
</script>

<style scoped>
/* ═══════════════ 页面布局 ═══════════════ */
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

/* ═══════════════ 页面头 ═══════════════ */
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
.count-tag {
  font-size: 11px;
  opacity: 0.7;
}

/* ═══════════════ 筛选表单 ═══════════════ */
.filter-form {
  flex-shrink: 0;
  margin-bottom: 14px;
  display: flex;
  flex-wrap: wrap;
  gap: 4px 12px;
}
.f-search { width: 220px; }
.f-sel    { width: 130px; }

/* select/input 输入框：统一 dark-slate 风格 */
:deep(.f-sel .el-select__wrapper),
:deep(.f-sel .el-input__wrapper),
:deep(.f-search .el-input__wrapper) {
  background: rgba(15, 22, 32, 0.75) !important;
  border: 1px solid rgba(0,255,255,0.12) !important;
  box-shadow: none !important;
}
:deep(.f-sel .el-input__inner) {
  color: #cbd5e1 !important;
}
:deep(.f-sel .el-select__caret) {
  color: rgba(0,255,255,0.4) !important;
}

/* ═══════════════ 表格容器 ═══════════════ */
.table-wrap {
  flex: 1;
  min-height: 0;
  border-radius: 14px;
  border: 1px solid rgba(0, 255, 255, 0.08);
  overflow: hidden;
  background: rgba(10, 14, 23, 0.65);
  backdrop-filter: blur(12px);
}

/* ═══════════════ 表格：背景强制透明 + 表头暗色 ═══════════════ */
:deep(.el-table),
:deep(.el-table__inner-wrapper),
:deep(.el-table__body-wrapper),
:deep(.el-table__header-wrapper),
:deep(.el-table__fixed),
:deep(.el-table__fixed-right),
:deep(.el-table__fixed-header-wrapper),
:deep(.el-table__fixed-body-wrapper),
:deep(.el-table__fixed-right-patch),
:deep(.el-table__empty-block),
:deep(.el-table__empty-text),
:deep(.el-table__body),
:deep(.el-table__body tr),
:deep(.el-table__body tr td) {
  background: transparent !important;
  border-bottom: none !important;
}

:deep(.el-table__header th),
:deep(.el-table__fixed-header-wrapper th),
:deep(.el-table.is-scrolling-none th.el-table-fixed-column--right) {
  background: radial-gradient(
circle at top,
rgba(41,70,255,.08),
transparent 35%
),
#090B12;
  color: #94a3b8 !important;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.3px;
  border-bottom: 1px solid rgba(0,255,255,0.08) !important;
  padding: 14px 8px;
}

:deep(.task-row),
:deep(.task-row td) {
  background: transparent !important;
  border-bottom: 1px solid rgba(0,255,255,0.04) !important;
  color: #cbd5e1;
  font-size: 13px;
  padding: 12px 0;
}

:deep(.el-table__body tr:hover > td) {
  background: rgba(0, 255, 255, 0.03) !important;
}
:deep(.el-table__body tr.current-row > td) {
  background: rgba(0, 255, 255, 0.05) !important;
}

/* 移除表格底部白线（Element Plus 伪元素边框） */
:deep(.el-table::before),
:deep(.el-table::after),
:deep(.el-table__inner-wrapper::before),
:deep(.el-table__inner-wrapper::after) {
  display: none !important;
}

/* ── 单元格内容 ── */
.cell-proj { display: flex; align-items: center; gap: 10px; }
.proj-icon {
  width: 34px; height: 34px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(135deg, rgba(0,255,255,0.08), rgba(255,0,255,0.04));
  border: 1px solid rgba(0,255,255,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}
.proj-icon i { font-size: 13px; color: rgba(255,255,255,0.15); }
.icon-img { width: 100%; height: 100%; }
.proj-name {
  font-size: 14px;
  font-weight: 500;
  color: #e2e8f0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cell-plat {
  font-size: 12px;
  color: #94a3b8;
  display: flex;
  align-items: center;
  gap: 5px;
  justify-content: center;
}
.cell-plat i { font-size: 11px; opacity: 0.5; }

.cell-prog {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: center;
  flex-wrap: nowrap;
}
.prog-num {
  font-size: 11px;
  color: #94a3b8;
  font-weight: 600;
  flex-shrink: 0;
  white-space: nowrap;
}

.cell-time {
  font-size: 12px;
  color: #94a3b8;
  font-family: 'Fira Code', monospace;
}

.cell-acts {
  display: flex;
  gap: 8px;
  justify-content: center;
  align-items: center;
}

/* ═══════════════ 空状态 ═══════════════ */
.empty-place {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

/* ═══════════════ 分页 ═══════════════ */
.pager-wrap {
  display: flex;
  justify-content: center;
  margin-top: 14px;
  flex-shrink: 0;
}

/* ═══════════════ 响应式 ═══════════════ */
@media (max-width: 1100px) {
  .task-main { padding: 16px; }
  .f-search { width: 100%; }
  .f-sel    { width: 110px; }
}

/* ═══ AppSidebar 插槽：任务统计概览 ═══ */
.sidebar-slot-content { padding: 0 2px; }
.slot-section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  font-weight: 600;
  color: rgba(255,255,255,0.25);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 10px;
}
.slot-section-title i { font-size: 11px; }
.slot-stats {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.slot-stat-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  border-radius: 8px;
  transition: background 0.2s;
}
.slot-stat-row:hover {
  background: rgba(255,255,255,0.03);
}
.slot-stat-left {
  display: flex;
  align-items: center;
  gap: 8px;
}
.slot-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}
.sdot--success { background: #34d399; }
.sdot--running { background: #facc15; }
.sdot--pending { background: #60a5fa; }
.sdot--failed  { background: #f87171; }
.slot-stat-label {
  font-size: 12px;
  color: rgba(255,255,255,0.45);
}
.slot-stat-num {
  font-size: 13px;
  font-weight: 700;
  color: rgba(255,255,255,0.6);
  font-variant-numeric: tabular-nums;
}

</style>

<!-- ═══════════════ 非 scoped：仅覆盖 teleported 下拉面板 ═══════════════ -->
<style>
.task-popper {
  background: #0f1a1f !important;
  border: 1px solid rgba(0, 255, 255, 0.2) !important;
  box-shadow: 0 8px 32px rgba(0,0,0,0.7) !important;

  /* 注入 CSS 变量确保 option 不白 */
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
.task-popper .el-select-dropdown__item.selected::after {
  display: none;
}
.task-popper .el-popper__arrow::before {
  background: #0f1a1f !important;
  border-color: rgba(0, 255, 255, 0.2) !important;
}
</style>
