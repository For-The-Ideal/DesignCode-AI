<template>
  <div class="detail-page">
    <AppSidebar
      v-model="sidebarOpen"
      brand="任务中心"
      :navItems="navItems"
      @nav-click="handleNavClick"
    >
      <TasksSidebar
        :tasks="allTasks"
        v-model:expandedId="sidebarExpandedId"
        :sidebarOpen="true"
        :showTitle="false"
        :showDetail="false"
        :collapsible="false"
        @navigate="switchTask"
      />
    </AppSidebar>
    
    <main class="detail-main">

      <!-- ═══ 顶部操作栏 ═══ -->
      <div class="top-bar">
        <button class="btn-back" @click="goBack">
          <i class="fas fa-arrow-left"></i>
          <span>返回列表</span>
        </button>
      </div>

      <!-- ═══ 未找到 ═══ -->
      <div v-if="!task" class="not-found">
        <div class="nf-icon"><i class="fas fa-ghost"></i></div>
        <h2>任务不存在</h2>
        <p>未找到 ID 为 "{{ route.params.id }}" 的任务</p>
        <button class="btn-back-lg" @click="goBack">
          <i class="fas fa-arrow-left"></i> 返回任务列表
        </button>
      </div>

      <template v-if="task">
        <!-- ═══ 内容区：预览 + 代码并排 ═══ -->
        <div class="content-row">
          <div class="content-left">
            <div class="pv-card">
             
              <div class="pv-body">
                <PreviewTemplate
                  v-if="taskDetail?.previewHtml"
                  :html="taskDetail.previewHtml"
                  :showBottomNav="task.platform === 'mobile'"
                />
                <div v-else class="pv-empty">
                  <i class="fas fa-eye-slash"></i>
                  <span>暂无预览数据</span>
                </div>
              </div>
            </div>
          </div>
          <div class="content-right">
            <div class="code-card">
             
              <div class="code-body">
                <CodeEditor
                  v-if="taskDetail"
                  :modelValue="taskDetail.code"
                  :language="codeLang(task.framework)"
                  :readonly="true"
                  height="100%"
                />
                <div v-else class="code-empty">
                  <i class="fas fa-code"></i>
                  <span>暂无代码数据</span>
                </div>
              </div>
            </div>
          </div>
        </div>

      </template>
    </main>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter, onBeforeRouteUpdate } from 'vue-router'
import AppSidebar from '~/components/layout/AppSidebar.vue'
import TasksSidebar from '~/components/tasks/TasksSidebar.vue'
import CodeEditor from '~/components/code/CodeEditor.vue'
import PreviewTemplate from '~/components/previewTempLate/index.vue'

definePageMeta({ key: route => route.fullPath })

const route = useRoute()
const router = useRouter()
const sidebarOpen = ref(true)
const sidebarExpandedId = ref(null)

const switchTask = (id) => router.push(`/detail/${id}`)

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

const goBack = () => router.push('/tasks')

// ═══ 工具函数 ═══
const platformLabel = (p) => ({ mobile: '手机', desktop: '桌面', tablet: '平板' }[p] || p)
const platformIcon = (p) => ({
  mobile: 'fas fa-mobile-alt', desktop: 'fas fa-desktop', tablet: 'fas fa-tablet-alt',
}[p] || 'fas fa-question-circle')
const statusLabel = (s) => ({ success: '已完成', running: '运行中', pending: '等待中', failed: '失败' }[s] || s)
const statusIcon = (s) => ({
  success: 'fas fa-check-circle', running: 'fas fa-spinner fa-pulse', pending: 'fas fa-clock', failed: 'fas fa-exclamation-circle',
}[s] || 'fas fa-question-circle')

const codeLang = (fw) => ({ Flutter: 'dart', React: 'typescript', Vue3: 'html', Vue: 'html' }[fw] || 'dart')

// ═══ 任务查找（watch + onBeforeRouteUpdate 双保险）
const task = ref(null)
const taskDetail = ref(null)

function loadTask(id) {
  if (!id) { task.value = null; taskDetail.value = null; return }
  const found = allTasks.find(t => t.id === id) || null
  task.value = found
  taskDetail.value = found ? (detailMap[found.id] || null) : null
}

// 1) watch 路由参数变化
watch(() => route.params.id, (id) => {
  loadTask(id)
})

// 2) onBeforeRouteUpdate 兜底（组件复用时）
onBeforeRouteUpdate((to) => {
  loadTask(to.params.id)
})

const copyCode = () => {
  if (!taskDetail.value?.code) return
  navigator.clipboard.writeText(taskDetail.value.code)
}

const downloadCode = () => {
  if (!taskDetail.value?.code) return
  const lang = codeLang(task.value?.framework || 'Flutter')
  const ext = { dart: '.dart', typescript: '.tsx', html: '.vue' }[lang] || '.txt'
  const blob = new Blob([taskDetail.value.code], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `generated${ext}`
  a.click()
  URL.revokeObjectURL(url)
}

// ═══ Mock 任务数据 ═══
const allTasks = [
  { id: 'T-20240601-001', images: [{ desc: '电商首页' }], framework: 'Flutter', platform: 'mobile', status: 'success', progress: 100, duration: '2分38秒', createdAt: '2024-06-01 15:36',
    options: ['响应式', '组件化'],
    steps: [
      { title: '上传设计稿', time: '15:36:20', completed: true, icon: 'fas fa-upload' },
      { title: 'AI 视觉分析', time: '15:36:28', completed: true, icon: 'fas fa-brain' },
      { title: '生成设计结构', time: '15:36:47', completed: true, icon: 'fas fa-layer-group' },
      { title: '生成代码', time: '15:37:13', completed: true, icon: 'fas fa-code' },
      { title: '渲染预览', time: '15:38:58', completed: true, icon: 'fas fa-eye' },
      { title: '任务完成', time: '15:39:02', completed: true, icon: 'fas fa-flag-checkered' },
    ],
  },
  // { id: 'T-20240601-002', images: [{ desc: '登录注册页' }], framework: 'React', platform: 'desktop', status: 'running', progress: 68, duration: '—', createdAt: '2024-06-01 15:40',
  //   options: ['表单验证', 'OAuth'],
  //   steps: [
  //     { title: '上传设计稿', time: '15:40:10', completed: true, icon: 'fas fa-upload' },
  //     { title: 'AI 视觉分析', time: '15:40:18', completed: true, icon: 'fas fa-brain' },
  //     { title: '生成设计结构', time: '15:40:42', completed: true, icon: 'fas fa-layer-group' },
  //     { title: '生成代码', time: '15:41:05', completed: true, icon: 'fas fa-code' },
  //     { title: '渲染预览', time: '15:42:30', completed: false, icon: 'fas fa-eye' },
  //     { title: '任务完成', time: '—', completed: false, icon: 'fas fa-flag-checkered' },
  //   ],
  // },
  // { id: 'T-20240601-003', images: [{ desc: '商品详情' }], framework: 'Vue3', platform: 'mobile', status: 'failed', progress: 45, duration: '1分20秒', createdAt: '2024-06-01 15:38',
  //   options: ['商品卡片', '轮播图'],
  //   steps: [
  //     { title: '上传设计稿', time: '15:38:05', completed: true, icon: 'fas fa-upload' },
  //     { title: 'AI 视觉分析', time: '15:38:12', completed: true, icon: 'fas fa-brain' },
  //     { title: '生成设计结构', time: '15:38:30', completed: true, icon: 'fas fa-layer-group' },
  //     { title: '生成代码', time: '15:38:55', completed: false, icon: 'fas fa-code' },
  //     { title: '渲染预览', time: '—', completed: false, icon: 'fas fa-eye' },
  //     { title: '任务完成', time: '—', completed: false, icon: 'fas fa-flag-checkered' },
  //   ],
  // },
  // { id: 'T-20240602-004', images: [{ desc: '个人中心' }], framework: 'Flutter', platform: 'mobile', status: 'success', progress: 100, duration: '3分12秒', createdAt: '2024-06-02 09:15',
  //   options: ['用户信息', '设置'],
  //   steps: [
  //     { title: '上传设计稿', time: '09:15:20', completed: true, icon: 'fas fa-upload' },
  //     { title: 'AI 视觉分析', time: '09:15:28', completed: true, icon: 'fas fa-brain' },
  //     { title: '生成设计结构', time: '09:15:50', completed: true, icon: 'fas fa-layer-group' },
  //     { title: '生成代码', time: '09:16:15', completed: true, icon: 'fas fa-code' },
  //     { title: '渲染预览', time: '09:17:48', completed: true, icon: 'fas fa-eye' },
  //     { title: '任务完成', time: '09:18:32', completed: true, icon: 'fas fa-flag-checkered' },
  //   ],
  // },
  // { id: 'T-20240602-005', images: [{ desc: '数据仪表盘' }], framework: 'React', platform: 'desktop', status: 'success', progress: 100, duration: '4分05秒', createdAt: '2024-06-02 10:22',
  //   options: ['图表', '数据表格'],
  //   steps: [
  //     { title: '上传设计稿', time: '10:22:10', completed: true, icon: 'fas fa-upload' },
  //     { title: 'AI 视觉分析', time: '10:22:18', completed: true, icon: 'fas fa-brain' },
  //     { title: '生成设计结构', time: '10:22:45', completed: true, icon: 'fas fa-layer-group' },
  //     { title: '生成代码', time: '10:23:20', completed: true, icon: 'fas fa-code' },
  //     { title: '渲染预览', time: '10:25:10', completed: true, icon: 'fas fa-eye' },
  //     { title: '任务完成', time: '10:26:27', completed: true, icon: 'fas fa-flag-checkered' },
  //   ],
  // },
  // { id: 'T-20240602-006', images: [{ desc: '设置页面' }], framework: 'Vue3', platform: 'tablet', status: 'pending', progress: 0, duration: '—', createdAt: '2024-06-02 11:08', options: [], steps: [] },
  // { id: 'T-20240603-007', images: [{ desc: '聊天界面' }], framework: 'Flutter', platform: 'mobile', status: 'running', progress: 32, duration: '—', createdAt: '2024-06-03 14:30', options: [], steps: [] },
  // { id: 'T-20240603-008', images: [{ desc: '订单列表' }], framework: 'React', platform: 'mobile', status: 'success', progress: 100, duration: '2分55秒', createdAt: '2024-06-03 15:10', options: [], steps: [] },
]

const detailMap = {
  'T-20240601-001': {
    code: `import 'package:flutter/material.dart';\n\nvoid main() { runApp(const ECommerceApp()); }\nclass ECommerceApp extends StatelessWidget {\n  const ECommerceApp({super.key});\n  @override\n  Widget build(BuildContext context) {\n    return MaterialApp(\n      title: 'ShopNow',\n      theme: ThemeData(colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF6C63FF)), useMaterial3: true),\n      home: const HomePage(),\n    );\n  }\n}\nclass HomePage extends StatefulWidget {\n  const HomePage({super.key});\n  @override\n  State<HomePage> createState() => _HomePageState();\n}\nclass _HomePageState extends State<HomePage> {\n  int _selectedIndex = 0;\n  final List<String> _banners = ['夏日大促 · 全场5折', '新品首发 · 限时特惠', '会员专享 · 积分翻倍'];\n  final List<Map<String, dynamic>> _categories = [\n    {'icon': Icons.checkroom, 'label': '服装'},\n    {'icon': Icons.watch, 'label': '配饰'},\n    {'icon': Icons.devices, 'label': '数码'},\n    {'icon': Icons.chair, 'label': '家居'},\n    {'icon': Icons.sports_soccer, 'label': '运动'},\n    {'icon': Icons.auto_awesome, 'label': '美妆'},\n  ];\n  @override\n  Widget build(BuildContext context) {\n    return Scaffold(\n      appBar: AppBar(title: const Text('ShopNow'), actions: [\n        IconButton(icon: const Icon(Icons.search), onPressed: () {}),\n        IconButton(icon: const Icon(Icons.shopping_cart), onPressed: () {}),\n      ]),\n      body: SingleChildScrollView(\n        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [\n          Container(height: 160, margin: const EdgeInsets.all(16),\n            decoration: BoxDecoration(borderRadius: BorderRadius.circular(16),\n              gradient: const LinearGradient(colors: [Color(0xFF6C63FF), Color(0xFFE91E63)])),\n            child: PageView.builder(\n              itemCount: _banners.length,\n              itemBuilder: (_, i) => Center(child: Text(_banners[i], style: const TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.bold))),\n            ),\n          ),\n          Padding(\n            padding: const EdgeInsets.symmetric(horizontal: 16),\n            child: GridView.count(crossAxisCount: 3, shrinkWrap: true, physics: const NeverScrollableScrollPhysics(),\n              mainAxisSpacing: 12, crossAxisSpacing: 12, childAspectRatio: 2.2,\n              children: _categories.map((c) => Container(\n                decoration: BoxDecoration(color: Colors.grey[100], borderRadius: BorderRadius.circular(12)),\n                child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [\n                  Icon(c['icon'] as IconData, color: const Color(0xFF6C63FF)),\n                  const SizedBox(height: 4),\n                  Text(c['label'] as String, style: const TextStyle(fontSize: 13)),\n                ]),\n              )).toList(),\n            ),\n          ),\n        ]),\n      ),\n      bottomNavigationBar: NavigationBar(selectedIndex: _selectedIndex,\n        onDestinationSelected: (i) => setState(() => _selectedIndex = i),\n        destinations: const [\n          NavigationDestination(icon: Icon(Icons.home), label: '首页'),\n          NavigationDestination(icon: Icon(Icons.category), label: '分类'),\n          NavigationDestination(icon: Icon(Icons.favorite), label: '收藏'),\n          NavigationDestination(icon: Icon(Icons.person), label: '我的'),\n        ],\n      ),\n    );\n  }\n}`,
    previewHtml: `<div style="font-family: -apple-system, sans-serif; background: #f5f5f7; min-height: 100%;">\n  <div style="background: linear-gradient(135deg, #6C63FF, #E91E63); color: #fff; text-align: center; padding: 32px 20px; margin: 12px; border-radius: 16px;">\n    <h2 style="margin:0; font-size:18px;">夏日大促 · 全场5折</h2>\n    <p style="margin:8px 0 0; opacity:0.8;">限时抢购，不容错过</p>\n  </div>\n  <div style="display:grid; grid-template-columns:repeat(3,1fr); gap:10px; padding:0 12px;">\n    ${['服装','配饰','数码','家居','运动','美妆'].map(c => '<div style="background:#fff; border-radius:12px; padding:16px 8px; text-align:center; box-shadow:0 1px 3px rgba(0,0,0,0.06);"><span style="font-size:24px;">' + (c==='服装'?'👗':c==='配饰'?'💍':c==='数码'?'📱':c==='家居'?'🛋️':c==='运动'?'⚽':'💄') + '</span><p style="margin:6px 0 0; font-size:13px; color:#333;">' + c + '</p></div>').join('')}\n  </div>\n</div>`,
  },
  'T-20240602-004': {
    code: `import 'package:flutter/material.dart';\nclass ProfilePage extends StatelessWidget {\n  const ProfilePage({super.key});\n  @override\n  Widget build(BuildContext context) {\n    return Scaffold(\n      backgroundColor: const Color(0xFFF5F5F7),\n      appBar: AppBar(title: const Text('个人中心'), centerTitle: true, backgroundColor: Colors.white),\n      body: SingleChildScrollView(\n        child: Column(children: [\n          Container(padding: const EdgeInsets.symmetric(vertical: 32), color: Colors.white,\n            child: Column(children: [\n              CircleAvatar(radius: 40, backgroundColor: const Color(0xFF6C63FF).withOpacity(0.1),\n                child: const Icon(Icons.person, size: 40, color: Color(0xFF6C63FF))),\n              const SizedBox(height: 12),\n              const Text('用户昵称', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600)),\n              const SizedBox(height: 4),\n              Text('ID: 20240602', style: TextStyle(fontSize: 13, color: Colors.grey[500])),\n            ]),\n          ),\n          const SizedBox(height: 12),\n          Container(color: Colors.white, child: Column(children: [\n            _menuItem(Icons.receipt_long, '我的订单'),\n            _menuItem(Icons.favorite_border, '我的收藏'),\n            _menuItem(Icons.location_on_outlined, '收货地址'),\n            _menuItem(Icons.payment, '支付设置'),\n          ])),\n          const SizedBox(height: 12),\n          Container(color: Colors.white, child: Column(children: [\n            _menuItem(Icons.settings_outlined, '设置'),\n            _menuItem(Icons.help_outline, '帮助与反馈'),\n          ])),\n        ]),\n      ),\n    );\n  }\n  Widget _menuItem(IconData icon, String title) {\n    return ListTile(leading: Icon(icon, color: Colors.grey[700]), title: Text(title, style: const TextStyle(fontSize: 15)), trailing: Icon(Icons.chevron_right, color: Colors.grey[400]), onTap: () {});\n  }\n}`,
    previewHtml: `<div style="font-family: -apple-system, sans-serif; background: #f5f5f7; min-height: 100%;"><div style="background:#fff; padding:28px 20px; text-align:center;"><div style="width:64px; height:64px; border-radius:32px; background:rgba(108,99,255,0.1); margin:0 auto; display:flex; align-items:center; justify-content:center;"><span style="font-size:28px;">👤</span></div><h3 style="margin:10px 0 4px; font-size:17px;">用户昵称</h3><p style="margin:0; color:#999; font-size:12px;">ID: 20240602</p></div></div>`,
  },
  'T-20240602-005': {
    code: `import React, { useState } from 'react';\nconst data = [{ name: '周一', uv: 4000, pv: 2400 },{ name: '周二', uv: 3000, pv: 1398 },{ name: '周三', uv: 5000, pv: 3800 },{ name: '周四', uv: 2780, pv: 3908 },{ name: '周五', uv: 4890, pv: 4800 },{ name: '周六', uv: 2390, pv: 3800 },{ name: '周日', uv: 3490, pv: 4300 }];\nexport default function Dashboard() {\n  return (\n    <div>\n      <h1>Dashboard</h1>\n      <div style={{display:'flex', gap:12, marginBottom:16}}>\n        <div style={{flex:1, background:'#fff', borderRadius:12, padding:16}}><p style={{margin:0, fontSize:20, fontWeight:700, color:'#6C63FF'}}>¥128,430</p><p style={{margin:'4px 0 0', fontSize:12, color:'#999'}}>总营收</p></div>\n        <div style={{flex:1, background:'#fff', borderRadius:12, padding:16}}><p style={{margin:0, fontSize:20, fontWeight:700, color:'#00C9A7'}}>2,847</p><p style={{margin:'4px 0 0', fontSize:12, color:'#999'}}>总订单</p></div>\n        <div style={{flex:1, background:'#fff', borderRadius:12, padding:16}}><p style={{margin:0, fontSize:20, fontWeight:700, color:'#FF6B6B'}}>98.2%</p><p style={{margin:'4px 0 0', fontSize:12, color:'#999'}}>满意度</p></div>\n      </div>\n    </div>\n  );\n}`,
    previewHtml: `<div style="font-family: -apple-system, sans-serif; background: #f0f2f5; min-height: 100%; padding: 16px;"><h2 style="margin:0 0 16px; font-size:20px;">📊 数据仪表盘</h2><div style="display:flex; gap:12px; margin-bottom:16px;"><div style="flex:1; background:#fff; border-radius:12px; padding:16px; box-shadow:0 1px 3px rgba(0,0,0,0.06);"><p style="margin:0; font-size:20px; font-weight:700; color:#6C63FF;">¥128,430</p><p style="margin:4px 0 0; font-size:12px; color:#999;">总营收</p></div><div style="flex:1; background:#fff; border-radius:12px; padding:16px; box-shadow:0 1px 3px rgba(0,0,0,0.06);"><p style="margin:0; font-size:20px; font-weight:700; color:#00C9A7;">2,847</p><p style="margin:4px 0 0; font-size:12px; color:#999;">总订单</p></div><div style="flex:1; background:#fff; border-radius:12px; padding:16px; box-shadow:0 1px 3px rgba(0,0,0,0.06);"><p style="margin:0; font-size:20px; font-weight:700; color:#FF6B6B;">98.2%</p><p style="margin:4px 0 0; font-size:12px; color:#999;">满意度</p></div></div></div>`,
  },
  'T-20240603-008': {
    code: `import React, { useState } from 'react';\nconst orders = [\n  { id: 'ORD001', product: '无线蓝牙耳机', price: '¥299', status: '已发货', date: '2024-06-03' },\n  { id: 'ORD002', product: '机械键盘 K8', price: '¥599', status: '配送中', date: '2024-06-02' },\n  { id: 'ORD003', product: '4K显示器 27寸', price: '¥2,499', status: '已完成', date: '2024-06-01' },\n  { id: 'ORD004', product: 'Type-C 扩展坞', price: '¥199', status: '待付款', date: '2024-06-03' },\n];\nexport default function OrderList() {\n  const [filter, setFilter] = useState('all');\n  const filteredOrders = filter === 'all' ? orders : orders.filter(o => o.status === filter);\n  return (\n    <div>\n      <h1>我的订单</h1>\n      ${''/* short display */}\n      <div>{filteredOrders.map(o => <div key={o.id} style={{padding:12, borderBottom:'1px solid #eee'}}><p>{o.product}</p><p>{o.price} / {o.status}</p></div>)}</div>\n    </div>\n  );\n}`,
    previewHtml: `<div style="font-family: -apple-system, sans-serif; background: #f5f5f7; min-height: 100%; padding: 16px;"><h2 style="margin:0 0 16px;">📋 订单列表</h2>${['无线蓝牙耳机','机械键盘 K8','4K显示器 27寸','Type-C 扩展坞'].map((p,i) => `<div style="background:#fff; border-radius:12px; padding:12px 16px; margin-bottom:8px; display:flex; justify-content:space-between; align-items:center;"><span>${p}</span><span style="color:#999;">¥${['299','599','2,499','199'][i]}</span></div>`).join('')}</div>`,
  },
}

// 手动加载初始数据（替代 watch immediate，避免 TDZ）
sidebarExpandedId.value = route.params.id
loadTask(route.params.id)
</script>

<style scoped>
/* ═══ 布局 ═══ */
.detail-page {
  display: flex;
  height: 100vh;
  background: #0a0a0f;
  color: #e2e8f0;
  overflow: hidden;
}

.detail-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 4px 28px 28px;
  overflow-y: auto;
  min-width: 0;
}
.detail-main::-webkit-scrollbar { width: 4px; }
.detail-main::-webkit-scrollbar-track { background: transparent; }
.detail-main::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.06); border-radius: 2px; }

/* ═══ 顶部返回 ═══ */
.top-bar { padding: 14px 0 8px; }
.btn-back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 16px;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 8px;
  background: transparent;
  color: rgba(255,255,255,0.35);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}
.btn-back:hover {
  border-color: rgba(0,255,255,0.2);
  color: #e2e8f0;
  background: rgba(255,255,255,0.02);
}

/* ═══ 信息卡片 ═══ */
.info-card {
  background: rgba(15, 22, 32, 0.75);
  border: 1px solid rgba(0,255,255,0.06);
  border-radius: 14px;
  padding: 0;
  margin-bottom: 18px;
  overflow: hidden;
}
.info-head {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px 28px 18px;
}
.info-thumb {
  width: 46px; height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}
.thumb--success { background: rgba(52,211,153,0.1); color: #34d399; }
.thumb--running { background: rgba(0,216,255,0.1); color: #00d8ff; }
.thumb--pending { background: rgba(255,255,255,0.03); color: rgba(255,255,255,0.3); }
.thumb--failed  { background: rgba(255,90,108,0.1); color: #ff5a6c; }

.info-main { flex: 1; min-width: 0; }
.info-title {
  font-size: 22px;
  font-weight: 700;
  color: #f1f5f9;
  margin: 0 0 10px;
  line-height: 1.3;
}
.info-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: rgba(255,255,255,0.3);
}
.fw-tag { font-size: 11px !important; }
.meta-sep { color: rgba(255,255,255,0.1); font-size: 10px; }
.meta-plat i { margin-right: 4px; }
.meta-id {
  font-family: 'Fira Code', monospace;
  font-size: 11px;
  color: rgba(255,255,255,0.2);
  background: rgba(255,255,255,0.03);
  padding: 2px 6px;
  border-radius: 4px;
}

.info-badge { flex-shrink: 0; }
.status-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 14px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}
.slabel--success { background: rgba(52,211,153,0.08); color: #34d399; border: 1px solid rgba(52,211,153,0.15); }
.slabel--running { background: rgba(0,216,255,0.08); color: #00d8ff; border: 1px solid rgba(0,216,255,0.15); }
.slabel--pending { background: rgba(255,255,255,0.02); color: rgba(255,255,255,0.35); border: 1px solid rgba(255,255,255,0.06); }
.slabel--failed  { background: rgba(255,90,108,0.08); color: #ff5a6c; border: 1px solid rgba(255,90,108,0.15); }
.status-dot {
  width: 6px; height: 6px;
  border-radius: 50%;
  display: inline-block;
}
.slabel--success .status-dot { background: #34d399; }
.slabel--running .status-dot { background: #00d8ff; animation: dotPulse 1.5s ease-in-out infinite; }
.slabel--pending .status-dot { background: rgba(255,255,255,0.3); }
.slabel--failed .status-dot  { background: #ff5a6c; }
@keyframes dotPulse { 0%,100%{opacity:1} 50%{opacity:0.3} }

/* 状态栏 */
.info-stats {
  display: flex;
  align-items: center;
  gap: 0;
  padding: 14px 28px;
  border-top: 1px solid rgba(255,255,255,0.04);
  background: rgba(0,0,0,0.12);
}
.istat {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 24px;
}
.istat:first-child { padding-left: 0; }
.istat-lbl {
  font-size: 10px;
  color: rgba(255,255,255,0.25);
  text-transform: uppercase;
  letter-spacing: 0.3px;
  white-space: nowrap;
}
.istat-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 120px;
}
.istat-bar .el-progress { flex: 1; }
:deep(.dp .el-progress-bar__outer) { background: rgba(255,255,255,0.05) !important; }
.istat-pct { font-size: 12px; color: rgba(255,255,255,0.35); width: 28px; text-align: right; flex-shrink: 0; }
.istat-val { font-size: 13px; color: rgba(255,255,255,0.55); font-weight: 500; white-space: nowrap; }
.istat-divider {
  width: 1px; height: 18px;
  background: rgba(255,255,255,0.06);
}

/* ═══ 内容区：预览 + 代码并排 ═══ */
.content-row {
  display: flex;
  gap: 18px;
  flex: 1;
  min-height: 0;
  margin-bottom: 18px;
}
.content-left {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}
.content-right {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

/* 代码面板（在 content-right 内） */
.code-card {
  background: rgba(15, 22, 32, 0.6);
  border: 1px solid rgba(255,255,255,0.04);
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

/* 预览卡片 */
.pv-card {
  height: 100%;
  background: rgba(15, 22, 32, 0.6);
  border: 1px solid rgba(255,255,255,0.04);
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.pv-header {
  padding: 10px 16px;
  background: rgba(0,0,0,0.2);
  border-bottom: 1px solid rgba(255,255,255,0.04);
}
.pv-title-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: rgba(255,255,255,0.35);
}
.pv-title-row i { font-size: 13px; }
.pv-title { font-weight: 600; }
.pv-device {
  margin-left: auto;
  font-size: 10px;
  color: rgba(255,255,255,0.2);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}
.pv-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}
.pv-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: rgba(255,255,255,0.06);
}
.pv-empty i { font-size: 28px; }
.pv-empty span { font-size: 12px; color: rgba(255,255,255,0.15); }

.code-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: rgba(0,0,0,0.25);
  border-bottom: 1px solid rgba(255,255,255,0.04);
}
.code-dots { display: flex; gap: 6px; }
.pdot { width: 9px; height: 9px; border-radius: 50%; }
.pdot-red    { background: #ff5f56; }
.pdot-yellow { background: #ffbd2e; }
.pdot-green  { background: #27c93f; }

.code-lang {
  font-size: 11px;
  color: rgba(255,255,255,0.25);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}
.code-actions { margin-left: auto; display: flex; gap: 4px; }
.code-act {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 6px;
  background: transparent;
  color: rgba(255,255,255,0.3);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}
.code-act:hover {
  background: rgba(255,255,255,0.03);
  border-color: rgba(0,255,255,0.15);
  color: #e2e8f0;
}
.code-body {
  flex: 1;
  min-height: 0;
}
:deep(.code-editor-wrapper) {
  border: none !important;
  border-radius: 0 !important;
}
:deep(.code-editor-wrapper:hover) { box-shadow: none !important; }
.code-empty {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: rgba(255,255,255,0.06);
}
.code-empty i { font-size: 28px; }
.code-empty span { font-size: 12px; color: rgba(255,255,255,0.15); }

/* ═══ Not Found ═══ */
.not-found {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-top: -60px;
}
.nf-icon { font-size: 52px; color: rgba(255,255,255,0.05); margin-bottom: 4px; animation: ghostFloat 3s ease-in-out infinite; }
@keyframes ghostFloat { 0%,100%{transform:translateY(0)} 50%{transform:translateY(-8px)} }
.not-found h2 { font-size: 20px; color: rgba(255,255,255,0.3); margin: 0; }
.not-found p { font-size: 13px; color: rgba(255,255,255,0.15); margin: 0; }
.btn-back-lg {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 20px;
  border: 1px solid rgba(0,255,255,0.12);
  border-radius: 8px;
  background: transparent;
  color: rgba(255,255,255,0.4);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 6px;
  font-family: inherit;
}
.btn-back-lg:hover {
  background: rgba(255,255,255,0.03);
  border-color: rgba(0,255,255,0.25);
  color: #e2e8f0;
}

/* ═══ 响应式 ═══ */
@media (max-width: 1100px) {
  .content-row { flex-direction: column; }
  .detail-main { padding: 4px 16px 16px; }
}
</style>
