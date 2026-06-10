<template>
  <el-dialog
    v-model="visible"
    :width="width"
    :center="true"
    :show-close="false"
    :close-on-click-modal="closeOnClickModal"
    :close-on-press-escape="closeOnPressEscape"
    class="auth-dialog"
    :before-close="handleBeforeClose">
      <slot />
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue'
const props = defineProps({
  width: { type: String, default: '460px' },
  closeOnClickModal: { type: Boolean, default: true },
  closeOnPressEscape: { type: Boolean, default: true },
})
const emit = defineEmits(['close'])

const visible = ref(false)

const open = () => { visible.value = true }
const close = () => { visible.value = false; emit('close') }
const handleBeforeClose = (done) => { done(); emit('close') }

defineExpose({ open, close, visible })
</script>


<style>
/* ========== DialogModel 暗色弹窗全局样式 ========== */
.auth-dialog {
  background: rgba(12, 20, 28, 0.35) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 255, 255, 0.3) !important;
  border-radius: 32px !important;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}
.auth-dialog .el-dialog__header {
  display: none;
}
.auth-dialog .el-dialog__body {
  padding: 40px 32px;
}

.progress-dialog .el-dialog__header {
  display: none;
}
.progress-dialog .el-dialog__body {
  padding: 0 !important;
}

.desc-dialog .el-dialog__header {
  display: none;
}
.desc-dialog .el-dialog__body {
  padding: 0 !important;
}
</style>