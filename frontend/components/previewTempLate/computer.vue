<template>
  <div class="desktop-mock">
    <!-- 浏览器窗口 chrome -->
    <div class="browser-chrome">
      <div class="browser-dots">
        <div class="dot dot-red"></div>
        <div class="dot dot-yellow"></div>
        <div class="dot dot-green"></div>
      </div>
      <div class="browser-url">
        <i class="fas fa-lock"></i>
        <span>localhost:3000</span>
        <i class="fas fa-refresh"></i>
      </div>
    </div>

    <!-- 页面内容 or 加载中 -->
    <div class="browser-screen">
      <div v-if="html" class="desktop-page-content" v-html="html"></div>
      <div v-else class="desktop-page-content desktop-loading">
        <div class="scan-line"></div>
        <div class="ring-pulse">
          <div class="ring ring-1"></div>
          <div class="ring ring-2"></div>
          <div class="ring ring-3"></div>
          <div class="ring-center">
            <i class="fas fa-microchip"></i>
          </div>
        </div>
        <div class="loading-text">
          <span class="loading-label">AI 生成引擎</span>
          <div class="loading-bar">
            <div class="loading-bar-inner"></div>
          </div>
          <span class="loading-desc">正在解析设计稿，构建代码框架...</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  html: { type: String, default: '' },
})
</script>

<style lang="scss" scoped>
.desktop-mock {
  width: 100%;
  max-width: 800px;
  background: #1e1e1e;
  overflow: hidden;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.5);
}

/* 浏览器 chrome */
.browser-chrome {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: #2d2d2d;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.browser-dots {
  display: flex;
  gap: 7px;
  flex-shrink: 0;
}

.dot {
  width: 11px;
  height: 11px;
  border-radius: 50%;
}
.dot-red { background: #ff5f56; }
.dot-yellow { background: #ffbd2e; }
.dot-green { background: #27c93f; }

.browser-url {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 12px;
  background: #1a1a1a;
  border-radius: 6px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.browser-url i:first-child {
  font-size: 10px;
  color: #27c93f;
}

.browser-url span {
  flex: 1;
  color: rgba(255, 255, 255, 0.6);
}

.browser-url i:last-child {
  font-size: 11px;
  cursor: pointer;
  transition: color 0.2s;
}
.browser-url i:last-child:hover { color: rgba(255, 255, 255, 0.8); }

/* 屏幕 */
.browser-screen {
  height: 480px;
  overflow: auto;
}

/* 页面内容 */
.desktop-page-content {
  min-height: 100%;
  background: #fff;
}

/* 加载态 */
.desktop-loading {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 28px;
  background: linear-gradient(160deg, #050a18 0%, #0a1628 40%, #0d1f3c 100%);
  overflow: hidden;
}

.scan-line {
  position: absolute;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(0, 255, 255, 0.6), transparent);
  animation: scan 2.2s ease-in-out infinite;
  box-shadow: 0 0 12px rgba(0, 255, 255, 0.3);
}

@keyframes scan {
  0% { top: 0; opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { top: 100%; opacity: 0; }
}

.ring-pulse {
  position: relative;
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.ring {
  position: absolute;
  border-radius: 50%;
  border: 2px solid transparent;
}
.ring-1 {
  width: 100px; height: 100px;
  border-top-color: rgba(0, 255, 255, 0.7);
  border-right-color: rgba(255, 0, 255, 0.5);
  animation: spin 1.5s linear infinite;
}
.ring-2 {
  width: 76px; height: 76px;
  border-bottom-color: rgba(0, 255, 255, 0.5);
  border-left-color: rgba(255, 0, 255, 0.3);
  animation: spin 2s linear infinite reverse;
}
.ring-3 {
  width: 52px; height: 52px;
  border-top-color: rgba(0, 255, 255, 0.3);
  animation: spin 1.2s linear infinite;
}
.ring-center {
  width: 36px; height: 36px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(0,255,255,0.25), transparent);
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse-glow 2s ease-in-out infinite;
}
.ring-center i {
  font-size: 16px;
  color: #00ffff;
  filter: drop-shadow(0 0 6px #00ffff);
}

@keyframes spin { to { transform: rotate(360deg); } }
@keyframes pulse-glow {
  0%, 100% { box-shadow: 0 0 8px rgba(0,255,255,0.2); }
  50% { box-shadow: 0 0 24px rgba(0,255,255,0.5); }
}

.loading-text {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
.loading-label {
  font-size: 14px;
  font-weight: 700;
  color: #00ffff;
  letter-spacing: 3px;
  text-shadow: 0 0 10px rgba(0, 255, 255, 0.5);
  text-transform: uppercase;
}
.loading-bar {
  width: 140px;
  height: 3px;
  background: rgba(0, 255, 255, 0.15);
  border-radius: 2px;
  overflow: hidden;
}
.loading-bar-inner {
  width: 40%;
  height: 100%;
  background: linear-gradient(90deg, transparent, #00ffff, transparent);
  border-radius: 2px;
  animation: bar-slide 1.6s ease-in-out infinite;
}
@keyframes bar-slide {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(350%); }
}

.loading-desc {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  letter-spacing: 0.5px;
}
</style>
