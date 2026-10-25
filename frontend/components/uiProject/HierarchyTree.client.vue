<template>
  <div class="tree">
    <div class="tree-head">
      <span class="tree-title">层级结构</span>
      <span class="tree-count">{{ nodeListCount }}</span>
    </div>

    <div class="tree-body">
      <div
        v-for="{ node, depth } in flat"
        :key="node.nodeKey"
        class="tree-row"
        :class="{ 'tree-row--active': store.selectedNodeKey === node.nodeKey }"
        :style="{ paddingLeft: 8 + depth * 18 + 'px' }"
        @click="emit('select', node.nodeKey)"
      >
        <!-- 缩进连线 -->
        <span v-if="depth > 0" class="tree-guide" :style="{ left: 8 + (depth - 1) * 18 + 9 + 'px' }"></span>

        <!-- 展开 / 折叠 -->
        <span v-if="hasKids(node)" class="twisty" @click.stop="toggle(node.nodeKey)">
          {{ collapsedSet.has(node.nodeKey) ? '▸' : '▾' }}
        </span>
        <span v-else class="twisty-spacer"></span>

        <!-- 状态色点 -->
        <span class="dot" :style="{ background: STATUS_COLOR[node.status] || STATUS_COLOR.pending }"></span>

        <span class="node-name">{{ displayName(node) }}</span>

        <!-- 排列方向图标 -->
        <span v-if="node.layout" class="layout-chip" :title="'layout: ' + node.layout">
          {{ LAYOUT_ICON[node.layout] }}
        </span>

        <span class="row-actions">
          <button
            v-if="canAddChild(node)"
            class="mini"
            title="框选子区域"
            @click.stop="emit('addChild', node.nodeKey)"
          >+</button>
          <button class="mini mini--danger" title="删除" @click.stop="emit('deleteNode', node.nodeKey)">
            <i class="fa-regular fa-trash-can"></i>
          </button>
        </span>
      </div>

      <div v-if="!store.hasNodes" class="tree-empty">
        还没有节点，在中间图片上拖拽框选第一个区域
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUiProjectStore } from '~/stores/uiProject'

const store = useUiProjectStore()
const emit = defineEmits(['addChild', 'deleteNode', 'select'])

// 状态 -> 颜色条（灰/黄/绿/红）
const STATUS_COLOR = {
  pending: 'rgba(255,255,255,0.35)',
  ai_draft: '#ffcf4d',
  confirmed: '#37e08d',
  empty: '#ff4757',
}
// 排列方向 -> 图标
const LAYOUT_ICON = { row: '⇄', column: '⇅', grid: '⊞' }

const collapsedSet = ref(new Set())
const nodeListCount = computed(() => store.nodeList.length)

// 嵌套树 -> 扁平带 depth（DFS，children 已按 sortIndex 排序；折叠的子树跳过但保留行）
const flat = computed(() => {
  const out = []
  const walk = (nodes, depth, hidden) => {
    nodes.forEach((n) => {
      if (!hidden) out.push({ node: n, depth })
      if (n.children && n.children.length) {
        walk(n.children, depth + 1, hidden || collapsedSet.value.has(n.nodeKey))
      }
    })
  }
  walk(store.tree, 0, false)
  return out
})

function toggle(key) {
  const s = new Set(collapsedSet.value)
  s.has(key) ? s.delete(key) : s.add(key)
  collapsedSet.value = s
}

const hasKids = (n) => !!(n.children && n.children.length)
// 任意节点都可再框子区域（叶子框入子节点后会自动转为容器）
const canAddChild = () => true

function displayName(n) {
  if (!n.name) return '未命名区域'
  return n.variantIndex > 0 ? `${n.name}-${n.variantIndex + 1}` : n.name
}
</script>

<style scoped>
.tree {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  color: rgba(255, 255, 255, 0.8);
}
.tree-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 6px 10px;
}
.tree-title { font-size: 13px; font-weight: 600; color: rgba(255, 255, 255, 0.6); }
.tree-count {
  font-size: 11px;
  min-width: 18px;
  height: 18px;
  line-height: 18px;
  text-align: center;
  padding: 0 5px;
  border-radius: 9px;
  background: rgba(0, 255, 255, 0.12);
  color: #00ffff;
}
.tree-body { flex: 1; overflow-y: auto; min-height: 0; }

.tree-row {
  position: relative;
  display: flex;
  align-items: center;
  gap: 7px;
  height: 34px;
  padding-right: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}
.tree-row:hover { background: rgba(0, 255, 255, 0.06); }
.tree-row--active {
  background: rgba(255, 0, 255, 0.08);
  box-shadow: inset 0 0 0 1px rgba(255, 0, 255, 0.35);
}
/* 层级竖线（思维导图观感的连线） */
.tree-guide {
  position: absolute;
  top: -6px;
  bottom: 50%;
  width: 1px;
  background: rgba(0, 255, 255, 0.18);
}

.twisty, .twisty-spacer {
  width: 14px;
  flex-shrink: 0;
  text-align: center;
  font-size: 11px;
  color: rgba(0, 255, 255, 0.6);
}
.twisty { cursor: pointer; }

.dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  flex-shrink: 0;
}
.node-name {
  flex: 1;
  min-width: 0;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.layout-chip {
  flex-shrink: 0;
  font-size: 12px;
  color: rgba(0, 255, 255, 0.7);
}

.row-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.15s;
}
.tree-row:hover .row-actions { opacity: 1; }
.mini {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.3);
  color: rgba(0, 255, 255, 0.8);
  font-size: 13px;
  line-height: 1;
  cursor: pointer;
}
.mini:hover { background: rgba(0, 255, 255, 0.12); color: #00ffff; }
.mini--danger { color: rgba(255, 255, 255, 0.5); border-color: rgba(255, 255, 255, 0.15); }
.mini--danger:hover { color: #ff4757; background: rgba(255, 71, 87, 0.12); border-color: rgba(255, 71, 87, 0.4); }

.tree-empty {
  padding: 24px 12px;
  text-align: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.25);
}
</style>
