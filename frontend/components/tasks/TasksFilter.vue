<template>
  <div class="filter-bar">
    <!-- 搜索 -->
    <div class="fb-search">
      <el-input
        :modelValue="modelValue.search"
        @update:modelValue="patch('search', $event)"
        placeholder="搜索任务名称、框架…"
        :prefix-icon="Search"
        clearable
      />
    </div>

    <!-- 框架 -->
    <div class="fb-field">
      <span class="fb-label">框架</span>
      <el-select
        class="fb-select"
        :modelValue="modelValue.framework"
        @update:modelValue="patch('framework', $event)"
        placeholder="全部"
        clearable
        popper-class="task-popper"
      >
        <el-option v-for="fw in frameworks" :key="fw.value" :label="fw.label" :value="fw.value" />
      </el-select>
    </div>

    <!-- 平台 -->
    <div class="fb-field">
      <span class="fb-label">平台</span>
      <el-select
        class="fb-select"
        :modelValue="modelValue.platform"
        @update:modelValue="patch('platform', $event)"
        placeholder="全部"
        clearable
        popper-class="task-popper"
      >
        <el-option v-for="pf in platforms" :key="pf.value" :label="pf.label" :value="pf.value" />
      </el-select>
    </div>

    <!-- 状态 -->
    <div class="fb-field">
      <span class="fb-label">状态</span>
      <el-select
        class="fb-select"
        :modelValue="modelValue.status"
        @update:modelValue="patch('status', $event)"
        placeholder="全部"
        clearable
        popper-class="task-popper"
      >
        <el-option v-for="s in statusOptions" :key="s.value" :label="s.label" :value="s.value">
          <span class="opt-dot" :style="{ background: s.color }"></span>{{ s.label }}
        </el-option>
      </el-select>
    </div>

    <!-- 重置：仅在有筛选时出现 -->
    <button v-if="hasFilter" class="fb-reset" @click="reset">
      <i class="fas fa-rotate-left"></i> 重置
    </button>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { storeToRefs } from 'pinia'
import { useConfigStore } from '~/stores/config'

const props = defineProps({
  modelValue: { type: Object, required: true },
})
const emit = defineEmits(['update:modelValue'])

// 框架/平台选项与生成面板同源（/config/options）
const configStore = useConfigStore()
const { frameworks, platforms } = storeToRefs(configStore)
onMounted(() => configStore.load())

// 状态为任务生命周期枚举，附带状态色点
const statusOptions = [
  { label: '已完成', value: 'success', color: '#34d399' },
  { label: '生成中', value: 'running', color: '#22d3ee' },
  { label: '排队中', value: 'pending', color: '#60a5fa' },
  { label: '失败', value: 'failed', color: '#f87171' },
]

const patch = (key, val) => emit('update:modelValue', { ...props.modelValue, [key]: val })
const reset = () => emit('update:modelValue', {
  ...props.modelValue, search: '', framework: '', platform: '', status: '',
})
const hasFilter = computed(() => {
  const m = props.modelValue || {}
  return !!(m.search || m.framework || m.platform || m.status)
})
</script>

<style scoped>
.filter-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  flex-shrink: 0;
  margin-bottom: 16px;
  padding: 12px 14px;
  background: rgba(10, 14, 23, 0.55);
  border: 1px solid rgba(0, 255, 255, 0.08);
  border-radius: 14px;
  backdrop-filter: blur(10px);
}

.fb-search { flex: 1 1 220px; min-width: 200px; max-width: 340px; }
.fb-field { display: flex; align-items: center; gap: 8px; }
.fb-label {
  font-size: 12px; font-weight: 600; color: #64748b;
  letter-spacing: 0.5px; white-space: nowrap;
}
.fb-select { width: 132px; }

/* 统一控件外观：搜索框 + 下拉触发器 */
.fb-search :deep(.el-input__wrapper),
.fb-select :deep(.el-select__wrapper) {
  min-height: 38px;
  background: rgba(15, 22, 32, 0.85) !important;
  border: 1px solid rgba(0, 255, 255, 0.12) !important;
  border-radius: 10px !important;
  box-shadow: none !important;
  transition: border-color 0.18s ease, box-shadow 0.18s ease;
}
.fb-search :deep(.el-input__wrapper:hover),
.fb-select :deep(.el-select__wrapper:hover) {
  border-color: rgba(0, 255, 255, 0.30) !important;
}
.fb-search :deep(.el-input__wrapper.is-focus),
.fb-select :deep(.el-select__wrapper.is-focused) {
  border-color: rgba(34, 211, 238, 0.55) !important;
  box-shadow: 0 0 0 3px rgba(34, 211, 238, 0.12) !important;
}

/* 文本 / 占位 / 图标配色 */
.fb-search :deep(.el-input__inner),
.fb-select :deep(.el-select__selected-item),
.fb-select :deep(.el-select__placeholder) {
  color: #e2e8f0 !important;
  font-size: 13px;
}
.fb-search :deep(.el-input__inner::placeholder),
.fb-select :deep(.el-select__placeholder.is-transparent) {
  color: #64748b !important;
}
.fb-search :deep(.el-input__prefix) { color: rgba(34, 211, 238, 0.5); }
.fb-select :deep(.el-select__caret) { color: rgba(34, 211, 238, 0.55) !important; }

/* 重置按钮 */
.fb-reset {
  display: inline-flex; align-items: center; gap: 6px;
  height: 38px; padding: 0 14px; margin-left: auto;
  font-size: 13px; color: #94a3b8;
  background: transparent;
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 10px; cursor: pointer;
  transition: all 0.18s ease;
}
.fb-reset:hover {
  color: #f87171;
  border-color: rgba(248, 113, 113, 0.4);
  background: rgba(248, 113, 113, 0.08);
}

@media (max-width: 1100px) {
  .fb-select { width: 118px; }
}
</style>
