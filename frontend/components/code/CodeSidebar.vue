<template>
  <aside class="code-sidebar">
    <!-- 导航菜单 -->
    <nav class="sidebar-nav">
      <NuxtLink
        v-for="item in navItems"
        :key="item.label"
        :to="item.to"
        class="nav-item"
        :class="{ 'nav-item--active': item.active }"
      >
        <i :class="item.icon" class="nav-icon"></i>
        {{ item.label }}
      </NuxtLink>
    </nav>

    <div class="sidebar-divider"></div>

    <!-- 底部区域 -->
    <div class="sidebar-bottom">
      <!-- 本月使用情况（SVG 圆环） -->
      <div class="usage-card">
        <div class="usage-label">本月使用情况</div>
        <div class="usage-row">
          <div class="usage-ring-wrap">
            <svg class="usage-ring" viewBox="0 0 80 80">
              <defs>
                <linearGradient id="ringGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#00ffff" />
                  <stop offset="100%" stop-color="#ff00ff" />
                </linearGradient>
              </defs>
              <circle cx="40" cy="40" r="34" class="ring-bg" />
              <circle
                cx="40" cy="40" r="34"
                class="ring-fill"
                :style="ringStyle"
              />
            </svg>
            <div class="ring-center">
              <span class="ring-value">{{ usagePercent }}%</span>
            </div>
          </div>
          <div class="usage-info">
            <span class="usage-remaining">
              剩余 <span class="usage-remaining-num">{{ usageRemaining }}</span> 次
            </span>
            <p class="usage-total">共 {{ usageTotal }} 次</p>
          </div>
        </div>
      </div>

      <!-- 升级会员 -->
      <div class="upgrade-card">
        <div class="upgrade-glow"></div>
        <div class="upgrade-title">💎 升级会员</div>
        <div class="upgrade-desc">解锁更多高级功能</div>
        <button class="upgrade-btn">立即升级</button>
      </div>
    </div>
  </aside>
</template>

<script setup>
const navItems = [
  { icon: 'fa-solid fa-code', label: '代码生成', active: true, to: '/code' },
  { icon: 'fa-regular fa-copy', label: '模板市场', to: '#' },
  { icon: 'fa-regular fa-folder', label: '我的项目', to: '#' },
]

const usagePercent = 68
const usageTotal = 100
const usageRemaining = 32

const ringStyle = computed(() => {
  const r = 34
  const circumference = 2 * Math.PI * r
  const offset = circumference * (1 - usagePercent / 100)
  return {
    strokeDasharray: `${circumference}`,
    strokeDashoffset: `${offset}`,
  }
})
</script>

<style scoped>
.code-sidebar {
  width: 225px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  padding: 16px 12px;
  background: rgba(10, 10, 15, 0.95);
  border-right: 1px solid rgba(0, 255, 255, 0.06);
}

/* 导航 */
.sidebar-nav {
  flex: 0;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.4);
  text-decoration: none;
  transition: all 0.2s;
  margin-bottom: 2px;
}
.nav-item:hover {
  color: rgba(255, 255, 255, 0.7);
  background: rgba(255, 255, 255, 0.04);
}
.nav-item--active {
  color: #e2e8f0;
  background: rgba(0, 255, 255, 0.06);
}
.nav-icon {
  width: 20px;
  text-align: center;
}
.nav-item--active .nav-icon {
  color: #00ffff;
}
.nav-item:not(.nav-item--active) .nav-icon {
  color: rgba(255, 255, 255, 0.35);
}

.sidebar-divider {
  height: 1px;
  background: rgba(0, 255, 255, 0.06);
  margin: 12px 0;
}

/* 底部 */
.sidebar-bottom {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* ═══ 使用量卡片（SVG 圆环） ═══ */
.usage-card {
  background: rgba(0, 255, 255, 0.03);
  border-radius: 14px;
  padding: 16px;
  border: 1px solid rgba(0, 255, 255, 0.08);
}
.usage-label {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  margin-bottom: 12px;
}

.usage-ring-wrap {
  position: relative;
  flex-shrink: 0;
}
.usage-row {
  display: flex;
  align-items: center;
  gap: 2px;
}
.usage-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  gap: 5px;
  margin-left: 12px;
}
.usage-ring {
  width: 80px;
  height: 80px;
  transform: rotate(-90deg);
}
.ring-bg {
  fill: none;
  stroke: rgba(255, 255, 255, 0.05);
  stroke-width: 4;
}
.ring-fill {
  fill: none;
  stroke: url(#ringGradient);
  stroke-width: 4;
  stroke-linecap: round;
  transition: stroke-dashoffset 0.6s ease;
}
.ring-center {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.ring-value {
  font-size: 18px;
  font-weight: 700;
  color: #00ffff;
}

.usage-remaining {
  font-size: 12px;
  color: #fff;
  text-align: left;
}
.usage-remaining-num {
  font-weight: 500;
  margin: 0px 3px;
}


.usage-total {
  font-size: 12px;
  color: #ccc;
}

/* ═══ 升级卡片（主题色版） ═══ */
.upgrade-card {
  position: relative;
  border-radius: 14px;
  padding: 18px 16px;
  text-align: center;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.04), rgba(255, 0, 255, 0.04));
  border: 1px solid rgba(0, 255, 255, 0.1);
}
.upgrade-glow {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(
    circle at 50% 50%,
    rgba(0, 255, 255, 0.04) 0%,
    transparent 60%
  );
  pointer-events: none;
}
.upgrade-title {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  position: relative;
}
.upgrade-desc {
  font-size: 12px;
  color: #fff;
  margin-top: 2px;
  position: relative;
}
.upgrade-btn {
  margin-top: 10px;
  width: 100%;
  padding: 8px 0;
  border-radius: 8px;
  border: none;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 2px 16px rgba(0, 255, 255, 0.15);
  transition: all 0.3s;
  position: relative;
  background: linear-gradient(90deg, #60a5fa, #818cf8);
}
.upgrade-btn:hover {
  box-shadow: 0 4px 28px rgba(0, 255, 255, 0.35);
  transform: translateY(-1px);
}

@media (max-width: 768px) {
  .code-sidebar {
    display: none;
  }
}
</style>