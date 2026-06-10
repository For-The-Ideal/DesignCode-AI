<template>
  <div ref="editorContainer" class="monaco-editor-core" :style="{ height }"></div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, computed } from 'vue'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'

if (process.client && !window.MonacoEnvironment) {
  window.MonacoEnvironment = {
    getWorker(_, label) {
      switch (label) {
        case 'typescript':
        case 'javascript':
          return new tsWorker()
        case 'json':
          return new jsonWorker()
        case 'css':
        case 'scss':
        case 'less':
          return new cssWorker()
        case 'html':
        case 'handlebars':
        case 'razor':
          return new htmlWorker()
        default:
          return new editorWorker()
      }
    }
  }
}

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
  },
  // 流式渲染时是否自动滚动到底部
  autoScroll: {
    type: Boolean,
    default: false
  },
})

const emit = defineEmits(['update:value', 'change', 'focus', 'blur', 'ready'])

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

// ═══ 流式滚动跟随状态 ═══
let rafId = null
let userScrolled = false  // 用户手动滚动时暂停自动跟随
let lastUserScrollTop = 0

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

      // 流式渲染自动滚动到底部
      if (props.autoScroll && !userScrolled) {
        scheduleScrollToBottom()
      }
    })

    // 检测用户手动滚动：向上滚动时暂停自动跟随
    editor.onDidScrollChange((e) => {
      if (!props.autoScroll) return
      const currentScrollTop = e.scrollTop
      // 用户向上滚动（远离底部）→ 暂停自动跟随
      if (currentScrollTop < lastUserScrollTop) {
        userScrolled = true
      }
      // 用户滚回底部 → 恢复自动跟随
      const model = editor.getModel()
      if (model) {
        const lastLine = model.getLineCount()
        const visibleRange = editor.getVisibleRanges()
        if (visibleRange.length > 0) {
          const lastVisibleLine = visibleRange[visibleRange.length - 1].endLineNumber
          if (lastVisibleLine >= lastLine - 1) {
            userScrolled = false
          }
        }
      }
      lastUserScrollTop = currentScrollTop
    })
    
    // 监听焦点事件
    editor.onDidFocusEditorText(() => {
      emit('focus')
    })
    
    editor.onDidBlurEditorText(() => {
      emit('blur')
    })

    // 通知父组件编辑器已就绪
    emit('ready')
    
  } catch (error) {
    console.error('Monaco Editor 加载失败:', error)
  }
}

// ═══ 流式滚动核心 ═══

/** 使用 RAF 节流滚动，避免高频调用卡顿 */
const scheduleScrollToBottom = () => {
  if (rafId) return // 已有待执行的滚动，跳过
  rafId = requestAnimationFrame(() => {
    rafId = null
    doScrollToBottom()
  })
}

/** 滚动到编辑器最后一行 */
const doScrollToBottom = () => {
  if (!editor) return
  const model = editor.getModel()
  if (!model) return
  const lastLine = model.getLineCount()
  editor.revealLine(lastLine, 1) // 平滑滚动到末尾（1 = Immediate）
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
  getEditor: () => editor,
  scrollToBottom: () => {
    userScrolled = false
    doScrollToBottom()
  },
})

onMounted(() => {
  initEditor()
})

onBeforeUnmount(() => {
  if (rafId) {
    cancelAnimationFrame(rafId)
    rafId = null
  }
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