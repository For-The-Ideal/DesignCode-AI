<template>
  <div class="flutter-mock">
    <div class="flutter-screen">
      <!-- 固定：刘海区 -->
      <div class="flutter-notch-area">
        <span class="notch-time">{{ currentTime }}</span>
        <div class="notch-cutout"></div>
        <div class="notch-icons">
          <i class="fas fa-signal"></i>
          <i class="fas fa-wifi"></i>
          <i class="fas fa-battery-full"></i>
        </div>
      </div>

      <!-- 动态：页面内容 or 加载中 -->
      <div v-if="html" class="flutter-page-content" v-html="html"></div>
      <div v-else class="flutter-page-content flutter-loading">
        <!-- 扫描线 -->
        <div class="scan-line"></div>
        <!-- 环形脉冲 -->
        <div class="ring-pulse">
          <div class="ring ring-1"></div>
          <div class="ring ring-2"></div>
          <div class="ring ring-3"></div>
          <div class="ring-center">
            <i class="fas fa-microchip"></i>
          </div>
        </div>
        <!-- 文字 -->
        <div class="loading-text">
          <span class="loading-label">AI 生成引擎</span>
          <div class="loading-bar">
            <div class="loading-bar-inner"></div>
          </div>
          <span class="loading-desc">正在解析设计稿，构建代码框架...</span>
        </div>
      </div>

      <!-- 固定：底部导航 -->
      <div class="flutter-bottom-nav" v-if="showBottomNav">
        <div v-for="tab in bottomNav" :key="tab.label" class="tab-item" :class="{ active: tab.active }">
          <i :class="tab.icon"></i>
          <span>{{ tab.label }}</span>
        </div>
      </div>
    </div>
    <!-- 固定：Home 指示条 -->
    <div class="flutter-home-bar"></div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const currentTime = computed(() => {
  const now = new Date()
  return now.getHours().toString().padStart(2, '0') + ':' + now.getMinutes().toString().padStart(2, '0')
})

const props = defineProps({
  html: { type: String, default: '' },
  showBottomNav: { type: Boolean, default: true },
  bottomNav: {
    type: Array,
    default: () => [
      { icon: 'fas fa-home', label: '首页', active: true },
      { icon: 'fas fa-compass', label: '发现', active: false },
      { icon: 'fas fa-shopping-bag', label: '购物袋', active: false },
      { icon: 'fas fa-user-circle', label: '我的', active: false },
    ],
  },
})
</script>

<style lang="scss" scoped>
/* ── 外壳 ── */
.flutter-mock {
  width: 100%;
  max-width: 375px;
  background: #010101;
  border-radius: 48px;
  padding: 10px 6px;
  box-shadow:
    inset 0 0 2px rgba(255,255,255,0.15),
    0 0 0 2px #1a1a1a,
    0 0 0 4px #0d0d0d,
    0 0 0 6px #1a1a1a,
    0 24px 48px rgba(0,0,0,0.55);
  position: relative;
}

/* ── 屏幕 ── */
.flutter-screen {
  background: #f2f2f7;
  border-radius: 38px;
  overflow: hidden;
  height: 660px;
  display: flex;
  flex-direction: column;
  position: relative;
}

/* ── 刘海 ── */
.flutter-notch-area {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px 4px;
  background: #f2f2f7;
  position: relative;
  min-height: 42px;
  flex-shrink: 0;
}
.notch-time {
  font-size: 13px;
  font-weight: 700;
  color: #000;
  z-index: 1;
}
.notch-icons {
  display: flex;
  gap: 5px;
  font-size: 11px;
  color: #000;
  z-index: 1;
}
.notch-cutout {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 126px;
  height: 28px;
  background: #010101;
  border-radius: 0 0 18px 18px;
}

/* ── 动态内容区 ── */
.flutter-page-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* ── 底部导航 ── */
.flutter-bottom-nav {
  display: flex;
  justify-content: space-around;
  background: #fff;
  border-top: 1px solid #e5e5ea;
  padding: 6px 0 10px;
  flex-shrink: 0;
}
.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  color: #8e8e93;
  font-size: 10px;
  transition: 0.2s;
}
.tab-item i { font-size: 18px; }
.tab-item.active { color: #1c1c1e; }

/* ── Home 指示条 ── */
.flutter-home-bar {
  display: flex;
  justify-content: center;
  padding: 10px 0 6px;
}
.flutter-home-bar::after {
  content: '';
  width: 120px;
  height: 5px;
  border-radius: 3px;
  background: rgba(255,255,255,0.2);
}

@media (max-width: 968px) {
  .flutter-mock { max-width: 320px; }
}

/* ── 赛博加载态 ── */
.flutter-loading {
  position: relative;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 28px;
  background: linear-gradient(160deg, #050a18 0%, #0a1628 40%, #0d1f3c 100%);
  overflow: hidden;
}

/* 扫描线 */
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

/* 环形脉冲 */
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
  width: 100px;
  height: 100px;
  border-top-color: rgba(0, 255, 255, 0.7);
  border-right-color: rgba(255, 0, 255, 0.5);
  animation: spin 1.5s linear infinite;
}
.ring-2 {
  width: 76px;
  height: 76px;
  border-bottom-color: rgba(0, 255, 255, 0.5);
  border-left-color: rgba(255, 0, 255, 0.3);
  animation: spin 2s linear infinite reverse;
}
.ring-3 {
  width: 52px;
  height: 52px;
  border-top-color: rgba(0, 255, 255, 0.3);
  animation: spin 1.2s linear infinite;
}
.ring-center {
  width: 36px;
  height: 36px;
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

@keyframes spin {
  to { transform: rotate(360deg); }
}

@keyframes pulse-glow {
  0%, 100% { box-shadow: 0 0 8px rgba(0,255,255,0.2); }
  50% { box-shadow: 0 0 24px rgba(0,255,255,0.5); }
}

/* 文字区域 */
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
