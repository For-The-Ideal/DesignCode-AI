<template>
  <!-- 画布区：有图渲染标注画布，无图渲染默认上传页（点击/拖拽上传） -->
  <div class="wb-canvas">
    <AnnotationCanvas
      v-if="hasImage"
      :mode="mode"
      :show-nav="false"
      :phone-shell="scene === 'single'"
      @drawn="(bbox) => emit('drawn', bbox)"
      @cancelMode="emit('cancelMode')"
    />

    <!-- 默认空态：大上传区（与页面列表/生成配置组成默认工作台） -->
    <div
      v-else
      class="canvas-empty"
      @dragover.prevent
      @drop.prevent="onDrop"
    >
      <div class="empty-icon" :class="{ 'is-loading': uploading, 'is-error': error }">
        <i class="fas" :class="iconClass"></i>
      </div>

      <div class="empty-drop">
        <h3 class="empty-title">{{ emptyTitle }}</h3>
        <button class="empty-btn" :disabled="uploading" @click="pick">
          <i class="fas fa-cloud-arrow-up"></i>
          点击上传图片
        </button>
      </div>

      <p class="empty-hint">{{ emptyHint }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import AnnotationCanvas from '~/components/uiProject/AnnotationCanvas.client.vue'

const props = defineProps({
  // draw | select | child，透传给画布
  mode: { type: String, default: 'draw' },
  // single | multi，仅用于空态文案
  scene: { type: String, default: 'single' },
  hasImage: { type: Boolean, default: false },
  uploading: { type: Boolean, default: false },
  error: { type: Boolean, default: false },
  maxCount: { type: Number, default: 1 },
})
const emit = defineEmits(['drawn', 'cancelMode', 'pick', 'drop'])

const iconClass = computed(() => {
  if (props.uploading) return 'fa-spinner fa-spin'
  if (props.error) return 'fa-circle-exclamation'
  return 'fa-cloud-arrow-up'
})
const emptyTitle = computed(() => {
  if (props.uploading) return '上传中…'
  if (props.error) return '上传失败，请重新上传'
  return props.scene === 'multi' ? '上传多页设计稿' : '上传设计稿'
})
const emptyHint = computed(() =>
  props.scene === 'multi'
    ? `支持一次性上传多张页面设计图（最多 ${props.maxCount} 张），多图需配置页面关联关系`
    : '上传后可框选区域，填写描述或交由 AI 生成',
)

function pick() {
  if (props.uploading) return
  emit('pick')
}
function onDrop(e) {
  const files = Array.from(e.dataTransfer?.files || []).filter((f) => f.type.startsWith('image/'))
  if (files.length) emit('drop', files)
}
</script>

<style scoped>
.wb-canvas {
  position: relative;
  height: 100%;
  min-height: 0;
  border-radius: 14px;
  border: 1px solid rgba(0, 255, 255, 0.08);
  background: rgba(10, 14, 23, 0.7);
  overflow: hidden;
}

/* 空态：居中大上传区 */
.canvas-empty {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 18px;
  padding: 24px;
}

.empty-icon {
  width: 84px;
  height: 84px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 22px;
  border: 1px dashed rgba(0, 255, 255, 0.28);
  background: rgba(0, 255, 255, 0.05);
  color: rgba(0, 255, 255, 0.7);
  font-size: 30px;
}
.empty-icon.is-loading { color: #00ffff; }
.empty-icon.is-error { border-color: rgba(255, 71, 87, 0.4); color: #ff6b7a; }

.empty-drop {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  min-width: min(420px, 80%);
  padding: 22px 32px;
  border-radius: 14px;
  border: 1px dashed rgba(0, 255, 255, 0.3);
  background: rgba(0, 255, 255, 0.03);
  transition: all 0.2s;
  cursor: pointer;
}
.empty-drop:hover {
  border-color: rgba(0, 255, 255, 0.55);
  background: rgba(0, 255, 255, 0.06);
}
.empty-title {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.88);
  letter-spacing: 0.5px;
}
.empty-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 26px;
  border-radius: 10px;
  border: none;
  background: linear-gradient(135deg, #2f6bff, #60a5fa);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}
.empty-btn:hover:not(:disabled) { filter: brightness(1.1); }
.empty-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.empty-hint {
  margin: 0;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.45);
  text-align: center;
}
</style>
