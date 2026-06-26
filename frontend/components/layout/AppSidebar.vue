<template>
  <aside class="app-sidebar" :class="{ collapsed: !sidebarOpen }">
    <!-- 头部：品牌名 + 展开/收起 -->
    <div class="sidebar-head">
      <span v-show="sidebarOpen" class="sidebar-brand">{{ brand }}</span>
      <button
        class="toggle-btn"
        @click="toggleSidebar"
        :title="sidebarOpen ? '收起菜单' : '展开菜单'"
      >
        <i :class="sidebarOpen ? 'fas fa-chevron-left' : 'fas fa-chevron-right'"></i>
      </button>
    </div>

    <!-- 导航菜单（收起时仅图标） -->
    <nav class="sidebar-nav">
      <a
        v-for="item in navItems"
        :key="item.label"
        href="#"
        class="nav-item"
        :class="{ active: item.active }"
        @click.prevent="$emit('nav-click', item)"
      >
        <i :class="item.icon" class="nav-icon"></i>
        <span v-show="sidebarOpen" class="nav-label">{{ item.label }}</span>
      </a>
    </nav>

    <div v-show="sidebarOpen" class="sidebar-slot">
        <!-- 分割线 -->
        <div class="sidebar-divider"></div>
        <!-- 内容插槽 -->
        <slot  />
    </div>
   
  </aside>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  brand: { type: String, default: '' },
  navItems: { type: Array, default: () => [] },
  modelValue: { type: Boolean, default: true },
  expandedWidth: { type: String, default: '280px' },
})

const emit = defineEmits(['update:modelValue', 'nav-click'])

const sidebarOpen = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  sidebarOpen.value = val
})

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
  emit('update:modelValue', sidebarOpen.value)
}
</script>

<style scoped>
.app-sidebar {
  width: v-bind('props.expandedWidth');
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  padding: 16px 12px;
  background: rgba(10, 10, 15, 0.95);
  border-right: 1px solid rgba(0, 255, 255, 0.06);
  transition: width 0.3s ease;
  position: relative;
}
.app-sidebar.collapsed {
  width: 56px;
  padding: 16px 8px;
}

/* 头部 */
.sidebar-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  flex-shrink: 0;
}
.collapsed .sidebar-head {
  justify-content: center;
}
.sidebar-brand {
  font-size: 14px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.6);
  letter-spacing: 0.5px;
}

/* 折叠按钮 */
.toggle-btn {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 255, 255, 0.06);
  color: #00cfff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  transition: all 0.2s;
  flex-shrink: 0;
}
.toggle-btn:hover {
  background: rgba(0, 255, 255, 0.14);
  border-color: #00ffff;
}

/* 导航 */
.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
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
  white-space: nowrap;
}
.collapsed .nav-item {
  justify-content: center;
  padding: 10px 0;
}
.nav-item:hover {
  color: rgba(255, 255, 255, 0.7);
  background: rgba(255, 255, 255, 0.04);
}
.nav-item.active {
  color: #e2e8f0;
  background: rgba(0, 255, 255, 0.06);
}
.nav-icon {
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}
.nav-item.active .nav-icon { color: #00ffff; }
.nav-item:not(.active) .nav-icon { color: rgba(255, 255, 255, 0.35); }
.nav-label {
  overflow: hidden;
}

.sidebar-divider {
  height: 1px;
  background: rgba(0, 255, 255, 0.08);
  margin: 12px 0;
  flex-shrink: 0;
}

/* 插槽容器：填满剩余空间，支持滚动 */
.sidebar-slot {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow-y: auto;
}
</style>
