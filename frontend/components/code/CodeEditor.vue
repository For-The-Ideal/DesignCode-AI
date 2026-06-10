<template>
  <ClientOnly>
    <div class="code-editor-wrapper">
      <div class="editor-header">
        <div class="editor-dots">
          <div class="dot dot-red"></div>
          <div class="dot dot-yellow"></div>
          <div class="dot dot-green"></div>
        </div>
      </div>

      <div v-show="isLoading" class="code-editor-loading">
        <div class="loading-spinner"></div>
        <span>加载编辑器中...</span>
      </div>

      <MonacoEditorCore
        v-show="!isLoading"
        ref="editorCoreRef"
        :value="modelValue"
        :language="language"
        :readonly="readonly"
        :line-numbers="lineNumbers"
        :word-wrap="wordWrap"
        :font-size="fontSize"
        :theme="theme"
        :height="height"
        :placeholder="placeholder"
        :auto-scroll="autoScroll"
        @ready="onEditorReady"
        @change="onEditorChange"
      />
    </div>
  </ClientOnly>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'
import MonacoEditorCore from './MonacoEditorCore.vue'

// ========== Props 定义 ==========
const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  language: {
    type: String,
    default: 'dart'
  },
  readonly: {
    type: Boolean,
    default: false
  },
  lineNumbers: {
    type: Boolean,
    default: true
  },
  wordWrap: {
    type: Boolean,
    default: true
  },
  fontSize: {
    type: Number,
    default: 13
  },
  theme: {
    type: String,
    default: 'vs-dark'
  },
  height: {
    type: [String, Number],
    default: '400px'
  },
  placeholder: {
    type: String,
    default: '// 代码将显示在这里...'
  },
  // 流式渲染时是否自动滚动到末尾
  autoScroll: {
    type: Boolean,
    default: false
  },
})

const emit = defineEmits(['update:modelValue', 'change', 'copy', 'format'])

// ========== 响应式数据 ==========
const editorCoreRef = ref(null)
const isLoading = ref(true)

// ========== 方法 ==========
function onEditorReady() {
  isLoading.value = false
}

function onEditorChange(value) {
  emit('update:modelValue', value)
  emit('change', value)
}

// ========== 复制代码 ==========
const copy = async () => {
  const code = editorCoreRef.value?.getValue() || props.modelValue
  if (!code) return

  try {
    await navigator.clipboard.writeText(code)
    emit('copy', code)
  } catch (err) {
    console.error('复制失败:', err)
  }
}

// ========== 获取编辑器实例 ==========
const getEditor = () => editorCoreRef.value?.getEditor()

// ========== 设置值 ==========
const setValue = (value) => {
  editorCoreRef.value?.setValue(value || '')
}

// ========== 获取值 ==========
const getValue = () => {
  return editorCoreRef.value?.getValue() || props.modelValue
}

// ========== 聚焦 ==========
const focus = () => {
  editorCoreRef.value?.focus()
}

// ========== 暴露方法 ==========
defineExpose({
  getEditor,
  setValue,
  getValue,
  focus,
  copy
})
</script>

<style scoped>
.code-editor-wrapper {
  position: relative;
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
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
}

.editor-dots {
  display: flex;
  gap: 8px;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  transition: all 0.2s;
}

.dot-red { background: #ff5f56; }
.dot-yellow { background: #ffbd2e; }
.dot-green { background: #27c93f; }

.code-editor-loading {
  z-index: 10;
  background: #0a0a0f;
  border-radius: 16px;
  display: flex;
  height: v-bind(height);
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: #888;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 2px solid rgba(0, 255, 255, 0.2);
  border-top-color: #00ffff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
