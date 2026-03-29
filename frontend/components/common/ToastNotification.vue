<!-- components/common/ToastNotification.vue -->
<template>
  <Teleport to="body">
    <transition name="toast">
      <div v-if="getVisible" class="toast-container" :class="getToastParams.type">
        <div class="toast-icon">
          <i :class="iconClass"></i>
        </div>
        <div class="toast-content">
          <div class="toast-title">{{ getToastParams.title }}</div>
          <div class="toast-message">{{ getToastParams.message }}</div>
        </div>
        <button class="toast-close" @click="close">
          <i class="fas fa-times"></i>
        </button>
      </div>
    </transition>
  </Teleport>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import {useToastNotificationStore} from "~/stores/toastNotification"
const toastNotificationStore = useToastNotificationStore()

const getToastParams = computed(() => toastNotificationStore.toastParams)
const getVisible = computed(() => toastNotificationStore.visible)
const emit = defineEmits(['close'])

const iconClass = computed(() => {
  const icons = {
    success: 'fas fa-check-circle',
    error: 'fas fa-times-circle',
    warning: 'fas fa-exclamation-triangle',
    info: 'fas fa-info-circle',
    loading: 'fas fa-spinner-circle'
  }
  return icons[getToastParams.value.type] || 'fas fa-check-circle'
})


const close = () => {
 toastNotificationStore.hideToast()
}

</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 80px;
  right: 24px;
  min-width: 320px;
  max-width: 420px;
  background: rgba(10, 20, 30, 0.95);
  backdrop-filter: blur(12px);
  border-radius: 16px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border-left: 4px solid;
  z-index: 1100;
  animation: slideInRight 0.3s ease;
}

.toast-container.success {
  border-left-color: #00ff00;
}

.toast-container.error {
  border-left-color: #ff4444;
}

.toast-container.warning {
  border-left-color: #ffaa00;
}

.toast-container.info {
  border-left-color: #00ffff;
}

.toast-icon i {
  font-size: 24px;
}

.toast-container.success .toast-icon i {
  color: #00ff00;
}

.toast-container.error .toast-icon i {
  color: #ff4444;
}

.toast-container.warning .toast-icon i {
  color: #ffaa00;
}

.toast-container.info .toast-icon i {
  color: #00ffff;
}

.toast-content {
  flex: 1;
}

.toast-title {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 4px;
}

.toast-message {
  font-size: 12px;
  color: #aaa;
}

.toast-close {
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  padding: 4px;
  transition: color 0.3s;
}

.toast-close:hover {
  color: #fff;
}

@keyframes slideInRight {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  transform: translateX(100%);
  opacity: 0;
}
</style>