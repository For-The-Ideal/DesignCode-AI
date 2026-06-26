<template>
  <!-- ═══ 任务列表（展开时显示） ═══ -->
  <Transition name="card-fade">
    <div v-show="sidebarOpen" class="task-list-wrapper">
          <div class="task-list-header" v-if="showTitle">
            <span class="task-count">共 {{ tasks.length }} 个任务</span>
          </div>

          <div
            v-for="task in tasks"
            :key="task.id"
            class="task-group"
          >
            <!-- ── 紧凑条目（始终可见） ── -->
            <div
              class="task-item"
              :class="{ expanded: expandedId === task.id }"
              @click="toggleExpand(task.id)"
            >
              <div class="task-item-top">
                <div class="task-item-title">{{ taskDisplayName(task) }}</div>
                <span class="task-status-badge" :class="'status--' + task.status">
                  <span class="status-dot"></span>
                  {{ statusLabel(task.status) }}
                </span>
              </div>
              <div v-if="task.status !== 'failed'" class="task-progress">
                <div class="task-progress-bar">
                  <div
                    class="task-progress-fill"
                    :class="{ animated: task.status === 'running' }"
                    :style="{ width: task.progress + '%' }"
                  ></div>
                </div>
                <span class="task-progress-num">{{ task.progress }}%</span>
              </div>
              <div v-if="task.status === 'failed'" class="task-error">
                <i class="fas fa-exclamation-circle"></i>
                {{ task.error || '生成失败' }}
              </div>
              <div class="task-item-meta">
                <span>平台：{{ task.framework }}</span>
                <span>创建时间：{{ task.createdAt || task.time }}</span>
              </div>
            </div>

            <!-- ── 展开详情卡片 ── -->
            <Transition name="detail-expand">
              <div v-if="expandedId === task.id" class="task-detail">
                <!-- 头部：缩略图 + 标题/标签 + 时间 -->
                <div class="card-head">
                  <div class="card-thumb">
                    <img
                      v-if="task.images?.[0]?.url"
                      :src="task.images[0].url"
                      @error="e => e.target.style.display = 'none'"
                      alt=""
                    />
                    <i v-if="!task.images?.[0]?.url" class="fa-solid fa-image"></i>
                  </div>
                  <div class="card-head-info">
                    <h3 class="card-title">{{ taskDisplayName(task) }}</h3>
                    <div class="card-tags">
                      <span v-for="(tag, i) in taskTags(task)" :key="i" class="card-tag">{{ tag }}</span>
                    </div>
                  </div>
                  <span class="card-time">{{ task.createdAt || task.time }}</span>
                </div>

                <!-- 进度条 -->
                <div class="progress-block">
                  <span class="progress-label">还原度 {{ task.progress }}%</span>
                  <div class="progress-bar">
                    <div class="progress-fill" :style="{ width: task.progress + '%' }">
                      <div class="progress-shine"></div>
                    </div>
                  </div>
                </div>

                <!-- 时间轴步骤 -->
                <div v-if="task.steps && task.steps.length" class="timeline">
                  <div class="timeline-line"></div>
                  <div
                    v-for="(step, idx) in task.steps"
                    :key="idx"
                    class="timeline-step"
                    :class="stepPhase(task, idx)"
                  >
                    <div class="step-left">
                      <div class="step-indicator" :class="{ done: step.completed }">
                        <i :class="step.icon" class="step-icon"></i>
                      </div>
                      <div class="step-body">
                        <div class="step-title">{{ step.title }}</div>
                        <div class="step-time">{{ step.time }}</div>
                      </div>
                    </div>
                    <div class="step-dot" :class="{ done: step.completed }">
                      <i v-if="step.completed" class="fas fa-check"></i>
                    </div>
                  </div>
                </div>

                <!-- 查看详情按钮 -->
                <button class="view-detail-btn" @click.stop="goToDetail(task.id)" v-if="showDetail">
                  <i class="fas fa-external-link-alt"></i> 查看完整详情
                </button>
              </div>
            </Transition>
          </div>

          <!-- 空态 -->
          <div v-if="tasks.length === 0" class="task-empty">
            <i class="fas fa-inbox"></i>
            <p>暂无进行中的任务</p>
          </div>
    </div>
  </Transition>
</template>

<script setup>
const props = defineProps({
  tasks: { type: Array, required: true },
  expandedId: { type: [String, null], default: null },
  sidebarOpen: { type: Boolean, default: true },
  collapsible: { type: Boolean, default: false },
  showDetail: { type: Boolean, default: true },
  showTitle: { type: Boolean, default: true },
})

const emit = defineEmits(['update:expandedId', 'navigate'])

const toggleExpand = (id) => {
  if (!props.collapsible) return
  emit('update:expandedId', props.expandedId === id ? null : id)
}

const goToDetail = (id) => {
  emit('navigate', id)
}

const STATUS_MAP = {
  pending: '排队中',
  running: '生成中',
  success: '已完成',
  failed: '失败',
}
const statusLabel = (s) => STATUS_MAP[s] || s

// 定位当前进度：返回 'phase-past' | 'phase-current' | 'phase-future'
const stepPhase = (task, idx) => {
  const firstPending = task.steps.findIndex(s => !s.completed)
  const currentIdx = firstPending === -1 ? task.steps.length - 1 : firstPending
  if (idx < currentIdx) return 'phase-past'
  if (idx === currentIdx) return 'phase-current'
  return 'phase-future'
}

// 任务显示标题：优先取第一张图片描述，否则 平台·框架
const taskDisplayName = (task) => {
  return task.images?.[0]?.desc || `${task.platform} · ${task.framework}`
}

// 任务标签列表：自动收集 framework、platform、options 等
const taskTags = (task) => {
  const tags = []
  if (task.framework) tags.push(task.framework)
  if (task.platform)  tags.push(task.platform)
  if (task.options && task.options.length) tags.push(...task.options)
  if (task.component_lib) tags.push(task.component_lib)
  return tags
}
</script>

<style scoped>
/* ═══ 任务列表 ═══ */
.task-list-wrapper {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.task-list-header {
  display: flex;
  align-items: center;
  padding: 0 2px;
}
.task-count {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}

/* ── 任务分组（紧凑条目 + 展开详情） ── */
.task-group {
  display: flex;
  flex-direction: column;
}

/* ── 紧凑条目（可点击） ── */
.task-item {
  background: rgba(0, 255, 255, 0.02);
  border-radius: 12px;
  padding: 10px 14px;
  border: 1px solid rgba(0, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.2s;
}
.task-item:hover {
  border-color: rgba(0, 255, 255, 0.18);
  background: rgba(0, 255, 255, 0.05);
}
.task-item.expanded {
  border-color: rgba(0, 255, 255, 0.2);
  background: rgba(0, 255, 255, 0.06);
  border-radius: 12px 12px 0 0;
  border-bottom-color: transparent;
}

.task-item-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}
.task-status-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
}
.task-status-badge .status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

/* 状态颜色 */
.status--pending {
  background: rgba(59, 130, 246, 0.1);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.2);
}
.status--pending .status-dot { background: #60a5fa; animation: pulse 2s ease-in-out infinite; }
.status--running {
  background: rgba(250, 204, 21, 0.1);
  color: #facc15;
  border: 1px solid rgba(250, 204, 21, 0.2);
}
.status--running .status-dot { background: #facc15; animation: pulse 1s ease-in-out infinite; }
.status--success {
  background: rgba(52, 211, 153, 0.1);
  color: #34d399;
  border: 1px solid rgba(52, 211, 153, 0.2);
}
.status--success .status-dot { background: #34d399; }
.status--failed {
  background: rgba(248, 113, 113, 0.1);
  color: #f87171;
  border: 1px solid rgba(248, 113, 113, 0.2);
}
.status--failed .status-dot { background: #f87171; }

.task-framework {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.3);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.task-item-title {
  font-size: 13px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.75);
  flex: 1;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 紧凑进度条 */
.task-progress {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}
.task-progress-bar {
  flex: 1;
  height: 5px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 3px;
  overflow: hidden;
}
.task-progress-fill {
  height: 100%;
  border-radius: 3px;
  background: linear-gradient(90deg, #3b82f6, #22d3ee);
  transition: width 0.5s ease-out;
}
.task-progress-fill.animated {
  background: linear-gradient(90deg, #3b82f6, #22d3ee, #3b82f6);
  background-size: 200% 100%;
  animation: progressPulse 1.5s ease-in-out infinite;
}
@keyframes progressPulse {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}
.task-progress-num {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  font-weight: 600;
  flex-shrink: 0;
}
.task-error {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: #f87171;
  margin-bottom: 6px;
}
.task-error i { font-size: 12px; flex-shrink: 0; }
.task-item-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.55);
}

/* ── 展开详情卡片 ── */
.task-detail {
  background: rgba(0, 255, 255, 0.03);
  border-radius: 0 0 14px 14px;
  padding: 16px 14px 14px;
  border: 1px solid rgba(0, 255, 255, 0.1);
  border-top: none;
  position: relative;
  overflow: hidden;
}

/* 卡片头部：缩略图 + 信息 + 时间 */
.card-head {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 14px;
}
.card-thumb {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(135deg, rgba(59,130,246,0.12), rgba(34,211,238,0.08));
  display: flex;
  align-items: center;
  justify-content: center;
}
.card-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.card-thumb i {
  font-size: 16px;
  color: rgba(148, 163, 184, 0.4);
}
.card-head-info {
  flex: 1;
  min-width: 0;
}
.card-title {
  font-size: 15px;
  font-weight: 700;
  color: #f1f5f9;
  line-height: 1.25;
  margin-bottom: 8px;
}
.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.card-tag {
  display: inline-block;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.8);
  background: rgba(148, 163, 184, 0.08);
  padding: 2px 10px;
  border-radius: 20px;
  border: 1px solid rgba(148, 163, 184, 0.12);
}
.card-time {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.3);
  font-family: 'Fira Code', monospace;
  flex-shrink: 0;
  padding-top: 2px;
}

/* 进度条 */
.progress-block {
  margin-bottom: 16px;
}
.progress-label {
  font-size: 11px;
  font-weight: 600;
  color: rgba(148, 163, 184, 0.6);
  letter-spacing: 0.3px;
}
.progress-bar {
  height: 5px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
  overflow: hidden;
  margin-top: 6px;
}
.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #2563eb, #22d3ee);
  border-radius: 3px;
  transition: width 1s ease-out;
  position: relative;
}
.progress-shine {
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 1px;
  background: rgba(255, 255, 255, 0.3);
}

/* 时间轴 */
.timeline {
  position: relative;
  padding-left: 8px;
}
.timeline-line {
  position: absolute;
  left: 19px;
  top: 6px;
  bottom: 6px;
  width: 1px;
  background: rgba(255, 255, 255, 0.08);
}
.timeline-step {
  display: flex;
  align-items: center;
  gap: 10px;
  justify-content: space-between;
  margin-bottom: 8px;
  position: relative;
  z-index: 1;
}
.timeline-step:last-child { margin-bottom: 0; }

/* 左侧分组（indicator + body） */
.step-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

/* 步骤指示器（仅图标） */
.step-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  transition: all 0.3s;
}
.step-icon {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.2);
  transition: all 0.3s;
}
.step-indicator.done {
  background: rgba(52, 211, 153, 0.1);
  border-color: rgba(52, 211, 153, 0.3);
}
.step-indicator.done .step-icon {
  color: rgba(52, 211, 153, 0.5);
}

.step-dot {
  width: 22px; height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(10, 10, 15, 0.95);
  border: 2px solid rgba(255, 255, 255, 0.1);
  font-size: 10px;
  color: transparent;
  transition: all 0.3s;
}
.step-dot.done {
  border-color: rgba(52, 211, 153, 0.5);
  background: rgba(52, 211, 153, 0.1);
  color: #34d399;
}
.step-body {
  min-width: 0;
  padding-top: 2px;
}
.step-title {
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
  transition: color 0.3s;
}
.step-time {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.25);
  font-family: 'Fira Code', monospace;
  transition: color 0.3s;
}

/* 阶段色：已完成 → 白色 */
.phase-past .step-title { color: rgba(255, 255, 255, 0.9); }
.phase-past .step-time  { color: rgba(255, 255, 255, 0.4); }

/* 阶段色：当前进行中 → 蓝色高亮 */
.phase-current .step-title { color: #60a5fa; }
.phase-current .step-time  { color: rgba(96, 165, 250, 0.6); }

/* 阶段色：未开始 → 保持默认灰 */
/* .phase-future 使用默认色 */

/* ═══ 卡片淡入过渡 ═══ */
.card-fade-enter-active {
  transition: opacity 0.35s ease, transform 0.35s ease;
}
.card-fade-leave-active {
  transition: opacity 0.15s ease;
}
.card-fade-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
.card-fade-leave-to { opacity: 0; }

/* ═══ 详情展开过渡 ═══ */
.detail-expand-enter-active {
  transition: all 0.3s ease-out;
}
.detail-expand-leave-active {
  transition: all 0.2s ease-in;
}
.detail-expand-enter-from {
  opacity: 0;
  max-height: 0;
  transform: translateY(-8px);
}
.detail-expand-enter-to {
  opacity: 1;
  max-height: 500px;
  transform: translateY(0);
}
.detail-expand-leave-from {
  opacity: 1;
  max-height: 500px;
}
.detail-expand-leave-to {
  opacity: 0;
  max-height: 0;
}

/* 空态 */
.task-empty {
  text-align: center;
  padding: 32px 16px;
  color: rgba(255, 255, 255, 0.15);
}
.task-empty i {
  font-size: 32px;
  margin-bottom: 8px;
  display: block;
}
.task-empty p { font-size: 12px; }

/* 查看详情按钮 */
.view-detail-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  width: 100%;
  padding: 6px 12px;
  margin-top: 12px;
  border: 1px solid rgba(0,255,255,0.12);
  border-radius: 8px;
  background: rgba(0,255,255,0.04);
  color: rgba(0,255,255,0.55);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
  justify-content: center;
}
.view-detail-btn:hover {
  background: rgba(0,255,255,0.08);
  border-color: rgba(0,255,255,0.25);
  color: #00d8ff;
}
</style>
