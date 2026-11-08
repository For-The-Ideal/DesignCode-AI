<template>
  <!-- 画布工具条：模式切换 / 子级提示 / AI 描述 -->
  <div class="wb-toolbar">
    <div class="tb-left">
      <div class="mode-seg">
        <button class="seg" :class="{ active: mode === 'draw' || mode === 'child' }" @click="emit('setMode', 'draw')">
          <i class="fas fa-draw-polygon"></i>
          框选
        </button>
        <button class="seg" :class="{ active: mode === 'select' }" @click="emit('setMode', 'select')">
          <i class="fas fa-mouse-pointer"></i>
          编辑
        </button>
      </div>
      <span v-if="mode === 'child'" class="child-hint">
        为「{{ parentName || '所选区域' }}」框选子区域
        <button class="link-btn" @click="emit('exitChild')">退出</button>
      </span>
    </div>

    <div class="tb-right">
      <!-- 唯一的确认动作：校验节点 → AI 回填描述 → 保存标注 -->
      <button
        class="tb-btn tb-btn--primary"
        :class="{ 'is-loading': aiLoading }"
        :disabled="!hasNodes || aiLoading"
        title="校验节点后 AI 回填描述，并保存当前标注"
        @click="emit('ai')"
      >
        <i class="fas" :class="aiLoading ? 'fa-spinner fa-spin' : 'fa-magic'"></i>
        {{ aiLoading ? 'AI 描述中…' : 'AI 描述' }}
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  // draw | select | child
  mode: { type: String, default: 'draw' },
  parentName: { type: String, default: '' },
  hasNodes: { type: Boolean, default: false },
  aiLoading: { type: Boolean, default: false },
})
const emit = defineEmits(['setMode', 'exitChild', 'ai'])
</script>

<style scoped>
.wb-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  padding: 6px 10px;
  border-radius: 12px;
  border: 1px solid rgba(0, 255, 255, 0.1);
  background: rgba(10, 14, 23, 0.85);
}
.tb-left { display: flex; align-items: center; gap: 10px; min-width: 0; }
.mode-seg {
  display: inline-flex;
  gap: 4px;
  padding: 3px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
.seg {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.55);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.seg:hover { color: rgba(255, 255, 255, 0.85); }
.seg.active { background: rgba(0, 255, 255, 0.12); color: #00ffff; }

.child-hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: rgba(255, 207, 77, 0.9);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.link-btn {
  border: none;
  background: none;
  color: #ff6b7a;
  font-size: 12px;
  cursor: pointer;
  padding: 0;
}
.link-btn:hover { text-decoration: underline; }

.tb-right { display: flex; align-items: center; gap: 10px; }
.tb-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.3);
  color: rgba(255, 255, 255, 0.7);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}
.tb-btn:hover:not(:disabled) { border-color: rgba(0, 255, 255, 0.5); color: #00ffff; }
.tb-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.tb-btn--primary {
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  color: #fff;
  border: none;
}
.tb-btn--primary:hover:not(:disabled) { filter: brightness(1.08); color: #fff; }
.tb-btn.is-loading:disabled {
  position: relative;
  overflow: hidden;
  opacity: 1;
  cursor: progress;
  border-color: rgba(0, 255, 255, 0.7);
  color: #00ffff;
}
.tb-btn.is-loading::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(100deg, transparent 20%, rgba(0, 255, 255, 0.25) 50%, transparent 80%);
  transform: translateX(-100%);
  animation: tb-sweep 1.1s linear infinite;
  pointer-events: none;
}
@keyframes tb-sweep { to { transform: translateX(100%); } }
</style>
