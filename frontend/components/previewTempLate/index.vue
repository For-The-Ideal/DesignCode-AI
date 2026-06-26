<template>
  <div class="preview-template">
    <!-- 预览 header -->
    <div class="preview-header">
      <div class="preview-title">
        <i class="fas fa-mobile-alt"></i>
        <span>实时预览</span>
      </div>
      <div class="device-tabs">
        <div class="tab-indicator" :class="device"></div>
        <button
          v-for="d in deviceOptions"
          :key="d.id"
          class="device-tab"
          :class="{ active: device === d.id }"
          @click="device = d.id"
        >
          {{ d.label }}
        </button>
      </div>
    </div>

    <!-- 预览区 -->
    <div class="preview-body" :class="'preview--' + device">
      <MobilePreview
        v-if="device === 'mobile'"
        :html="html"
        :showBottomNav="showBottomNav"
        :bottomNav="bottomNav"
      />
      <ComputerPreview
        v-else
        :html="html"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import MobilePreview from './mobile.vue'
import ComputerPreview from './computer.vue'

const deviceOptions = [
  { id: 'mobile', label: '手机' },
  { id: 'desktop', label: '电脑' },
]

const device = ref('mobile')


defineProps({
  html: { type: String, default: '' },
  showBottomNav: { type: Boolean, default: true },
  bottomNav: { type: Array, default: () => [
    { icon: 'fas fa-home', label: '首页', active: true },
    { icon: 'fas fa-compass', label: '发现', active: false },
    { icon: 'fas fa-shopping-bag', label: '购物袋', active: false },
    { icon: 'fas fa-user-circle', label: '我的', active: false },
  ]},
})
</script>

<style scoped>
.preview-template {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* 预览 header */
.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px 18px;
  flex-shrink: 0;
  background: rgba(0, 0, 0, 0.3);
  border-bottom: 1px solid rgba(0, 255, 255, 0.08);
}

.preview-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: #c8d6e5;
  letter-spacing: 0.3px;
}
.preview-title i {
  font-size: 14px;
  color: #00cfff;
  opacity: 0.8;
}

.device-tabs {
  position: relative;
  display: flex;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 20px;
  padding: 5px;
  width: 125px;
}

.tab-indicator {
  position: absolute;
  top: 3px;
  bottom: 3px;
  width: calc(50% - 3px);
  background: rgba(0, 255, 255, 0.15);
  border-radius: 18px;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.tab-indicator.mobile { transform: translateX(0); }
.tab-indicator.desktop { transform: translateX(calc(100% + 3px)); }

.device-tab {
  position: relative;
  z-index: 1;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2px 0;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.35);
  font-size: 12px;
  cursor: pointer;
  transition: color 0.2s;
  font-family: inherit;
}

.device-tab:hover {
  color: rgba(255, 255, 255, 0.6);
}

.device-tab.active {
  color: #00ffff;
}

/* 预览体 */
.preview-body {
  flex: 1;
  min-height: 0;
  display: flex;
  justify-content: center;
  overflow: hidden;
}

.preview--mobile {
  align-items: flex-start;
  padding: 16px 0;
}
.preview--desktop {
  align-items: flex-start;
  padding: 16px;
}
</style>
