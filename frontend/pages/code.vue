<template>
  <div class="code-page">
    <CodeSidebar />

    <main class="code-main" @click="requireLogin">
      <div class="code-header">
        <div class="header-left">
          <h2 class="header-title">AI 智能代码生成</h2>
          <p class="header-sub">上传设计稿，AI 自动识别并生成高质量代码</p>
        </div>
      </div>

      <section class="upload-section">
        <div class="upload-section-inner">
          <div class="upload-left">
            <div class="glow-section">
              <div class="section-label">
                <span class="label-num">1</span> 
                <span class="section-label-text">上传设计稿</span>
              </div>
              <UploadZone />
            </div>
          </div>

          <div class="upload-right">
            <div class="glow-section">
              <ConfigPanel />
            </div>
          </div>
        </div>
      </section>

      <div class="glow-section">
        <FlowSteps :activeStep="-1" />
      </div>
    </main>
  </div>
</template>

<script setup>
import { watch } from 'vue'
import CodeSidebar from '~/components/code/CodeSidebar.vue'
import UploadZone from '~/components/upload/UploadZone.vue'
import ConfigPanel from '~/components/code/ConfigPanel.vue'
import FlowSteps from '~/components/code/FlowSteps.vue'
import { useCodeStore } from '~/stores/code'
import { useUserStore } from '~/stores/user'
import { useCommonStore } from '~/stores/common'

const store = useCodeStore()
const userStore = useUserStore()
const commonStore = useCommonStore()

// ═══ 登录保护 ═══
const requireLogin = () => {
  if (!userStore.isLogin) {
    commonStore.setLoginModalVisible(true)
  }
}

// ═══ 初始化 ═══
watch(() => userStore.isLogin, (val) => {
  if (val) store.fetchUserTasks()
}, { immediate: true })
</script>

<style scoped>
.code-page {
  display: flex;
  min-height: calc(100vh - 140px);
  background: #0a0a0f;
  position: relative;
}

.code-main {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.glow-section {
  position: relative;
  border-radius: 16px;
  border: 1px solid rgba(0, 255, 255, 0.08);
  background: rgba(10, 14, 23, 0.7);
  backdrop-filter: blur(16px);
  padding: 20px 24px;
  transition: all 0.4s;
  overflow: hidden;
}
.glow-section::before {
  content: '';
  position: absolute;
  top: -1px; left: -1px; right: -1px; bottom: -1px;
  border-radius: 17px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.12), transparent 30%, transparent 70%, rgba(255, 0, 255, 0.08));
  z-index: 0;
  pointer-events: none;
}
.glow-section:hover {
  border-color: rgba(0, 255, 255, 0.18);
  box-shadow: 0 0 30px rgba(0, 255, 255, 0.04), inset 0 0 30px rgba(0, 255, 255, 0.02);
}
.glow-section > * { position: relative; z-index: 1; }

.code-header {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  padding: 20px 24px;
  border-radius: 16px;
}
.header-title { font-size: 28px; font-weight: 800; color: #e2e8f0; letter-spacing: -0.3px; }
.header-sub { font-size: 13px; color: rgba(255,255,255,0.45); margin-top: 4px; }

.upload-section-inner { display: flex; gap: 24px; }
.upload-left { flex: 1.5; display: flex; flex-direction: column; gap: 12px; }
.upload-left .glow-section,
.upload-right .glow-section { flex: 1; display: flex; flex-direction: column; padding: 20px; }

.section-label {
  font-size: 14px;
  font-weight: 600;
  color: #60a5fa;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  margin-bottom: 16px;
  display: inline-block;
  padding-left: 14px;
}
.section-label::before {
  content: '';
  position: absolute;
  left: 0; top: 2px; bottom: 2px;
  width: 2px;
  background: linear-gradient(135deg, #60a5fa, #818cf8);
  border-radius: 1px;
}
.label-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px; height: 24px;
  border-radius: 8px;
  font-size: 12px;
  background: #60a5fa;
  color: #fff;
  font-weight: 700;
  margin-right: 6px;
}
.section-label-text { font-weight: 700; color: #fff; }
.upload-right { flex: 1; display: flex; flex-direction: column; gap: 16px; }

@media (max-width: 1100px) {
  .upload-section-inner { flex-direction: column; }
}
</style>
