<template>
  <div class="grid-wrap">
    <!-- ═══ 卡片网格 ═══ -->
    <div class="tasks-grid" v-if="paginatedTasks.length">
      <article
        v-for="row in paginatedTasks"
        :key="row.id"
        class="t-card"
        :class="'st-' + row.status"
        @click="openDetail(row)"
      >
        <!-- 手机壳缩略图 -->
        <div class="t-thumb">
          <el-image v-if="row.images?.[0]?.url" :src="row.images[0].url" fit="cover" class="thumb-img" loading="lazy">
            <template #error><div class="thumb-ph"><i class="fas fa-image"></i></div></template>
          </el-image>
          <div v-else class="thumb-ph"><i class="fas fa-image"></i></div>

          <!-- 失败：hover 出现禁止蒙层（仅盖缩略图，不阻断删除） -->
          <div class="t-block" v-if="row.status === 'failed'">
            <i class="fas fa-ban"></i>
          </div>
        </div>

        <!-- 信息区 -->
        <div class="t-body">
          <!-- 标题行：状态点 + 名称 + 框架药丸 -->
          <div class="t-head">
            <div class="t-title-wrap">
              <span class="t-dot" :title="statusLabel(row.status)"></span>
              <h3 class="t-name" :title="taskName(row)">{{ taskName(row) }}</h3>
            </div>
            <span class="st-pill" :class="'sp-' + row.status">{{ statusLabel(row.status) }}</span>
          </div>

          <!-- 框架 chip -->
          <span class="fw-chip">{{ row.framework }}</span>

          <!-- 其他生成选项标签（逐项拆分，label 来自配置字典） -->
          <div class="t-tags" v-if="configStore.loaded && rowTags(row).length">
            <span class="t-tag" v-for="(tag, i) in rowTags(row)" :key="i">{{ tag }}</span>
          </div>

          <!-- 平台 · 时间 -->
          <div class="t-meta">
            <i :class="platformIcon(row.platform)"></i>
            <span>{{ platformLabel(row.platform) }}</span>
            <span class="dot-sep">·</span>
            <span class="m-time">{{ row.created_at }}</span>
          </div>

          <!-- 进度 -->
          <div class="t-prog" :class="'pb-' + row.status">
            <div class="pbar"><span class="pbar-fill" :style="{ width: row.progress + '%' }"></span></div>
            <div class="prow">
              <span class="p-step">{{ row.status === 'running' && row.current_step ? stepLabel(row.current_step) : '' }}</span>
              <span class="p-num">{{ row.progress }}%</span>
            </div>
          </div>

          <!-- 操作 -->
          <div class="t-acts" @click.stop>
            <button class="icon-btn" title="查看详情" @click="openDetail(row)">
              <i class="fas fa-arrow-up-right-from-square"></i>
            </button>
            <button class="icon-btn icon-btn--danger" title="删除" @click="$emit('delete', row)">
              <i class="fas fa-trash"></i>
            </button>
          </div>
        </div>
      </article>
    </div>

    <!-- ═══ 空状态 ═══ -->
    <div class="empty-place" v-else>
      <el-empty :imageSize="80" :description="isEmptyFiltered ? '未找到匹配的任务' : '暂无任务，快去创建一个吧'" />
    </div>

    <!-- ═══ 分页 ═══ -->
    <div class="pager-wrap" v-if="tasks.length > pageSize">
      <el-pagination
        v-model:currentPage="page"
        :pageSize="pageSize"
        :total="tasks.length"
        layout="prev, pager, next"
        background
        size="small"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useConfigStore } from '~/stores/config'

const props = defineProps({
  tasks: { type: Array, default: () => [] },
  isEmptyFiltered: { type: Boolean, default: false },
})

const emit = defineEmits(['detail', 'delete'])

const page = ref(1)
const pageSize = ref(8)

// 筛选结果变化时回到第一页
watch(() => props.tasks, () => { page.value = 1 })

const paginatedTasks = computed(() => {
  const s = (page.value - 1) * pageSize.value
  return props.tasks.slice(s, s + pageSize.value)
})

// 当前步骤 → 中文简述（与后端 Step* 常量对齐）
const STEP_LABELS = {
  PrepareImages: '准备图片', VisionAnalyzeSkill: '视觉解析', GenerateSkill: '生成代码',
  CodeToHtmlSkill: '转换预览', SaveResult: '保存结果', Done: '完成',
}
const stepLabel = s => STEP_LABELS[s] || s

// 进入详情：生成中/失败 不可查看，给出提示
const openDetail = (row) => {
  if (row.status === 'running') { ElMessage.info('任务正在生成中，暂不可查看详情'); return }
  if (row.status === 'failed') { ElMessage.warning('任务生成失败，无法查看详情'); return }
  emit('detail', row)
}

// 生成选项/高级选项：数组里每一项拆成独立标签；label 唯一来源为配置字典（共享 store，不写死）
const configStore = useConfigStore()
onMounted(() => configStore.load())

// 逐项拆分：options 每项一个标签（component 用具体组件库名），advanced 每项一个标签
const rowTags = (row) => {
  const t = []
  ;(row.options || []).forEach(v => {
    t.push(v === 'component' ? (row.component_lib || configStore.optionLabel(v)) : configStore.optionLabel(v))
  })
  ;(row.advanced || []).forEach(v => t.push(configStore.advancedLabel(v)))
  return t
}

// 模板保持简洁：字典回显统一封装为短方法（数据源仍是配置字典）
const taskName      = (row) => configStore.taskName(row)
const statusLabel   = (s) => configStore.statusLabel(s)
const platformLabel = (p) => configStore.platformLabel(p)
const platformIcon  = (p) => configStore.platformIcon(p)
</script>

<style scoped>
.grid-wrap { flex: 1; min-height: 0; display: flex; flex-direction: column; }

.tasks-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);   /* 一行固定 4 个 */
  gap: 14px;
  padding-bottom: 8px;
}

/* 窄屏时逐级减少每行列数，保证卡片不被压变形 */
@media (max-width: 1180px) { .tasks-grid { grid-template-columns: repeat(3, 1fr); } }
@media (max-width: 820px)  { .tasks-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 560px)  { .tasks-grid { grid-template-columns: 1fr; } }

/* ═══ 卡片（横版）═══ */
.t-card {
  display: flex;
  align-items: stretch;
  gap: 12px;
  padding: 12px;
  border-radius: 16px;
  border: 1px solid rgba(0, 255, 255, 0.08);
  background: rgba(10, 14, 23, 0.65);
  backdrop-filter: blur(12px);
  cursor: pointer;
  transition: transform 0.22s ease, border-color 0.22s ease, box-shadow 0.22s ease;
}
.t-card:hover {
  transform: translateY(-4px);
  border-color: rgba(0, 255, 255, 0.28);
  box-shadow: 0 14px 36px rgba(0, 0, 0, 0.5);
}
.t-card.st-running {
  border-color: rgba(34, 211, 238, 0.32);
  box-shadow: 0 0 22px rgba(34, 211, 238, 0.10);
}
.t-card.st-failed { cursor: not-allowed; }

/* ═══ 手机壳缩略图 ═══ */
.t-thumb {
  position: relative;
  flex-shrink: 0;
  width: 84px;
  aspect-ratio: 9 / 16;   /* 固定手机屏比例，所有缩略图尺寸完全一致 */
  align-self: flex-start;
  border-radius: 12px;
  border: 2px solid rgba(148, 163, 184, 0.22);
  background: #0b1020;
  overflow: hidden;
}
.thumb-img { width: 100%; height: 100%; display: block; }
.t-thumb :deep(.el-image__inner) {
  width: 100%; height: 100%; object-fit: cover; display: block;
}
.thumb-ph {
  width: 100%; height: 100%;
  display: flex; align-items: center; justify-content: center;
}
.thumb-ph i { font-size: 24px; color: rgba(255, 255, 255, 0.12); }

/* 失败禁止蒙层（hover 卡片时淡入） */
.t-block {
  position: absolute; inset: 0; z-index: 2;
  display: flex; align-items: center; justify-content: center;
  background: rgba(2, 6, 12, 0.62);
  backdrop-filter: blur(2px);
  color: #f87171; font-size: 24px;
  opacity: 0; transition: opacity 0.2s ease;
  pointer-events: none;
}
.t-card.st-failed:hover .t-block { opacity: 1; }

/* ═══ 信息区 ═══ */
.t-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.t-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; }
.t-title-wrap { display: flex; align-items: center; gap: 10px; min-width: 0; }
.t-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; background: #64748b; }
.st-success .t-dot { background: #34d399; box-shadow: 0 0 8px rgba(52, 211, 153, 0.6); }
.st-running .t-dot { background: #22d3ee; box-shadow: 0 0 8px rgba(34, 211, 238, 0.6); animation: blink 1.2s infinite; }
.st-pending .t-dot { background: #60a5fa; }
.st-failed .t-dot { background: #f87171; }

.t-name {
  margin: 0; font-size: 16px; font-weight: 800; color: #ffffff;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}

/* 状态药丸（右上）：按状态配色 */
.st-pill {
  flex-shrink: 0;
  font-size: 12px; font-weight: 700;
  padding: 4px 12px; border-radius: 999px;
  border: 1px solid transparent;
}
.sp-success { color: #34d399; background: rgba(52, 211, 153, 0.16); border-color: rgba(52, 211, 153, 0.4); }
.sp-running { color: #22d3ee; background: rgba(34, 211, 238, 0.16); border-color: rgba(34, 211, 238, 0.4); }
.sp-pending { color: #60a5fa; background: rgba(96, 165, 250, 0.16); border-color: rgba(96, 165, 250, 0.4); }
.sp-failed  { color: #f87171; background: rgba(248, 113, 113, 0.16); border-color: rgba(248, 113, 113, 0.4); }
/* 框架 chip（标题下） */
.fw-chip {
  align-self: flex-start;
  font-size: 11px; font-weight: 600; color: #bfdbfe;
  padding: 3px 10px; border-radius: 8px;
  background: rgba(59, 130, 246, 0.16);
  border: 1px solid rgba(59, 130, 246, 0.32);
}

/* 其他生成选项标签 */
.t-tags { display: flex; flex-wrap: wrap; gap: 6px; }
.t-tag {
  font-size: 11px; color: #94a3b8;
  padding: 2px 8px; border-radius: 6px;
  background: rgba(148, 163, 184, 0.10);
  border: 1px solid rgba(148, 163, 184, 0.16);
}

/* 平台 · 时间 */
.t-meta { display: flex; align-items: center; gap: 6px; font-size: 12px; color: #93c5fd; }
.t-meta i { font-size: 12px; opacity: 0.85; }
.dot-sep { color: #475569; }
.m-time { color: #94a3b8; font-family: 'Fira Code', monospace; font-size: 12px; }

/* ═══ 进度 ═══ */
.t-prog { display: flex; flex-direction: column; gap: 6px; }
.pbar {
  height: 8px; border-radius: 999px;
  background: rgba(148, 163, 184, 0.14);
  overflow: hidden;
}
.pbar-fill { display: block; height: 100%; border-radius: inherit; transition: width 0.5s ease; }
.pb-success .pbar-fill { background: linear-gradient(90deg, #34d399, #10b981); }
.pb-failed  .pbar-fill { background: linear-gradient(90deg, #f87171, #ef4444); }
.pb-pending .pbar-fill { background: linear-gradient(90deg, #60a5fa, #3b82f6); }
.pb-running .pbar-fill {
  background: linear-gradient(90deg, #22d3ee, #06b6d4);
  position: relative; overflow: hidden;
}
.pb-running .pbar-fill::after {
  content: ''; position: absolute; inset: 0;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.35), transparent);
  transform: translateX(-100%);
  animation: shimmer 1.4s infinite;
}
.prow { display: flex; align-items: center; justify-content: space-between; }
.p-step { font-size: 12px; color: #22d3ee; opacity: 0.9; }
.p-num { font-size: 13px; font-weight: 800; color: #e2e8f0; font-variant-numeric: tabular-nums; }

/* ═══ 操作（右下圆形按钮）═══ */
.t-acts { display: flex; justify-content: flex-end; gap: 10px; margin-top: auto; padding-top: 6px; }
.icon-btn {
  width: 34px; height: 34px; border-radius: 50%;
  display: inline-flex; align-items: center; justify-content: center;
  cursor: pointer; font-size: 14px;
  border: 1px solid rgba(148, 163, 184, 0.25);
  background: rgba(148, 163, 184, 0.08);
  color: #cbd5e1;
  transition: all 0.18s ease;
}
.icon-btn:hover {
  border-color: rgba(34, 211, 238, 0.5);
  color: #22d3ee;
  background: rgba(34, 211, 238, 0.12);
}
.icon-btn--danger {
  color: #f87171;
  background: rgba(248, 113, 113, 0.14);
  border-color: rgba(248, 113, 113, 0.45);
}
.icon-btn--danger:hover {
  color: #ffffff;
  background: rgba(248, 113, 113, 0.85);
  border-color: rgba(248, 113, 113, 0.9);
}

/* ═══ 空状态 / 分页 ═══ */
.empty-place { flex: 1; display: flex; align-items: center; justify-content: center; min-height: 300px; }
.pager-wrap { display: flex; justify-content: center; margin-top: 18px; flex-shrink: 0; }

@keyframes blink { 0%, 100% { opacity: 1; } 50% { opacity: 0.25; } }
@keyframes shimmer { to { transform: translateX(100%); } }
</style>
