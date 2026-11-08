<template>
  <!-- 页面关系段（多图，P2 起显示）：以当前页为主体，维护「本页由哪一页进入」 -->
  <div class="rel-sec">
    <div class="rel-list">
        <div
          v-for="l in rows"
          :key="l.id"
          class="rel-row"
          :class="{ active: l.id === editId }"
          @click="startEdit(l)"
        >
          <span class="rel-kind" :class="`rel-kind--${l.kind || 'route'}`">{{ kindLabel(l.kind) }}</span>
          <span class="rel-text">
            <b>{{ pageNameOf(l.fromPageId) }}</b>
            <template v-if="l.label"> · {{ l.label }}</template>
            <i class="fas fa-arrow-right-long rel-arrow"></i>
            <b>{{ currentLabel }}</b>
          </span>
          <button class="rel-del" title="删除关系" @click.stop="onRemove(l)">
            <i class="fa-regular fa-trash-can"></i>
          </button>
        </div>

        <div v-if="!rows.length" class="rel-empty">
          <i class="fas fa-diagram-project"></i>
          <span>本页还没有上层页面，点「新增关系」指定它与上页的关联关系</span>
        </div>
      </div>

      <button v-if="!editing" class="rel-add" :disabled="pages.length < 2" @click="startCreate()">
        <i class="fas fa-plus"></i> 新增关系
      </button>

      <div v-if="editing" class="rel-form">
        <label class="form-fld">
          <span>上层页面 <em>（从哪一页进入本页）</em></span>
          <select v-model="form.fromPageId" class="form-inp" @change="onFromChange">
            <option v-for="p in parentPages" :key="p.pageId" :value="p.pageId">{{ p.label }}</option>
          </select>
        </label>

        <label class="form-fld">
          <span>源组件 <em>（空 = 整页）</em></span>
          <select v-model="form.nodeKey" class="form-inp">
            <option value="">整页</option>
            <option v-for="n in fromNodes" :key="n.nodeKey" :value="n.nodeKey">{{ n.label }}</option>
          </select>
        </label>

        <div class="form-row">
          <label class="form-fld">
            <span>关系类型</span>
            <select v-model="form.kind" class="form-inp">
              <option v-for="k in KINDS" :key="k.value" :value="k.value">{{ k.label }}</option>
            </select>
          </label>
          <label class="form-fld">
            <span>说明</span>
            <input v-model="form.desc" class="form-inp" maxlength="40" placeholder="如 点击活动入口进入" />
          </label>
        </div>

        <p class="form-target">进入本页：<b>{{ currentLabel }}</b></p>

        <div class="form-actions">
          <button class="btn-cancel" @click="cancel">取消</button>
          <button class="btn-save" :disabled="!canSave" @click="save">{{ editId ? '保存' : '添加' }}</button>
        </div>
      </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useCodeStore } from '~/stores/code'
import { useUiProjectStore } from '~/stores/uiProject'
import { nodeLabel } from '~/utils/uiTree.js'

const store = useCodeStore()
const uiStore = useUiProjectStore()

// 关系类型：kind 是前端语义，提交时由 store 派生合法 trigger
const KINDS = [
  { value: 'route', label: '路由跳转' },
  { value: 'modal', label: '弹层' },
  { value: 'tab', label: '底部 Tab' },
  { value: 'back', label: '返回' },
  { value: 'auto', label: '自动' },
]
const kindLabel = (kind) => KINDS.find((k) => k.value === (kind || 'route'))?.label || '路由跳转'

// ═══ 当前页与页序 ═══
const pages = computed(() =>
  store.images.map((img, i) => ({
    pageId: img.pageId,
    label: `P${i + 1} ${img.pageName || (img.file?.name || '').replace(/\.[^.]+$/, '') || '未命名'}`,
  })),
)
const pageNameOf = (id) => pages.value.find((p) => p.pageId === id)?.label || id

const curImg = computed(() => store.currentImage)
const curIdx = computed(() => store.images.findIndex((i) => i.pageId === curImg.value?.pageId))
// 列表第 1 页 = 入口页：没有上层，也就没有页面关系
const isEntry = computed(() => curIdx.value <= 0)
const currentLabel = computed(() => (curImg.value ? pageNameOf(curImg.value.pageId) : '本页'))

// 本页的关系 = 指向本页的边（每页只声明自己的上层）
const rows = computed(() =>
  store.pageLinks.filter((l) => l.toPageId === curImg.value?.pageId),
)
// 上层候选：除本页外的所有页
const parentPages = computed(() => pages.value.filter((p) => p.pageId !== curImg.value?.pageId))

// ═══ 表单 ═══
const editing = ref(false)
const editId = ref('')
const form = reactive({ fromPageId: '', nodeKey: '', kind: 'route', desc: '' })

const canSave = computed(() => !!(form.fromPageId && form.fromPageId !== curImg.value?.pageId))

function flatTree(list, out = []) {
  (list || []).forEach((n) => { out.push(n); flatTree(n.children, out) })
  return out
}

// 源组件取「上层页面」的标注节点：当前页读 uiStore（含未写回的编辑），其它页读该页已确认的树
const fromNodes = computed(() => {
  const img = store.images.find((i) => i.pageId === form.fromPageId)
  if (!img) return []
  const list = form.fromPageId === store.currentImage?.pageId ? uiStore.nodeList : flatTree(img.annotationTree)
  return list.map((n) => ({
    nodeKey: n.nodeKey,
    label: nodeLabel(n),
  }))
})

/* 入口页不应进入本段（关系段在 P1 不渲染），做出防御避免误建边 */
function startCreate() {
  if (isEntry.value || store.images.length < 2) {
    ElMessage.warning('至少需要 2 个页面才能建立页面关系')
    return
  }
  const prev = store.images[curIdx.value - 1]
  form.fromPageId = prev?.pageId || parentPages.value[0]?.pageId || ''
  // 源组件属于上层页，默认整页；需要时再从上层页的节点里选
  form.nodeKey = ''
  form.kind = 'route'
  form.desc = ''
  editId.value = ''
  editing.value = true
}

function startEdit(link) {
  editId.value = link.id
  form.fromPageId = link.fromPageId
  form.nodeKey = link.nodeKey || ''
  form.kind = link.kind || 'route'
  form.desc = link.desc || ''
  editing.value = true
}

function onFromChange() {
  // 换上层页后原节点不属于新页，清空避免悬空引用
  if (!fromNodes.value.some((n) => n.nodeKey === form.nodeKey)) form.nodeKey = ''
}

function save() {
  if (!canSave.value) return
  const toPageId = curImg.value?.pageId
  const dup = store.pageLinks.some(
    (l) => l.id !== editId.value && l.fromPageId === form.fromPageId && l.toPageId === toPageId,
  )
  if (dup) {
    ElMessage.warning('这条页面关系已存在')
    return
  }
  const label = form.nodeKey ? (fromNodes.value.find((n) => n.nodeKey === form.nodeKey)?.label || '') : ''
  const extra = { kind: form.kind, nodeKey: form.nodeKey, label, desc: form.desc.trim() }
  if (editId.value) {
    store.updatePageLink(editId.value, { fromPageId: form.fromPageId, toPageId, ...extra })
    ElMessage.success('页面关系已更新')
  } else {
    store.addPageLink(form.fromPageId, toPageId, undefined, extra)
    ElMessage.success('页面关系已添加')
  }
  cancel()
}

function cancel() {
  editing.value = false
  editId.value = ''
}

// 切页时关闭表单：关系固定归属当前页，避免把在编辑的关系带到别的页
watch(() => curImg.value?.pageId, () => { if (editing.value) cancel() })

function onRemove(link) {
  store.removePageLink(link.id)
  if (editId.value === link.id) cancel()
}
</script>

<style scoped>
.rel-sec { display: flex; flex-direction: column; gap: 8px; }

.rel-list { display: flex; flex-direction: column; gap: 6px; }
.rel-row {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 6px 8px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.07);
  background: rgba(255, 255, 255, 0.02);
  cursor: pointer;
  transition: all 0.15s;
}
.rel-row:hover { border-color: rgba(0, 255, 255, 0.28); }
.rel-row.active { border-color: rgba(0, 255, 255, 0.45); background: rgba(0, 255, 255, 0.06); }

.rel-kind {
  flex-shrink: 0;
  padding: 1px 6px;
  border-radius: 5px;
  font-size: 10px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
  background: rgba(255, 255, 255, 0.08);
}
.rel-kind--route { color: #00ffff; background: rgba(0, 255, 255, 0.12); }
.rel-kind--modal { color: #c084fc; background: rgba(192, 132, 252, 0.14); }
.rel-kind--tab { color: #37e08d; background: rgba(55, 224, 141, 0.12); }
.rel-kind--back { color: #fbbf24; background: rgba(251, 191, 36, 0.12); }

.rel-text {
  flex: 1;
  min-width: 0;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.55);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.rel-text b { font-weight: 600; color: rgba(255, 255, 255, 0.85); }
.rel-arrow { margin: 0 4px; font-size: 10px; color: rgba(0, 255, 255, 0.55); }

.rel-del {
  flex-shrink: 0;
  width: 22px;
  height: 22px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.3);
  font-size: 11px;
  cursor: pointer;
}
.rel-del:hover { color: #ff6b7a; background: rgba(255, 71, 87, 0.12); }

.rel-empty {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px;
  border-radius: 8px;
  border: 1px dashed rgba(255, 255, 255, 0.1);
  font-size: 11px;
  color: rgba(255, 255, 255, 0.3);
}
.rel-empty i { color: rgba(0, 255, 255, 0.35); }

.rel-add {
  width: 100%;
  padding: 6px 0;
  border-radius: 8px;
  border: 1px dashed rgba(0, 255, 255, 0.3);
  background: rgba(0, 255, 255, 0.04);
  color: rgba(0, 255, 255, 0.75);
  font-size: 11px;
  cursor: pointer;
}
.rel-add:hover:not(:disabled) { background: rgba(0, 255, 255, 0.1); }
.rel-add:disabled { opacity: 0.4; cursor: not-allowed; }

.rel-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 9px;
  border-radius: 9px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.35);
}
.form-row { display: flex; gap: 8px; }
.form-row .form-fld { flex: 1; min-width: 0; }
.form-fld { display: flex; flex-direction: column; gap: 4px; }
.form-fld > span { font-size: 11px; color: rgba(255, 255, 255, 0.5); }
.form-fld > span em { font-style: normal; color: rgba(255, 255, 255, 0.28); }
.form-inp {
  width: 100%;
  padding: 5px 7px;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(0, 0, 0, 0.45);
  color: rgba(255, 255, 255, 0.85);
  font-size: 11px;
  outline: none;
}
.form-inp:focus { border-color: rgba(0, 255, 255, 0.45); }
.form-inp option { background: #0b0f17; }

/* 目标页固定为当前页，只做只读展示 */
.form-target {
  margin: 0;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.38);
}
.form-target b { color: rgba(0, 255, 255, 0.85); font-weight: 600; }

.form-actions { display: flex; justify-content: flex-end; gap: 8px; }
.btn-cancel, .btn-save {
  padding: 5px 12px;
  border-radius: 7px;
  font-size: 11px;
  cursor: pointer;
}
.btn-cancel {
  border: 1px solid rgba(255, 255, 255, 0.14);
  background: transparent;
  color: rgba(255, 255, 255, 0.55);
}
.btn-save {
  border: 1px solid rgba(0, 255, 255, 0.45);
  background: rgba(0, 255, 255, 0.14);
  color: #00ffff;
}
.btn-save:disabled { opacity: 0.4; cursor: not-allowed; }
</style>
