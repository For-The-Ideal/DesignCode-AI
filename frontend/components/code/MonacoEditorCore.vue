<template>
  <div ref="editorContainer" class="monaco-editor-core" :style="{ height }"></div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, computed } from 'vue'

const props = defineProps({
  // 代码内容
  value: {
    type: String,
    default: ''
  },
  // 代码语言
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
  // 主题
  theme: {
    type: String,
    default: 'vs-dark'
  },
  // 编辑器高度
  height: {
    type: [String, Number],
    default: '400px'
  },
  // 占位符
  placeholder: {
    type: String,
    default: '// 代码将显示在这里...'
  }
})

const emit = defineEmits(['update:value', 'change', 'focus', 'blur'])

// 语言映射
const languageMap = {
  dart: 'dart',
  flutter: 'dart',
  javascript: 'javascript',
  typescript: 'typescript',
  vue: 'html',
  html: 'html',
  css: 'css',
  json: 'json',
  python: 'python',
  go: 'go',
  rust: 'rust',
  java: 'java',
  cpp: 'cpp',
  php: 'php',
  sql: 'sql',
  yaml: 'yaml',
  markdown: 'markdown'
}

const getMonacoLanguage = () => {
  return languageMap[props.language] || 'dart'
}

const editorContainer = ref(null)
let editor = null
let monaco = null

// 初始化编辑器
const initEditor = async () => {
  if (!process.client) return
  
  try {
    const monacoModule = await import('monaco-editor')
    monaco = monacoModule.default || monacoModule
    
    if (!editorContainer.value) return
    
    editor = monaco.editor.create(editorContainer.value, {
      value: props.value,
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
      placeholder: props.placeholder,
      // 隐藏滚动条
      scrollbar: {
        vertical: 'hidden',
        horizontal: 'hidden',
        handleMouseWheel: true
      },
      // 隐藏右侧缩略图旁边的装饰
      overviewRulerLanes: 0,
      overviewRulerBorder: false
    })
    
    // 监听内容变化
    editor.onDidChangeModelContent(() => {
      const value = editor.getValue()
      emit('update:value', value)
      emit('change', value)
    })
    
    // 监听焦点事件
    editor.onDidFocusEditorText(() => {
      emit('focus')
    })
    
    editor.onDidBlurEditorText(() => {
      emit('blur')
    })
    
  } catch (error) {
    console.error('Monaco Editor 加载失败:', error)
  }
}

// 监听外部值变化
watch(() => props.value, (newValue) => {
  if (editor && editor.getValue() !== newValue) {
    editor.setValue(newValue || '')
  }
})

// 监听语言变化
watch(() => props.language, () => {
  if (editor && monaco) {
    const model = editor.getModel()
    if (model) {
      monaco.editor.setModelLanguage(model, getMonacoLanguage())
    }
  }
})

// 监听只读变化
watch(() => props.readonly, (newVal) => {
  if (editor) {
    editor.updateOptions({ readOnly: newVal })
  }
})

// 监听字体大小
watch(() => props.fontSize, (newVal) => {
  if (editor) {
    editor.updateOptions({ fontSize: newVal })
  }
})

// 监听自动换行
watch(() => props.wordWrap, (newVal) => {
  if (editor) {
    editor.updateOptions({ wordWrap: newVal ? 'on' : 'off' })
  }
})

// 监听行号
watch(() => props.lineNumbers, (newVal) => {
  if (editor) {
    editor.updateOptions({ lineNumbers: newVal ? 'on' : 'off' })
  }
})

// 暴露方法
const getValue = () => {
  return editor?.getValue() || props.value
}

const setValue = (val) => {
  if (editor) {
    editor.setValue(val || '')
  }
}

const format = () => {
  if (editor) {
    editor.getAction('editor.action.formatDocument')?.run()
  }
}

const focus = () => {
  editor?.focus()
}

defineExpose({
  getValue,
  setValue,
  format,
  focus,
  getEditor: () => editor
})

onMounted(() => {
  initEditor()
})

onBeforeUnmount(() => {
  if (editor) {
    editor.dispose()
  }
})
</script>

<style scoped>
.monaco-editor-core {
  width: 100%;
  background: #0a0a0f;
  text-align: left;
}

</style>