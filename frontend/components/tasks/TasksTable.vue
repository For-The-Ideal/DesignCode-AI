<template>
  <div class="table-wrap">
    <el-table
      :data="paginatedTasks"
      rowClassName="task-row"
      highlightCurrentRow
      style="width:100%"
    >
      <!-- ═══ 项目名称 ═══ -->
      <el-table-column label="项目名称" align="left">
        <template #default="{ row }">
          <div class="cell-proj">
            <div class="proj-icon">
              <el-image
                v-if="row.images?.[0]?.url"
                :src="row.images[0].url"
                fit="cover"
                class="icon-img"
              >
                <template #error><i class="fas fa-image"></i></template>
              </el-image>
              <i v-else class="fas fa-image"></i>
            </div>
            <span class="proj-name">{{ taskDisplayName(row) }}</span>
          </div>
        </template>
      </el-table-column>

      <!-- ═══ 框架 ═══ -->
      <el-table-column label="框架" align="center">
        <template #default="{ row }">
          <el-tag :type="fwTagType(row.framework)" size="small" effect="dark" round>
            {{ row.framework }}
          </el-tag>
        </template>
      </el-table-column>

      <!-- ═══ 平台 ═══ -->
      <el-table-column label="平台" align="center">
        <template #default="{ row }">
          <span class="cell-plat">
            <i :class="platformIcon(row.platform)"></i>
            {{ platformLabel(row.platform) }}
          </span>
        </template>
      </el-table-column>

      <!-- ═══ 状态 ═══ -->
      <el-table-column label="状态" align="center">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small" effect="dark" round>
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <!-- ═══ 进度 ═══ -->
      <el-table-column label="进度" align="center">
        <template #default="{ row }">
          <div class="cell-prog">
            <el-progress
              :percentage="row.progress"
              :status="row.status === 'success' ? 'success' : row.status === 'failed' ? 'exception' : ''"
              :strokeWidth="6"
              :showText="false"
              :color="progressColor(row.status)"
              style="flex:1"
            />
            <span class="prog-num">{{ row.progress }}%</span>
          </div>
        </template>
      </el-table-column>

      <!-- ═══ 创建时间 ═══ -->
      <el-table-column label="创建时间" align="center">
        <template #default="{ row }">
          <span class="cell-time">{{ row.createdAt }}</span>
        </template>
      </el-table-column>

      <!-- ═══ 操作 ═══ -->
      <el-table-column label="操作" width="120" align="center" fixed="right">
        <template #default="{ row }">
          <div class="cell-acts" @click.stop>
            <el-button link type="primary" size="small" @click="$emit('detail', row)">
              <span>详情</span>
            </el-button>
            <el-button link type="danger" size="small" @click="$emit('delete', row)">
              <span>删除</span>
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </div>

  <!-- 空状态 -->
  <div class="empty-place" v-if="tasks.length === 0">
    <el-empty
      :imageSize="80"
      :description="isEmptyFiltered ? '未找到匹配的任务' : '暂无任务，快去创建一个吧'"
    />
  </div>

  <!-- 分页 -->
  <div class="pager-wrap" v-if="tasks.length > pageSize">
    <el-pagination
      v-model:currentPage="page"
      :pageSize="pageSize"
      :total="tasks.length"
      layout="prev, pager, next"
      background
      small
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  platformLabel, platformIcon, statusLabel,
  statusTagType, fwTagType, progressColor, taskDisplayName,
} from '~/utils/taskHelpers'

const props = defineProps({
  tasks: { type: Array, default: () => [] },
  isEmptyFiltered: { type: Boolean, default: false },
})

defineEmits(['detail', 'delete'])

const page = ref(1)
const pageSize = ref(10)

const paginatedTasks = computed(() => {
  const s = (page.value - 1) * pageSize.value
  return props.tasks.slice(s, s + pageSize.value)
})
</script>

<style scoped>
.table-wrap {
  flex: 1;
  min-height: 0;
  border-radius: 14px;
  border: 1px solid rgba(0, 255, 255, 0.08);
  overflow: hidden;
  background: rgba(10, 14, 23, 0.65);
  backdrop-filter: blur(12px);
}

/* ═══ 表格背景强制透明 ═══ */
:deep(.el-table),
:deep(.el-table__inner-wrapper),
:deep(.el-table__body-wrapper),
:deep(.el-table__header-wrapper),
:deep(.el-table__fixed),
:deep(.el-table__fixed-right),
:deep(.el-table__fixed-header-wrapper),
:deep(.el-table__fixed-body-wrapper),
:deep(.el-table__fixed-right-patch),
:deep(.el-table__empty-block),
:deep(.el-table__empty-text),
:deep(.el-table__body),
:deep(.el-table__body tr),
:deep(.el-table__body tr td) {
  background: transparent !important;
  border-bottom: none !important;
}

:deep(.el-table__header th),
:deep(.el-table__fixed-header-wrapper th),
:deep(.el-table.is-scrolling-none th.el-table-fixed-column--right) {
  background: radial-gradient(circle at top, rgba(41,70,255,.08), transparent 35%), #090B12;
  color: #94a3b8 !important;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.3px;
  border-bottom: 1px solid rgba(0,255,255,0.08) !important;
  padding: 14px 8px;
}

:deep(.task-row),
:deep(.task-row td) {
  background: transparent !important;
  border-bottom: 1px solid rgba(0,255,255,0.04) !important;
  color: #cbd5e1;
  font-size: 13px;
  padding: 12px 0;
}

:deep(.el-table__body tr:hover > td) {
  background: rgba(0, 255, 255, 0.03) !important;
}
:deep(.el-table__body tr.current-row > td) {
  background: rgba(0, 255, 255, 0.05) !important;
}

:deep(.el-table::before),
:deep(.el-table::after),
:deep(.el-table__inner-wrapper::before),
:deep(.el-table__inner-wrapper::after) {
  display: none !important;
}

/* ═══ 单元格内容 ═══ */
.cell-proj  { display: flex; align-items: center; gap: 10px; }
.proj-icon {
  width: 34px; height: 34px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(135deg, rgba(0,255,255,0.08), rgba(255,0,255,0.04));
  border: 1px solid rgba(0,255,255,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}
.proj-icon i { font-size: 13px; color: rgba(255,255,255,0.15); }
.icon-img { width: 100%; height: 100%; }
.proj-name {
  font-size: 14px;
  font-weight: 500;
  color: #e2e8f0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cell-plat {
  font-size: 12px;
  color: #94a3b8;
  display: flex;
  align-items: center;
  gap: 5px;
  justify-content: center;
}
.cell-plat i { font-size: 11px; opacity: 0.5; }

.cell-prog {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: center;
  flex-wrap: nowrap;
}
.prog-num {
  font-size: 11px;
  color: #94a3b8;
  font-weight: 600;
  flex-shrink: 0;
  white-space: nowrap;
}

.cell-time { font-size: 12px; color: #94a3b8; font-family: 'Fira Code', monospace; }

.cell-acts {
  display: flex;
  gap: 8px;
  justify-content: center;
  align-items: center;
}

/* ═══ 空状态 ═══ */
.empty-place {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

/* ═══ 分页 ═══ */
.pager-wrap {
  display: flex;
  justify-content: center;
  margin-top: 14px;
  flex-shrink: 0;
}

/* ═══ 响应式 ═══ */
@media (max-width: 1100px) {
  :deep(.f-sel) { width: 110px; }
}
</style>
