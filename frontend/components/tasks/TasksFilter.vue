<template>
  <el-form :model="modelValue" inline class="filter-form">
    <el-form-item label="搜索">
      <el-input
        :modelValue="modelValue.search"
        @update:modelValue="$emit('update:modelValue', { ...modelValue, search: $event })"
        placeholder="任务名称、框架..."
        :prefixIcon="Search"
        clearable
        class="f-search"
      />
    </el-form-item>
    <el-form-item label="框架">
      <el-select
        :modelValue="modelValue.framework"
        @update:modelValue="$emit('update:modelValue', { ...modelValue, framework: $event })"
        placeholder="全部"
        clearable
        class="f-sel"
        popperClass="task-popper"
      >
        <el-option label="Flutter" value="Flutter" />
        <el-option label="React" value="React" />
        <el-option label="Vue3" value="Vue3" />
      </el-select>
    </el-form-item>
    <el-form-item label="平台">
      <el-select
        :modelValue="modelValue.platform"
        @update:modelValue="$emit('update:modelValue', { ...modelValue, platform: $event })"
        placeholder="全部"
        clearable
        class="f-sel"
        popperClass="task-popper"
      >
        <el-option label="📱 手机" value="mobile" />
        <el-option label="🖥️ 桌面" value="desktop" />
        <el-option label="📋 平板" value="tablet" />
      </el-select>
    </el-form-item>
    <el-form-item label="状态">
      <el-select
        :modelValue="modelValue.status"
        @update:modelValue="$emit('update:modelValue', { ...modelValue, status: $event })"
        placeholder="全部"
        clearable
        class="f-sel"
        popperClass="task-popper"
      >
        <el-option label="✅ 已完成" value="success" />
        <el-option label="🔄 生成中" value="running" />
        <el-option label="⏳ 排队中" value="pending" />
        <el-option label="❌ 失败" value="failed" />
      </el-select>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { Search } from '@element-plus/icons-vue'

defineProps({
  modelValue: { type: Object, required: true },
})

defineEmits(['update:modelValue'])
</script>

<style scoped>
.filter-form {
  flex-shrink: 0;
  margin-bottom: 14px;
  display: flex;
  flex-wrap: wrap;
  gap: 4px 12px;
}
.f-search { width: 220px; }
.f-sel    { width: 130px; }

:deep(.f-sel .el-select__wrapper),
:deep(.f-sel .el-input__wrapper),
:deep(.f-search .el-input__wrapper) {
  background: rgba(15, 22, 32, 0.75) !important;
  border: 1px solid rgba(0,255,255,0.12) !important;
  box-shadow: none !important;
}
:deep(.f-sel .el-input__inner) {
  color: #cbd5e1 !important;
}
:deep(.f-sel .el-select__caret) {
  color: rgba(0,255,255,0.4) !important;
}
</style>
