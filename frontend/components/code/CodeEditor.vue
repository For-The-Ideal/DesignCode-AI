<!-- components/code/CodeEditor.vue -->
<template>
  <ClientOnly>
    <div class="code-editor-wrapper">
      <div class="editor-header">
        <div class="editor-dots">
          <div class="dot dot-red"></div>
          <div class="dot dot-yellow"></div>
          <div class="dot dot-green"></div>
        </div>
        <div class="editor-lang">output.dart</div>
        <div class="editor-actions">
          <button class="editor-btn" @click="handleCopy" :disabled="!codeContent">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M16 1H4C2.9 1 2 1.9 2 3V17H4V3H16V1ZM19 5H8C6.9 5 6 5.9 6 7V21C6 22.1 6.9 23 8 23H19C20.1 23 21 22.1 21 21V7C21 5.9 20.1 5 19 5ZM19 21H8V7H19V21Z" fill="currentColor"/>
            </svg>
            复制
          </button>
        </div>
      </div>
      
     


      <!-- <div class="code-editor-loading">
        <div class="loading-spinner"></div>
        <span>加载编辑器中...</span>
      </div> -->

      
      <GeneratingOverlay
      v-if="showGenerating"
      :visible="showGenerating"
      :initial-progress="generatingProgress"
      @complete="onGenerateComplete"
      @progress="onGenerateProgress"
      ref="generatingOverlayRef"
    />

      <!-- <div ref="editorContainer" v-else class="editor-container"></div> -->

      <MonacoEditorCore
        v-else
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
        @update:value="handleValueUpdate"
        @change="handleChange"
      />

    </div>
    

    
      
  </ClientOnly>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import GeneratingOverlay from './GeneratingOverlay.vue'
import MonacoEditorCore from './MonacoEditorCore.vue'

// ========== Props 定义 ==========
const props = defineProps({
  // 代码内容
  modelValue: {
    type: String,
    default: ''
  },
  // 代码语言 (dart, javascript, typescript, html, vue, css, json)
  language: {
    type: String,
    default: 'dart'
  },
  // 是否只读
  readonly: {
    type: Boolean,
    default: false
  },
  // 是否显示行号
  lineNumbers: {
    type: Boolean,
    default: true
  },
  // 是否自动换行
  wordWrap: {
    type: Boolean,
    default: true
  },
  // 字体大小
  fontSize: {
    type: Number,
    default: 13
  },
  // 主题 (vs-dark, vs, hc-black)
  theme: {
    type: String,
    default: 'vs-dark'
  },
  // 编辑器高度
  height: {
    type: [String, Number],
    default: '400px'
  },
  // 占位符文本
  placeholder: {
    type: String,
    default: '// 代码将显示在这里...'
  },
  generating:{ // 是否正在生成中
    type: Boolean,
    default: false
  },
  // 生成进度
  generatingProgress: {
    type: Number,
    default: 80
  },
})

const emit = defineEmits(['update:modelValue', 'change', 'copy', 'format'])

// ========== 响应式数据 ==========
const editorContainer = ref(null)
let editor = null
let monaco = null

const showGenerating = computed(() => props.generating)

// 语言映射（用于显示）
const languageMap = {
  dart: { label: 'Dart', monaco: 'dart' },
  flutter: { label: 'Flutter', monaco: 'dart' },
  javascript: { label: 'JavaScript', monaco: 'javascript' },
  js: { label: 'JavaScript', monaco: 'javascript' },
  typescript: { label: 'TypeScript', monaco: 'typescript' },
  ts: { label: 'TypeScript', monaco: 'typescript' },
  vue: { label: 'Vue', monaco: 'html' },
  html: { label: 'HTML', monaco: 'html' },
  css: { label: 'CSS', monaco: 'css' },
  json: { label: 'JSON', monaco: 'json' },
  python: { label: 'Python', monaco: 'python' },
  go: { label: 'Go', monaco: 'go' },
  rust: { label: 'Rust', monaco: 'rust' },
  java: { label: 'Java', monaco: 'java' },
  cpp: { label: 'C++', monaco: 'cpp' },
  csharp: { label: 'C#', monaco: 'csharp' },
  php: { label: 'PHP', monaco: 'php' },
  sql: { label: 'SQL', monaco: 'sql' },
  yaml: { label: 'YAML', monaco: 'yaml' },
  markdown: { label: 'Markdown', monaco: 'markdown' }
}

// 计算显示的代码内容
const codeContent = computed(() => props.modelValue || '')

// 计算显示的语言名称
const displayLanguage = computed(() => {
  const lang = languageMap[props.language] || languageMap.dart
  return lang.label
})

// 获取 Monaco 语言标识
const getMonacoLanguage = () => {
  const lang = languageMap[props.language] || languageMap.dart
  return lang.monaco
}

// ========== 初始化 Monaco Editor ==========
const initMonaco = async () => {
  if (!process.client) return
  
  try {
    // 动态导入 Monaco Editor
    const monacoModule = await import('monaco-editor')
    monaco = monacoModule.default || monacoModule
    
    if (!editorContainer.value) return
    
    // 创建编辑器实例
    editor = monaco.editor.create(editorContainer.value, {
      value: codeContent.value,
      language: getMonacoLanguage(),
      theme: props.theme,
      readOnly: props.readonly,
      minimap: { enabled: false },
      fontSize: props.fontSize,
      fontFamily: 'Fira Code, Consolas, monospace',
      lineNumbers: props.lineNumbers ? 'on' : 'off',
      scrollBeyondLastLine: false,
      automaticLayout: true,
      tabSize: 2,
      wordWrap: props.wordWrap ? 'on' : 'off',
      renderWhitespace: 'boundary',
      renderLineHighlight: 'all',
      cursorBlinking: 'smooth',
      cursorSmoothCaretAnimation: 'on',
      formatOnPaste: true,
      formatOnType: false,
      placeholder: props.placeholder
    })
    
    // 监听内容变化
    editor.onDidChangeModelContent(() => {
      const value = editor.getValue()
      if (value !== props.modelValue) {
        emit('update:modelValue', value)
        emit('change', value)
      }
    })
    
    // 监听失去焦点
    editor.onDidBlurEditorText(() => {
      emit('change', editor.getValue())
    })
    
  } catch (error) {
    console.error('Monaco Editor 加载失败:', error)
  }
}

// ========== 监听外部值变化 ==========
watch(() => props.modelValue, (newValue) => {
  if (editor && editor.getValue() !== newValue) {
    editor.setValue(newValue || '')
  }
})

// ========== 监听语言变化 ==========
watch(() => props.language, () => {
  if (editor && monaco) {
    const model = editor.getModel()
    if (model) {
      monaco.editor.setModelLanguage(model, getMonacoLanguage())
    }
  }
})

// ========== 监听只读变化 ==========
watch(() => props.readonly, (newVal) => {
  if (editor) {
    editor.updateOptions({ readOnly: newVal })
  }
})

// ========== 监听主题变化 ==========
watch(() => props.theme, (newVal) => {
  if (editor && monaco) {
    monaco.editor.setTheme(newVal)
  }
})

// ========== 监听字体大小变化 ==========
watch(() => props.fontSize, (newVal) => {
  if (editor) {
    editor.updateOptions({ fontSize: newVal })
  }
})

// ========== 监听自动换行变化 ==========
watch(() => props.wordWrap, (newVal) => {
  if (editor) {
    editor.updateOptions({ wordWrap: newVal ? 'on' : 'off' })
  }
})

const onGenerateProgress  = ()=>{

}

// ========== 复制代码 ==========
const handleCopy = async () => {
  const code = editor?.getValue() || props.modelValue
  if (!code) return
  
  try {
    await navigator.clipboard.writeText(code)
    emit('copy', code)
    
    // 显示提示
    const btn = document.querySelector('.editor-actions .editor-btn:first-child')
    if (btn) {
      const originalHTML = btn.innerHTML
      btn.innerHTML = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none"><path d="M9 16.17L4.83 12L3.41 13.41L9 19L21 7L19.59 5.59L9 16.17Z" fill="currentColor"/></svg> 已复制'
      setTimeout(() => {
        btn.innerHTML = originalHTML
      }, 2000)
    }
  } catch (err) {
    console.error('复制失败:', err)
  }
}

// ========== 格式化代码 ==========
const handleFormat = () => {
  if (!editor) return
  
  try {
    editor.getAction('editor.action.formatDocument').run()
    emit('format', editor.getValue())
  } catch (error) {
    console.error('格式化失败:', error)
  }
}

// ========== 获取编辑器实例（供父组件调用）=========
const getEditor = () => editor

// ========== 设置值 ==========
const setValue = (value) => {
  if (editor) {
    editor.setValue(value || '')
  }
}

// ========== 获取值 ==========
const getValue = () => {
  return editor?.getValue() || props.modelValue
}

// ========== 聚焦 ==========
const focus = () => {
  editor?.focus()
}

// ========== 生命周期 ==========
onMounted(() => {
  initMonaco()
})

onBeforeUnmount(() => {
  if (editor) {
    editor.dispose()
  }
})

// 暴露方法给父组件
defineExpose({
  getEditor,
  setValue,
  getValue,
  focus,
  format: handleFormat,
  copy: handleCopy
})
</script>

<style scoped>
.code-editor-wrapper {
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

.editor-lang {
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  font-weight: 500;
  color: #00ffff;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  background: rgba(0, 255, 255, 0.1);
  padding: 4px 10px;
  border-radius: 20px;
}

.editor-actions {
  display: flex;
  gap: 8px;
}

.editor-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  background: rgba(0, 255, 255, 0.08);
  border: 1px solid rgba(0, 255, 255, 0.25);
  border-radius: 8px;
  color: #00ffff;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.editor-btn:hover:not(:disabled) {
  background: rgba(0, 255, 255, 0.18);
  border-color: rgba(0, 255, 255, 0.5);
  transform: translateY(-1px);
}

.editor-btn:active {
  transform: translateY(0);
}

.editor-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.editor-container {
  width: 100%;
  height: v-bind(height);
}

.code-editor-loading {
  background: #0a0a0f;
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 16px;
  height: v-bind(height);
  display: flex;
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