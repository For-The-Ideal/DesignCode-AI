<template>
  <div class="screen" ref="screenRef">
    <div class="content">
      <div class="header-title">
        <h2>AI实时代码生成</h2>
      </div>
      <div class="sse-layout">
        <!-- 左侧代码区 -->
        <div class="code-area">
          <CodeEditor
            v-model="template.templateCode"
            language="dart"
            :readonly="true"
            :auto-scroll="true"
            height="700px"
            placeholder="// 滚动到此处将自动开始流式生成 Flutter 代码..."
          />
        </div>
        <!-- 右侧手机模拟器 -->
        <div class="preview-area">
          <FlutterTemplate :html="template.previewCode" />
        </div>
      </div>
      <div class="scroll-hint"><span>▼ 继续滚动查看 AI 智能诊断 ▼</span></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useGeneration } from '@/composables/useGeneration'
import CodeEditor from '@/components/code/CodeEditor.vue'
import FlutterTemplate from '@/components/common/FlutterTemplate.vue'

const screenRef = ref(null)

const {
  template,
  initTemplateData,
  pauseStreaming,
  cleanup,
} = useGeneration()

let observer = null

onMounted(() => {
  observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        initTemplateData()
      } else {
        pauseStreaming()
      }
    })
  }, { threshold: 0.1 })

  if (screenRef.value) observer.observe(screenRef.value)
})

onUnmounted(() => {
  cleanup()
  if (observer) observer.disconnect()
})
</script>

<style scoped>
.screen {
  scroll-snap-align: start;
  padding: 20px;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
}
.content {
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
}
.header-title {
  text-align: center;
  margin-bottom: 2rem;
}
.header-title h2 {
  font-family: 'Orbitron', monospace;
  font-size: 2rem;
}
.sse-layout {
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
}
.code-area {
  flex: 1.2;
  max-height: 720px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}
.preview-area {
  flex: 1;
  max-height: 752px;
  backdrop-filter: blur(16px);
  border-radius: 40px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow: hidden;
}
.code-display {
  background: #01050e;
  border-radius: 40px;
  padding: 1rem;
  margin-top: 1rem;
  max-height: 752px;
  overflow-y: auto;
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  flex: 1;
}

.scroll-hint {
  position: absolute;
  bottom: 1.2rem;
  left: 50%;
  transform: translateX(-50%);
  text-align: center;
  font-size: 0.8rem;
  color: #6f9eff;
  z-index: 1;
  animation: bounce 1.8s infinite;
}
@keyframes bounce {
  0%, 100% { transform: translateX(-50%) translateY(0); }
  50% { transform: translateX(-50%) translateY(8px); }
}
@media (max-width: 968px) {
  .sse-layout { flex-direction: column; }
}
</style>
