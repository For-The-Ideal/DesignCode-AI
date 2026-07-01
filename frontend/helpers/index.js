/**
 * taskHelpers.js — 任务模块公共常量与工具函数
 *
 * 使用位置：tasks.vue / TasksTable.vue / TasksFilter.vue / TasksSidebar.vue / detail/[id].vue
 */

// ── 平台映射 ──
export const PLAT = { mobile: '手机', desktop: '桌面', tablet: '平板' }
export const PLAT_ICON = { mobile: 'fas fa-mobile-alt', desktop: 'fas fa-desktop', tablet: 'fas fa-tablet-alt' }

// ── 状态映射 ──
export const STATUS_MAP = { pending: '排队中', running: '生成中', success: '已完成', failed: '失败' }

export const platformLabel = p => PLAT[p] || p || '—'
export const platformIcon = p => PLAT_ICON[p] || 'fas fa-question-circle'
export const statusLabel  = s => STATUS_MAP[s] || s

export const statusTagType = s =>
  ({ success: 'success', running: 'warning', pending: 'info', failed: 'danger' }[s] || 'info')

export const fwTagType = fw =>
  ({ Flutter: 'primary', React: '', Vue3: 'success' }[fw] || 'info')

export const progressColor = s =>
  ({ success: '#34d399', failed: '#f87171', running: '#facc15' }[s])

/**
 * 任务显示名称：优先第一张图片描述 → 平台·框架 → id
 */
export const taskDisplayName = t => t.images?.[0]?.desc || `${PLAT[t.platform] || t.platform || '未知'} · ${t.framework || '未知'}`
