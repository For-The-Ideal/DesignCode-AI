<template>
  <div class="code-editor-wrapper" :style="{ height: pxHeight }">
    <div class="editor-header">
      <div class="editor-dots">
        <div class="dot dot-red"></div>
        <div class="dot dot-yellow"></div>
        <div class="dot dot-green"></div>
      </div>
      <div class="editor-title">{{ language }}</div>
      <div class="editor-actions">
        <button class="editor-btn" @click="copy" title="复制代码">
          <i class="fas fa-copy"></i>
          <span>复制</span>
        </button>
      </div>
    </div>

    <!-- 轻量只读代码视图：秒开、无 Monaco、无“加载编辑器中…”卡死 -->
    <div
      ref="bodyRef"
      class="code-body"
      :class="{ 'code-body--wrap': wordWrap }"
      :style="{ fontSize: fontSize + 'px' }"
    >
      <div v-if="!hasContent" class="code-placeholder">{{ placeholder }}</div>
      <div v-else class="code-lines">
        <div class="ln-row" v-for="(line, i) in lines" :key="i">
          <span v-if="lineNumbers" class="ln-num">{{ i + 1 }}</span>
          <span class="ln-code" v-html="htmlLines[i] || ' '"></span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'

// ========== Props（与原 Monaco 版保持一致，替换不影响调用方） ==========
const props = defineProps({
  modelValue: { type: String, default: '' },
  language: { type: String, default: 'dart' },
  readonly: { type: Boolean, default: false },
  lineNumbers: { type: Boolean, default: true },
  wordWrap: { type: Boolean, default: true },
  fontSize: { type: Number, default: 13 },
  theme: { type: String, default: 'vs-dark' },
  height: { type: [String, Number], default: '400px' },
  placeholder: { type: String, default: '// 代码将显示在这里...' },
  // 流式渲染时是否自动滚动到末尾
  autoScroll: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue', 'change', 'copy', 'format'])

const bodyRef = ref(null)

const pxHeight = computed(() => (typeof props.height === 'number' ? `${props.height}px` : props.height))
const lines = computed(() => (props.modelValue || '').split('\n'))
const hasContent = computed(() => (props.modelValue || '').length > 0)

// ── 轻量逐行分词高亮（无第三方库，保持行号/换行对齐）──
const KEYWORDS = new Set(('abstract as async await break case catch class const continue default defer do dynamic else '
  + 'enum export extends external factory final finally for function get if implements import in instanceof interface '
  + 'is late let library new of operator part private protected public readonly required rethrow return set static '
  + 'super switch sync this throw true false try typedef typeof undefined null var void while with yield from').split(' '))
const escapeHtml = (s) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
const TOKEN_RE = /(\/\/[^\n]*)|(['"`])(?:\\.|(?!\2).)*?\2|(\b\d[\d_]*(?:\.\d+)?\b)|(<\/?[A-Za-z][\w-]*)|(@[A-Za-z_][\w-]*)|([A-Za-z_$][\w$]*)/g
const classifyId = (id) => KEYWORDS.has(id) ? 'tk-kw' : (/^[A-Z]/.test(id) ? 'tk-type' : '')
function tokenize(line) {
  if (!line) return ''
  let out = '', last = 0, m
  TOKEN_RE.lastIndex = 0
  while ((m = TOKEN_RE.exec(line)) !== null) {
    if (m.index > last) out += escapeHtml(line.slice(last, m.index))
    const text = escapeHtml(m[0])
    let cls = ''
    if (m[1] != null) cls = 'tk-com'
    else if (m[2] != null) cls = 'tk-str'
    else if (m[3] != null) cls = 'tk-num'
    else if (m[4] != null) cls = 'tk-tag'
    else if (m[5] != null) cls = 'tk-anno'
    else if (m[6] != null) cls = classifyId(m[6])
    out += cls ? `<span class="${cls}">${text}</span>` : text
    last = m.index + m[0].length
    if (m[0] === '') TOKEN_RE.lastIndex++
  }
  out += escapeHtml(line.slice(last))
  return out
}
const htmlLines = computed(() => lines.value.map(tokenize))

// 流式输出：内容增长时自动滚到底部（等价原 autoScroll 行为）
watch(
  () => props.modelValue,
  () => {
    if (!props.autoScroll) return
    nextTick(() => {
      const el = bodyRef.value
      if (el) el.scrollTop = el.scrollHeight
    })
  }
)

// ========== 复制代码 ==========
const copy = async () => {
  const code = props.modelValue
  if (!code) return
  try {
    await navigator.clipboard.writeText(code)
    emit('copy', code)
  } catch (err) {
    console.error('复制失败:', err)
  }
}

// ========== 只读展示：格式化无意义，保留接口兼容 ==========
const format = () => emit('format')

// ========== 暴露方法（保持与原组件一致的对外接口） ==========
defineExpose({
  getEditor: () => null,
  setValue: (value) => emit('update:modelValue', value || ''),
  getValue: () => props.modelValue,
  focus: () => {},
  copy,
  format,
})
</script>

<style scoped>
.code-editor-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
  background: #0a0a0f;
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s;
}

.code-editor-wrapper:hover {
  border-color: rgba(0, 255, 255, 0.5);
  box-shadow: 0 0 20px rgba(0, 255, 255, 0.1);
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
}

.editor-dots {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.dot-red { background: #ff5f56; }
.dot-yellow { background: #ffbd2e; }
.dot-green { background: #27c93f; }

.editor-title {
  flex: 1;
  text-align: center;
  font-size: 14px;
  color: #fff;
  text-transform: capitalize;
}

.editor-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.editor-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.45);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
  white-space: nowrap;
}
.editor-btn i { font-size: 11px; }
.editor-btn:hover {
  background: rgba(0, 255, 255, 0.1);
  border-color: rgba(0, 255, 255, 0.3);
  color: #00ffff;
}

/* ═══ 代码主体 ═══ */
.code-body {
  flex: 1;
  min-height: 0;
  overflow: auto;
  background: #0a0a0f;
  font-family: 'Fira Code', Consolas, monospace;
  line-height: 1.6;
  color: #d6deeb;
}
.code-body::-webkit-scrollbar { width: 8px; height: 8px; }
.code-body::-webkit-scrollbar-track { background: transparent; }
.code-body::-webkit-scrollbar-thumb { background: rgba(0, 255, 255, 0.15); border-radius: 4px; }

.code-placeholder {
  padding: 16px;
  color: rgba(255, 255, 255, 0.2);
  font-style: italic;
}

.code-lines {
  padding: 12px 0;
}
.ln-row {
  display: flex;
  align-items: flex-start;
}
.ln-num {
  flex-shrink: 0;
  width: 3.2em;
  padding-right: 1em;
  text-align: right;
  color: rgba(255, 255, 255, 0.18);
  user-select: none;
}
.ln-code {
  white-space: pre;
  flex: 1;
  min-width: 0;
}
.code-body--wrap .ln-code {
  white-space: pre-wrap;
  word-break: break-word;
}

/* ── 语法高亮色板（v-html 插入的子元素需 :deep 才生效）── */
.ln-code :deep(.tk-com)  { color: #5c6a86; font-style: italic; }
.ln-code :deep(.tk-str)  { color: #c3e88d; }
.ln-code :deep(.tk-num)  { color: #f78c6c; }
.ln-code :deep(.tk-kw)   { color: #c792ea; }
.ln-code :deep(.tk-type) { color: #ffcb6b; }
.ln-code :deep(.tk-tag)  { color: #f07178; }
.ln-code :deep(.tk-anno) { color: #82aaff; }
</style>
