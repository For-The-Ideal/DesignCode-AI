<template>
  <div class="phone-mock">
    <div class="phone-screen">
      <!-- 固定：刘海区 -->
      <div class="phone-notch-area">
        <span class="notch-time">9:41</span>
        <div class="notch-cutout"></div>
        <div class="notch-icons">
          <i class="fas fa-signal"></i>
          <i class="fas fa-wifi"></i>
          <i class="fas fa-battery-full"></i>
        </div>
      </div>

      <!-- 动态：页面内容 -->
      <div class="phone-page-content" v-html="html"></div>

      <!-- 固定：底部导航 -->
      <div class="phone-bottom-nav">
        <div class="tab-item active"><i class="fas fa-home"></i><span>首页</span></div>
        <div class="tab-item"><i class="fas fa-compass"></i><span>发现</span></div>
        <div class="tab-item"><i class="fas fa-shopping-bag"></i><span>购物袋</span></div>
        <div class="tab-item"><i class="fas fa-user-circle"></i><span>我的</span></div>
      </div>
    </div>
    <!-- 固定：Home 指示条 -->
    <div class="phone-home-bar"></div>
  </div>
</template>

<script setup>
defineProps({
  html: { type: String, default: '' }
})
</script>

<style lang="scss" scoped>
/* ── 外壳 ── */
.phone-mock {
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
.phone-screen {
  background: #f2f2f7;
  border-radius: 38px;
  overflow: hidden;
  height: 660px;
  display: flex;
  flex-direction: column;
  position: relative;
}

/* ── 刘海 ── */
.phone-notch-area {
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
.phone-page-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* ── 底部导航 ── */
.phone-bottom-nav {
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
.phone-home-bar {
  display: flex;
  justify-content: center;
  padding: 10px 0 6px;
}
.phone-home-bar::after {
  content: '';
  width: 120px;
  height: 5px;
  border-radius: 3px;
  background: rgba(255,255,255,0.2);
}

@media (max-width: 968px) {
  .phone-mock { max-width: 320px; }
}
</style>
