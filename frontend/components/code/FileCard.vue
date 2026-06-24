<template>
  <div class="file-card">
    <i class="fa-regular fa-file-image file-icon"></i>
    <div class="file-info">
      <div class="file-name">
        {{ file.name }}
        <span v-if="file.index" class="file-index">{{ file.index }}</span>
      </div>
      <div v-if="file.desc" class="file-desc">{{ file.desc }}</div>
      <div class="file-meta">{{ file.dimensions || '1920×1080' }}</div>
    </div>
    <button v-if="$slots.actions" class="file-actions">
      <slot name="actions" />
    </button>
  </div>
</template>

<script setup>
defineProps({
  file: { type: Object, required: true },
})
</script>

<style scoped>
.file-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 12px;
  border: 1px solid rgba(0, 255, 255, 0.06);
  background: rgba(0, 255, 255, 0.02);
  transition: all 0.2s;
  cursor: pointer;
}
.file-card:hover {
  border-color: rgba(0, 255, 255, 0.15);
  background: rgba(0, 255, 255, 0.04);
}

.file-icon {
  font-size: 22px;
  color: #00ffff;
  margin-top: 2px;
  flex-shrink: 0;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.85);
}
.file-index {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  margin-left: 4px;
}

.file-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  margin-top: 2px;
}

.file-meta {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  margin-top: 4px;
}

.file-actions {
  flex-shrink: 0;
}
</style>
