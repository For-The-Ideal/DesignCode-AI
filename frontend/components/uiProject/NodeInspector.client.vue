<template>
  <div v-if="node" class="inspector">
    <div class="ins-head">
      <span class="ins-title">{{ displayName(node) }}</span>
      <span class="status-chip" :style="{ background: STATUS_COLOR[node.status] }">{{ statusText(node.status) }}</span>
    </div>

    <div class="ins-body">
      <!-- 顶层：从枚举内联选类型（取代旧的图上浮层） -->
      <div v-if="isTop" class="field">
        <span class="field-label">区域类型<span v-if="!form.name" class="need-pick">（必选其一）</span></span>
        <div class="type-grid">
          <button
            v-for="t in TOP_REGION_ENUM"
            :key="t"
            class="type-chip"
            :class="{ active: form.name === t, used: otherUsed[t] }"
            @click="pickType(t)"
          >{{ t }}<span v-if="otherUsed[t]" class="chip-used">已用</span></button>
        </div>
      </div>

      <div v-if="isContainer" class="field field--row">
        <span class="field-label">排列方式
          <span class="help-tip" :data-tip="LAYOUT_TIP">
            <i class="fas fa-circle-question help-badge"></i>
          </span>
        </span>
        <div class="layout-seg">
          <button
            v-for="op in LAYOUT_OPTIONS"
            :key="op.value"
            class="seg-btn"
            :class="{ active: form.layout === op.value }"
            @click="setLayout(op.value)"
          >{{ op.label }}</button>
        </div>
      </div>

      <!-- 非顶层节点：自由命名（无论叶子还是容器都可命名；顶层走上方枚举类型） -->
      <label v-if="!isTop" class="field">
        <span class="field-label">区域名</span>
        <input v-model="form.name" class="inp" maxlength="24" placeholder="子节点自由命名" />
      </label>

      <label class="field">
        <span class="field-label">一句话描述</span>
        <textarea v-model="form.summary" class="inp txt" rows="2" placeholder="这个区域是什么、含哪些元素"></textarea>
      </label>

      <label class="field">
        <span class="field-label">视觉要点</span>
        <textarea v-model="form.styleHints" class="inp txt" rows="2" placeholder="背景色 / 高度 / 分割线 / 字号 / 间距"></textarea>
      </label>

      <label v-if="isTop || hasChildren" class="field">
        <span class="field-label">交互行为</span>
        <input v-model="form.behavior" class="inp" placeholder="如 固定吸顶 / 下拉刷新，可空" />
      </label>
    </div>
  </div>
  <div v-else class="inspector inspector--empty">
    选中左侧框或右侧节点后可编辑描述
  </div>
</template>

<script setup>
import { reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useUiProjectStore } from '~/stores/uiProject'
import { TOP_REGION_ENUM } from '~/utils/uiTree.js'

const store = useUiProjectStore()

const STATUS_COLOR = {
  pending: 'rgba(255,255,255,0.3)',
  ai_draft: '#ffcf4d',
  confirmed: '#37e08d',
  empty: '#ff4757',
}
const STATUS_TEXT = { pending: '待描述', ai_draft: 'AI 初稿', confirmed: '已确认', empty: '待补全' }
const LAYOUT_OPTIONS = [
  { value: 'row', label: '横排 ⇄' },
  { value: 'column', label: '竖排 ⇅' },
  { value: 'grid', label: '网格 ⊞' },
]

// 排列方式的优先级说明：用 data-tip 交给纯 CSS 气泡渲染（不用 v-title，避免与原生 title 气泡叠层）
const LAYOUT_TIP = '排列方式取值优先级：手动 ＞ 自动推断 ＞ AI 建议。画完框选了类型时，若顶层 AI 给了排列方向会自动选中；AI 没给就留空不选。增删子区域后会按子框位置重新自动推断，手动选定过就不再被覆盖。'
const LAYOUT_SRC_TEXT = { manual: '手动', ai: 'AI 建议', auto: '自动推断' }

const node = computed(() => store.selectedNode)
const isTop = computed(() => !!node.value && !node.value.parentId)
// 是否有子节点：非顶层但已下钻子级时，它充当父容器，同样需要交互行为字段
const hasChildren = computed(() => !!node.value && (node.value.childrenIds?.length || 0) > 0)
// 其它顶层节点已占用的枚举（用于“已用”标记，仍可点选自动生成序号）
const otherUsed = computed(() => {
  const set = {}
  if (!node.value) return set
  store.nodeList
    .filter((n) => !n.parentId && n.nodeKey !== node.value.nodeKey && n.name)
    .forEach((n) => { set[n.name] = (set[n.name] || 0) + 1 })
  return set
})
const isContainer = computed(() =>
  !!(node.value && (!node.value.parentId || node.value.childrenIds?.length || node.value.nodeType === 'container')),
)

const statusText = (s) => STATUS_TEXT[s] || s
function displayName(n) {
  if (!n.name) return '未命名区域'
  return n.variantIndex > 0 ? `${n.name}-${n.variantIndex + 1}` : n.name
}

// 表单本地态，改动 debounce 回写 store（触发 markEdited / manual layout）
const form = reactive({ name: '', summary: '', styleHints: '', behavior: '', layout: '' })
let timer = null

function pullFromNode(n) {
  form.name = n?.name || ''
  form.summary = n?.summary || ''
  form.styleHints = n?.styleHints || ''
  form.behavior = n?.behavior || ''
  form.layout = n?.layout || ''
}

watch(node, (n) => { if (n) pullFromNode(n) }, { immediate: true })

function pushToStore() {
  if (!node.value) return
  store.updateNode(node.value.nodeKey, {
    name: form.name,
    summary: form.summary,
    styleHints: form.styleHints,
    behavior: form.behavior,
  })
}

// name/layout 单独处理：layout 走 manual，name 走 updateNode
watch(() => [form.name, form.summary, form.styleHints, form.behavior], () => {
  clearTimeout(timer)
  timer = setTimeout(pushToStore, 300)
})

function setLayout(val) {
  if (!node.value) return
  form.layout = val
  store.updateNode(node.value.nodeKey, { layout: val })
}

// 顶层选类型：指派名称 + 匹配 topRegions 自动填充，刷新表单；位置异常软提示
function pickType(name) {
  if (!node.value) return
  const res = store.assignTopName(node.value.nodeKey, name)
  pullFromNode(node.value)
  if (res && res.ok && res.hints && res.hints.length) ElMessage.warning(`「${name}」${res.hints[0]}`)
}
</script>

<style scoped>
.inspector {
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  background: rgba(0, 255, 255, 0.04);
  border: 1px solid rgba(0, 255, 255, 0.14);
  overflow: hidden;
}
.inspector--empty {
  padding: 20px;
  text-align: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.28);
}
.ins-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid rgba(0, 255, 255, 0.1);
}
.ins-title { font-size: 13px; font-weight: 600; color: rgba(255, 255, 255, 0.85); }
.status-chip {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  color: #0a0a0f;
  font-weight: 700;
}
.ins-body { padding: 12px; display: flex; flex-direction: column; gap: 12px; overflow-y: auto; container-type: inline-size; }
.field { display: flex; flex-direction: column; gap: 5px; }
.field--row { flex-direction: row; flex-wrap: wrap; align-items: center; gap: 8px; }
.field-label { font-size: 12px; color: rgba(255, 255, 255, 0.45); }
.inp {
  width: 100%;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(0, 255, 255, 0.22);
  border-radius: 8px;
  padding: 7px 9px;
  color: #fff;
  font-size: 13px;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s;
}
.inp:focus { border-color: #00ffff; box-shadow: 0 0 8px rgba(0, 255, 255, 0.1); }
.txt { resize: vertical; }
.field--row .inp { flex: none; width: auto; }
.layout-seg { display: flex; gap: 6px; }
.seg-btn {
  padding: 4px 10px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.3);
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  cursor: pointer;
}
.seg-btn.active { border-color: #00ffff; color: #00ffff; background: rgba(0, 255, 255, 0.1); }
.src-chip {
  margin-left: 6px;
  padding: 1px 6px;
  border-radius: 8px;
  font-size: 10px;
  color: rgba(0, 255, 255, 0.75);
  background: rgba(0, 255, 255, 0.08);
  border: 1px solid rgba(0, 255, 255, 0.22);
}
/* 字段说明气泡：圆形徽章 + 自绘 tooltip（跟随文字换行，不依赖指令/原生 title） */
/* cursor 不用 help：Windows 的 help 光标自带一个小问号，会与徽章上的问号重复 */
.help-tip { position: relative; display: inline-flex; vertical-align: middle; margin-left: 5px; cursor: pointer; }
.help-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  font-size: 10px;
  color: #00ffff;
  background: rgba(0, 255, 255, 0.12);
  border: 1px solid rgba(0, 255, 255, 0.45);
  transition: all 0.2s;
}
.help-tip:hover .help-badge { background: rgba(0, 255, 255, 0.25); box-shadow: 0 0 8px rgba(0, 255, 255, 0.35); }
.help-tip::after {
  content: attr(data-tip);
  position: absolute;
  left: 0;
  bottom: calc(100% + 8px);
  /* 限定在滚动面板内（.ins-body 为 inline-size 容器），避免比面板宽被 overflow 裁掉 */
  width: min(240px, calc(100cqw - 12px));
  padding: 8px 10px;
  border-radius: 8px;
  background: rgba(12, 14, 22, 0.97);
  border: 1px solid rgba(0, 255, 255, 0.32);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.5);
  color: rgba(255, 255, 255, 0.82);
  font-size: 12px;
  line-height: 1.6;
  white-space: normal;
  opacity: 0;
  visibility: hidden;
  transform: translateY(4px);
  transition: opacity 0.15s, transform 0.15s, visibility 0s linear 0.15s;
  z-index: 30;
  pointer-events: none;
}
.help-tip:hover::after { opacity: 1; visibility: visible; transform: translateY(0); transition-delay: 0s; }

/* 顶层类型内联胶囊 */
.need-pick { color: #ffcf4d; margin-left: 4px; }
.type-grid { display: flex; flex-wrap: wrap; gap: 6px; }
.type-chip {
  position: relative;
  padding: 5px 10px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.3);
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.type-chip:hover { border-color: rgba(0, 255, 255, 0.5); color: #00ffff; }
.type-chip.active { border-color: #00ffff; color: #00ffff; background: rgba(0, 255, 255, 0.12); }
.type-chip.used { opacity: 0.62; }
.chip-used { margin-left: 4px; font-size: 10px; color: rgba(255, 207, 77, 0.85); }
</style>
