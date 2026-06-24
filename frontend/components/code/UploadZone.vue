<template>
  <div
    class="upload-zone"
    @click="$emit('click')"
    @dragover.prevent
    @drop.prevent="$emit('drop', $event)"
  >
    <i class="fa-solid fa-cloud-arrow-up upload-icon"></i>
    <div class="upload-text">
      <span class="upload-text-bold">拖拽图片到这里，或点击上传</span>
    </div>
    <div class="upload-hint">
      {{ hint }}
    </div>
  </div>
</template>

<script setup>
defineProps({
  hint: { type: String, default: '支持 PNG、JPG、JPEG，最多 5 张，单张不超过 10MB' },
})
defineEmits(['click', 'drop'])
</script>

<style scoped>
.upload-zone {
  border: 2px dashed rgba(0, 255, 255, 0.1);
  border-radius: 16px;
  padding: 40px 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.upload-zone:hover {
  border-color: rgba(0, 255, 255, 0.3);
  background: rgba(0, 255, 255, 0.02);
}

.upload-icon {
  font-size: 40px;
  color: rgba(0, 255, 255, 0.25);
  margin-bottom: 12px;
}

.upload-text {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 6px;
}
.upload-text-bold {
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
}

.upload-hint {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}
</style>
