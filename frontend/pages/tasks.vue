<template>
    <div class="tasks-page">
        <!-- ═══ 左侧可折叠菜单 ═══ -->
        <AppSidebar
            v-model="sidebarOpen"
            brand="任务中心"
            :nav-items="navItems"
            @nav-click="handleNavClick"
        >
            <TasksSidebar
                v-model:expanded-id="expandedId"
                :tasks="tasks"
                :sidebar-open="sidebarOpen"
            />
        </AppSidebar>

        <!-- ═══ 右侧主内容区：代码编辑器 + 实时预览 ═══ -->
        <main class="task-main">
            <!-- 未选中任务时的占位提示 -->
            <div v-if="!expandedId" class="main-placeholder">
                <i class="fas fa-code"></i>
                <p>点击左侧任务查看生成的代码与预览</p>
            </div>

            <!-- 选中任务时：左侧编辑器 + 右侧模拟器 -->
            <template v-else>
                <div class="result-editor">
                    <CodeEditor
                        v-model="selectedTask.code"
                        :language="codeLanguage"
                        :readonly="true"
                        height="500px"
                        placeholder="// 选择任务查看生成代码..."
                    />
                </div>
                <div class="result-preview">
                    <FlutterTemplate :html="selectedTask?.preview || ''" :showBottomNav="false" />
                </div>
            </template>
        </main>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import AppSidebar from '~/components/layout/AppSidebar.vue'
import TasksSidebar from '~/components/tasks/TasksSidebar.vue'
import CodeEditor from '~/components/code/CodeEditor.vue'
import FlutterTemplate from '~/components/previewTempLate/index.vue'
import { handleCopy } from '~/utils/index.js'

const router = useRouter()
const sidebarOpen = ref(true)
const expandedId = ref('task_001')

// 当前选中的任务
const selectedTask = computed(() => {
    return tasks.value.find(t => t.id === expandedId.value) || null
})

// 框架 → 语言映射
const langMap = { Flutter: 'dart', React: 'typescript', Vue3: 'html', Vue: 'html' }
const codeLanguage = computed(() => langMap[selectedTask.value?.framework] || 'dart')

const handleCopyCode = () => {
    if (selectedTask.value?.code) {
        handleCopy(selectedTask.value.code)
    }
}

const navItems = [
    { icon: 'fa-solid fa-code', label: '代码生成', active: false, to: '/code' },
    { icon: 'fa-regular fa-copy', label: '模板市场', active: false, to: '/templates' },
    { icon: 'fa-regular fa-folder', label: '任务列表', active: true, to: '/tasks' },
    { icon: 'fa-regular fa-file', label: '我的项目', active: false, to: '/projects' },
]

const handleNavClick = (item) => {
    if (item.active) return
    router.push(item.to)
}

// ═══ Mock 数据（后续接入 API） ═══
const tasks = ref([
    {
        id: 'task_001',
        images: [{ desc: '电商首页', url: 'https://example.com/1.png' }],
        framework: 'Flutter',
        platform: 'mobile',
        status: 'success',
        progress: 100,
        time: '15:39:02',
        createdAt: '2024-05-26 15:36:20',
        duration: '2分38秒',
        error: '',
        code: `import 'package:flutter/material.dart';

void main() => runApp(const ECommerceApp());

class ECommerceApp extends StatelessWidget {
  const ECommerceApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '电商首页',
      theme: ThemeData(primarySwatch: Colors.blue, useMaterial3: true),
      home: const HomePage(),
    );
  }
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _products = [
    {'name': '商品A', 'price': 99.9},
    {'name': '商品B', 'price': 199.9},
    {'name': '商品C', 'price': 299.9},
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('电商首页'),
        actions: [
          IconButton(icon: const Icon(Icons.search), onPressed: () {}),
          IconButton(icon: const Icon(Icons.shopping_cart), onPressed: () {}),
        ],
      ),
      body: Column(
        children: [
          Container(
            height: 180,
            margin: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [Color(0xFF667eea), Color(0xFF764ba2)],
              ),
              borderRadius: BorderRadius.circular(16),
            ),
            child: const Center(
              child: Text('🔥 限时特惠', style: TextStyle(
                color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold)),
            ),
          ),
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: _products.length,
              itemBuilder: (_, i) => Card(
                margin: const EdgeInsets.only(bottom: 12),
                child: ListTile(
                  title: Text(_products[i]['name']),
                  subtitle: Text('¥\${_products[i]['price']}'),
                  trailing: const Icon(Icons.chevron_right),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}`,
        preview: '<div style="padding:20px;font-family:sans-serif;background:#f5f5f5;min-height:100%"><div style="background:linear-gradient(135deg,#667eea,#764ba2);border-radius:16px;padding:30px;color:#fff;text-align:center;margin-bottom:16px"><h2 style="margin:0">🔥 限时特惠</h2><p style="margin:8px 0 0;opacity:.8">全场低至5折</p></div><div style="background:#fff;border-radius:12px;margin-bottom:10px;padding:16px;display:flex;align-items:center;gap:12px"><div style="font-size:36px">🛍️</div><div style="flex:1"><div style="font-weight:600">商品A</div><div style="color:#e74c3c;font-weight:700">¥99.9</div></div><span style="color:#ccc">›</span></div><div style="background:#fff;border-radius:12px;margin-bottom:10px;padding:16px;display:flex;align-items:center;gap:12px"><div style="font-size:36px">👟</div><div style="flex:1"><div style="font-weight:600">商品B</div><div style="color:#e74c3c;font-weight:700">¥199.9</div></div><span style="color:#ccc">›</span></div><div style="background:#fff;border-radius:12px;padding:16px;display:flex;align-items:center;gap:12px"><div style="font-size:36px">📱</div><div style="flex:1"><div style="font-weight:600">商品C</div><div style="color:#e74c3c;font-weight:700">¥299.9</div></div><span style="color:#ccc">›</span></div></div>',
        steps: [
            { title: '上传设计稿', time: '15:36:20', completed: true, icon: 'fa-solid fa-cloud-arrow-up' },
            { title: 'AI 视觉分析', time: '15:36:28', completed: true, icon: 'fa-solid fa-eye' },
            { title: '生成设计结构', time: '15:36:47', completed: true, icon: 'fa-solid fa-layer-group' },
            { title: '生成代码', time: '15:37:13', completed: true, icon: 'fa-solid fa-code' },
            { title: '渲染预览', time: '15:38:58', completed: true, icon: 'fa-solid fa-display' },
            { title: '任务完成', time: '15:39:02', completed: true, icon: 'fa-solid fa-flag-checkered' },
        ],
    },
    {
        id: 'task_002',
        images: [{ desc: '登录注册', url: 'https://example.com/2.png' }],
        framework: 'React',
        platform: 'desktop',
        status: 'pending',
        progress: 0,
        time: '15:38:02',
        createdAt: '2024-05-26 15:38:02',
        duration: '—',
        error: '',
        code: '',
        preview: '',
        steps: [
            { title: '上传设计稿', time: '15:38:02', completed: true, icon: 'fa-solid fa-cloud-arrow-up' },
            { title: 'AI 视觉分析', time: '—', completed: false, icon: 'fa-solid fa-eye' },
            { title: '生成设计结构', time: '—', completed: false, icon: 'fa-solid fa-layer-group' },
            { title: '生成代码', time: '—', completed: false, icon: 'fa-solid fa-code' },
            { title: '渲染预览', time: '—', completed: false, icon: 'fa-solid fa-display' },
            { title: '任务完成', time: '—', completed: false, icon: 'fa-solid fa-flag-checkered' },
        ],
    },
    {
        id: 'task_003',
        images: [{ desc: '商品详情', url: 'https://example.com/3.png' }],
        framework: 'Vue3',
        platform: 'mobile',
        status: 'failed',
        progress: 45,
        time: '15:35:50',
        createdAt: '2024-05-26 15:35:10',
        duration: '1分20秒',
        error: '视觉分析超时，请重试',
        code: '',
        preview: '',
        steps: [
            { title: '上传设计稿', time: '15:35:10', completed: true, icon: 'fa-solid fa-cloud-arrow-up' },
            { title: 'AI 视觉分析', time: '15:35:28', completed: true, icon: 'fa-solid fa-eye' },
            { title: '生成设计结构', time: '15:35:47', completed: false, icon: 'fa-solid fa-layer-group' },
            { title: '生成代码', time: '—', completed: false, icon: 'fa-solid fa-code' },
            { title: '渲染预览', time: '—', completed: false, icon: 'fa-solid fa-display' },
            { title: '任务完成', time: '—', completed: false, icon: 'fa-solid fa-flag-checkered' },
        ],
    },
])
</script>

<style scoped>
.tasks-page {
    display: flex;
    min-height: calc(100vh - 140px);
    background: #0a0a0f;
}

/* ═══ 右侧主内容 ═══ */
.task-main {
    flex: 1;
    display: flex;
    gap: 24px;
    padding: 20px 24px;
    min-height: 0;
}

/* 未选中任务占位 */
.main-placeholder {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: rgba(255, 255, 255, 0.2);
}
.main-placeholder i {
    font-size: 48px;
    margin-bottom: 12px;
    display: block;
}
.main-placeholder p {
    font-size: 14px;
}

/* ── 结果展示：编辑器 + 预览 ── */
.result-editor {
    flex: 1.55;
    min-width: 0;
    display: flex;
    flex-direction: column;
    background: rgba(15, 20, 30, 0.6);
    backdrop-filter: blur(12px);
    border: 1px solid rgba(0, 255, 255, 0.18);
    border-radius: 18px;
    overflow: hidden;
}
.result-preview {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    background: rgba(15, 20, 30, 0.6);
    backdrop-filter: blur(12px);
    border: 1px solid rgba(0, 255, 255, 0.18);
    border-radius: 18px;
    overflow: hidden;
}
.result-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 18px;
    background: rgba(0, 0, 0, 0.3);
    border-bottom: 1px solid rgba(0, 255, 255, 0.1);
    font-size: 14px;
    font-weight: 600;
    color: #00cfff;
    width: 100%;
}
.result-header i { font-size: 15px; margin-right: 6px; }

.act-btn {
    display: flex;
    align-items: center;
    gap: 5px;
    padding: 5px 12px;
    background: rgba(0, 255, 255, 0.07);
    border: 1px solid rgba(0, 255, 255, 0.18);
    border-radius: 18px;
    color: #00cfff;
    font-size: 12px;
    cursor: pointer;
    transition: all 0.2s;
    font-family: inherit;
}
.act-btn i { font-size: 12px; }
.act-btn:hover { background: rgba(0, 255, 255, 0.14); border-color: #00ffff; }

.device-badge {
    font-size: 11px;
    color: #6b7280;
    letter-spacing: 0.3px;
    background: rgba(255,255,255,0.04);
    padding: 3px 10px;
    border-radius: 20px;
}

@media (max-width: 1100px) {
    .task-main { flex-direction: column; }
    .result-editor { height: 480px; }
    .result-preview { flex: unset; }
}
</style>
