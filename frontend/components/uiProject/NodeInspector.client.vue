<template>
  <div v-if="node" class="inspector">
    <div class="ins-head">
      <span class="ins-title">{{ displayName(node) }}</span>
      <div class="ins-actions">
        <span class="status-chip" :style="{ background: STATUS_COLOR[node.status] }">{{ statusText(node.status) }}</span>
      </div>
    </div>

    <div class="ins-body">
      <!-- 顶层：从枚举内联选类型（取代旧的图上浮层） -->
      <div v-if="isTop" class="field">
        <span class="field-label">区域类型<span v-if="!form.name" class="need-pick">（必选其一）</span></span>
        <div class="type-grid">
          <button
            v-for="t in TOP_REGION_ENUM"
            :key="t"
            class="type-chip"
            :class="{ active: form.name === t, used: otherUsed[t] }"
            @click="pickType(t)"
          >{{ t }}<span v-if="otherUsed[t]" class="chip-used">已用</span></button>
        </div>
      </div>

      <div v-if="isContainer" class="field field--row">
        <span class="field-label">排列方式
          <span class="help-tip" :data-tip="LAYOUT_TIP">
            <i class="fas fa-circle-question help-badge"></i>
          </span>
        </span>
        <div class="layout-seg">
          <button
            v-for="op in LAYOUT_OPTIONS"
            :key="op.value"
            class="seg-btn"
            :class="{ active: form.layout === op.value }"
            @click="setLayout(op.value)"
          >{{ op.label }}</button>
        </div>
      </div>

      <!-- 非顶层节点：自由命名（无论叶子还是容器都可命名；顶层走上方枚举类型）；
           已关联附图时名称随素材自动生成，禁止手动输入，取消关联后恢复可编辑 -->
      <label v-if="!isTop" class="field">
        <span class="field-label">
          区域名
          <span v-if="form.assetId" class="name-locked"><i class="fas fa-lock"></i> 随素材自动生成</span>
        </span>
        <input
          v-model="form.name"
          class="inp"
          maxlength="24"
          :disabled="!!form.assetId"
          :placeholder="form.assetId ? '已关联附图，名称随素材自动生成' : '子节点自由命名'"
        />
      </label>

      <label class="field">
        <span class="field-label">一句话描述</span>
        <textarea v-model="form.summary" class="inp txt" rows="2" placeholder="这个区域是什么、含哪些元素"></textarea>
      </label>

      <!-- 关联素材：引用当前主图的附图（Logo/图标/插画等真实素材）；编辑器内可直接补传，与外部托盘同一份数据 -->
      <div class="field">
        <div class="field-head">
          <span class="field-label">
            关联素材
            <span class="help-tip" :data-tip="ASSET_TIP"><i class="fas fa-circle-question help-badge"></i></span>
          </span>
          <button v-if="assetOptions.length" class="asset-add" :disabled="!canAddAsset" @click="pickAsset">
            <i class="fas fa-plus"></i> 上传
          </button>
        </div>
        <div v-if="assetOptions.length" class="asset-pick">
          <div
            v-for="a in assetOptions"
            :key="a.assetId"
            class="asset-opt"
            :class="{ active: form.assetId === a.assetId, 'asset-opt--err': !!a.uploadError }"
            :title="a.uploadError || a.label"
            @click="toggleAsset(a)"
          >
            <div class="asset-opt-thumb">
              <img v-if="a.preview" :src="a.preview" alt="" />
              <i v-else class="fas fa-image asset-opt-ph"></i>
              <!-- 上传中 / 失败蒙层：与外部素材托盘一致 -->
              <span v-if="a.uploading" class="asset-opt-mask"><i class="fas fa-spinner fa-spin"></i></span>
              <span v-else-if="a.uploadError" class="asset-opt-mask asset-opt-mask--err">
                <button class="mini-btn" @click.stop="codeStore.retryAsset(a.assetId)">重试</button>
              </span>
              <!-- hover 删除：与外部素材托盘一致；已选中的素材不显示删除，避免误触（点击为取消关联） -->
              <span v-if="form.assetId !== a.assetId" class="asset-opt-hover">
                <button class="asset-opt-del" title="删除素材" @click.stop="confirmRemoveAsset(a)">
                  <i class="fa-regular fa-trash-can"></i>
                </button>
              </span>
            </div>
            <span class="asset-opt-name">{{ a.label }}</span>
            <span v-if="a.ownerId !== store.assetScopeId" class="asset-opt-owner">其他页</span>
            <span v-if="form.assetId === a.assetId" class="asset-opt-check"><i class="fa-solid fa-check"></i></span>
          </div>
        </div>
        <!-- 无素材时也保留入口：标注过程中可直接补传，不必退出编辑器 -->
        <div v-else class="asset-empty" @click="pickAsset">
          <i class="fa-regular fa-images"></i>
          <span>暂无附图素材，点击上传 Logo / 图标 / 插画</span>
        </div>
        <span v-if="form.assetId" class="asset-link-tip">
          生成时该区域直接使用素材，不再由 AI 重绘
          <button class="asset-clear" @click="clearAsset">取消关联</button>
        </span>
        <input ref="assetInput" type="file" accept="image/*" multiple hidden @change="onAssetFiles" />
      </div>

      <label class="field">
        <span class="field-label">视觉要点</span>
        <textarea v-model="form.styleHints" class="inp txt" rows="2" placeholder="背景色 / 高度 / 分割线 / 字号 / 间距"></textarea>
      </label>

      <label v-if="isTop || hasChildren" class="field">
        <span class="field-label">交互行为</span>
        <input v-model="form.behavior" class="inp" placeholder="如 固定吸顶 / 下拉刷新，可空" />
      </label>
    </div>
  </div>
  <div v-else class="inspector inspector--empty">
    选中左侧框或右侧节点后可编辑描述
  </div>
</template>

<script setup>
import { reactive, computed, watch, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUiProjectStore } from '~/stores/uiProject'
import { useCodeStore } from '~/stores/code'
import { MAX_ASSETS } from '~/stores/code/upload'
import { TOP_REGION_ENUM } from '~/utils/uiTree.js'

const store = useUiProjectStore()
const codeStore = useCodeStore()

const STATUS_COLOR = {
  pending: 'rgba(255,255,255,0.3)',
  ai_draft: '#ffcf4d',
  confirmed: '#37e08d',
  empty: '#ff4757',
}
const STATUS_TEXT = { pending: '待描述', ai_draft: 'AI 初稿', confirmed: '已确认', empty: '待补全' }
const LAYOUT_OPTIONS = [
  { value: 'row', label: '横排 ⇄' },
  { value: 'column', label: '竖排 ⇅' },
  { value: 'grid', label: '网格 ⊞' },
]

// 排列方式的优先级说明：用 data-tip 交给纯 CSS 气泡渲染（不用 v-title，避免与原生 title 气泡叠层）
const LAYOUT_TIP = '排列方式取值优先级：手动 ＞ 自动推断 ＞ AI 建议。画完框选了类型时，若顶层 AI 给了排列方向会自动选中；AI 没给就留空不选。增删子区域后会按子框位置重新自动推断，手动选定过就不再被覆盖。'
const LAYOUT_SRC_TEXT = { manual: '手动', ai: 'AI 建议', auto: '自动推断' }
// 关联素材说明：（只在当前主图有可用附图时显示该字段）
const ASSET_TIP = '选择上传区附带的真实素材（Logo / 图标 / 插画 / 照片），生成时该区域直接使用这张素材，不再由 AI 重绘。'

const node = computed(() => store.selectedNode)
const isTop = computed(() => !!node.value && !node.value.parentId)
// 可引用/可管理的附图：展示全部素材，上传中与失败也入列表（loading/重试/删除与外部托盘一致）
const assetOptions = computed(() => codeStore.assets)
// 是否有子节点：非顶层但已下钻子级时，它充当父容器，同样需要交互行为字段
const hasChildren = computed(() => !!node.value && (node.value.childrenIds?.length || 0) > 0)
// 其它顶层节点已占用的枚举（用于“已用”标记，仍可点选自动生成序号）
const otherUsed = computed(() => {
  const set = {}
  if (!node.value) return set
  store.nodeList
    .filter((n) => !n.parentId && n.nodeKey !== node.value.nodeKey && n.name)
    .forEach((n) => { set[n.name] = (set[n.name] || 0) + 1 })
  return set
})
const isContainer = computed(() =>
  !!(node.value && (!node.value.parentId || node.value.childrenIds?.length || node.value.nodeType === 'container')),
)

const statusText = (s) => STATUS_TEXT[s] || s
function displayName(n) {
  if (!n.name) return '未命名区域'
  return n.variantIndex > 0 ? `${n.name}-${n.variantIndex + 1}` : n.name
}

// 表单本地态，改动 debounce 回写 store（触发 markEdited / manual layout）
const form = reactive({ name: '', summary: '', styleHints: '', behavior: '', layout: '', assetId: '' })
let timer = null

function pullFromNode(n) {
  form.name = n?.name || ''
  form.summary = n?.summary || ''
  form.styleHints = n?.styleHints || ''
  form.behavior = n?.behavior || ''
  form.layout = n?.layout || ''
  form.assetId = n?.assetId || ''
}

watch(node, (n) => { if (n) pullFromNode(n) }, { immediate: true })
// 外部改动（如顶部工具栏「AI 描述」批量回填）时，同步当前节点的描述到表单
watch(() => node.value?.summary, (s) => {
  if (node.value && form.summary !== (s || '')) form.summary = s || ''
})

function pushToStore() {
  if (!node.value) return
  store.updateNode(node.value.nodeKey, {
    name: form.name,
    summary: form.summary,
    styleHints: form.styleHints,
    behavior: form.behavior,
  })
}

// name/layout 单独处理：layout 走 manual，name 走 updateNode
watch(() => [form.name, form.summary, form.styleHints, form.behavior], () => {
  clearTimeout(timer)
  timer = setTimeout(pushToStore, 300)
})

function setLayout(val) {
  if (!node.value) return
  form.layout = val
  store.updateNode(node.value.nodeKey, { layout: val })
}

// 关联 / 取消关联附图素材：点击即回写（绑定视为已编辑确认）
// 素材归属其他主图时先改绑到当前主图，保证提交时它在 attached_assets 里
// 子级节点（有父级）关联素材时自动命名：附图 - 素材名；已手动命名的不覆盖
const isDefaultName = (name) => !name || name.startsWith('附图 - ')
const assetName = (a) => a.label || a.assetId || '素材'
function toggleAsset(a) {
  if (!node.value) return
  // 上传中 / 上传失败的素材不能关联（先等完成或重试）
  if (a.uploading || a.uploadError) return
  const next = form.assetId === a.assetId ? '' : a.assetId
  if (next && a.ownerId !== store.assetScopeId) {
    codeStore.setAssetOwner(a.assetId, store.assetScopeId)
    ElMessage.info(`素材「${a.label}」已归属到当前主图`)
  }
  form.assetId = next
  const patch = { assetId: next }
  if (next && node.value.parentId && isDefaultName(node.value.name)) {
    const auto = `附图 - ${assetName(a)}`
    patch.name = auto
    form.name = auto
  }
  // 取消勾选：自动生成的「附图 - 素材名」区域名一并清空，恢复未命名
  if (!next && isDefaultName(node.value.name)) {
    patch.name = ''
    form.name = ''
  }
  store.updateNode(node.value.nodeKey, patch)
}

function clearAsset() {
  if (!node.value) return
  form.assetId = ''
  const patch = { assetId: '' }
  // 解除关联时，自动生成的「附图 - 素材名」区域名同样清空
  if (isDefaultName(node.value.name)) {
    patch.name = ''
    form.name = ''
  }
  store.updateNode(node.value.nodeKey, patch)
}

// 删除素材：与外部托盘同一套确认流程；先解绑所有引用该素材的区域，再删除素材
async function confirmRemoveAsset(a) {
  try {
    await ElMessageBox.confirm(
      `确定删除素材「${a.label}」？所有引用该素材的区域将同时解绑。`,
      '删除附图',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' },
    )
  } catch {
    return
  }
  // 解绑所有引用该素材的区域（可能被多个节点引用），完整后再删除素材
  store.nodeList.forEach((n) => {
    if (n.assetId === a.assetId) store.updateNode(n.nodeKey, { assetId: '' })
  })
  if (form.assetId === a.assetId) form.assetId = ''
  codeStore.removeAsset(a.assetId)
}

// 编辑器内补传素材：与外部素材托盘同一份 store，上传完成自动出现在上方列表
const assetInput = ref(null)
const canAddAsset = computed(() => !codeStore.isSubmitting && codeStore.assets.length < MAX_ASSETS)
function pickAsset() {
  if (codeStore.isSubmitting) return
  if (codeStore.assets.length >= MAX_ASSETS) {
    ElMessage.warning(`附图最多 ${MAX_ASSETS} 张，请先删除不需要的素材`)
    return
  }
  assetInput.value?.click()
}
function onAssetFiles(e) {
  const files = Array.from(e.target.files || [])
  // 归属当前编辑主图（打开编辑器时已 setAssetScope，多图模式即当前页）
  if (files.length) codeStore.addAssets(files, store.assetScopeId)
  if (assetInput.value) assetInput.value.value = ''
}

// 顶层选类型：指派名称，刷新表单；位置异常软提示
function pickType(name) {
  if (!node.value) return
  const res = store.assignTopName(node.value.nodeKey, name)
  pullFromNode(node.value)
  if (res && res.ok && res.hints && res.hints.length) ElMessage.warning(`「${name}」${res.hints[0]}`)
}
</script>

<style scoped>
.inspector {
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  background: rgba(0, 255, 255, 0.04);
  border: 1px solid rgba(0, 255, 255, 0.14);
  overflow: hidden;
}
.inspector--empty {
  padding: 20px;
  text-align: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.28);
}
.ins-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid rgba(0, 255, 255, 0.1);
}
.ins-title { font-size: 13px; font-weight: 600; color: rgba(255, 255, 255, 0.85); }
.ins-actions { display: flex; align-items: center; gap: 8px; }
.status-chip {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  color: #0a0a0f;
  font-weight: 700;
}
.ins-body { padding: 12px; display: flex; flex-direction: column; gap: 12px; overflow-y: auto; container-type: inline-size; }
.field { display: flex; flex-direction: column; gap: 5px; }
.field--row { flex-direction: row; flex-wrap: wrap; align-items: center; gap: 8px; }
.field-label { font-size: 12px; color: rgba(255, 255, 255, 0.45); }
.inp {
  width: 100%;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(0, 255, 255, 0.22);
  border-radius: 8px;
  padding: 7px 9px;
  color: #fff;
  font-size: 13px;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s;
}
.inp:focus { border-color: #00ffff; box-shadow: 0 0 8px rgba(0, 255, 255, 0.1); }
/* 已关联附图：区域名由素材自动生成，输入框禁用 */
.inp:disabled { opacity: 0.55; cursor: not-allowed; border-color: rgba(255, 255, 255, 0.12); }
.name-locked { margin-left: 6px; font-size: 10px; color: rgba(0, 255, 255, 0.6); display: inline-flex; align-items: center; gap: 3px; }
.txt { resize: vertical; }
.field--row .inp { flex: none; width: auto; }
.layout-seg { display: flex; gap: 6px; }
.seg-btn {
  padding: 4px 10px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.3);
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  cursor: pointer;
}
.seg-btn.active { border-color: #00ffff; color: #00ffff; background: rgba(0, 255, 255, 0.1); }
.src-chip {
  margin-left: 6px;
  padding: 1px 6px;
  border-radius: 8px;
  font-size: 10px;
  color: rgba(0, 255, 255, 0.75);
  background: rgba(0, 255, 255, 0.08);
  border: 1px solid rgba(0, 255, 255, 0.22);
}
/* 字段说明气泡：圆形徽章 + 自绘 tooltip（跟随文字换行，不依赖指令/原生 title） */
/* cursor 不用 help：Windows 的 help 光标自带一个小问号，会与徽章上的问号重复 */
.help-tip { position: relative; display: inline-flex; vertical-align: middle; margin-left: 5px; cursor: pointer; }
.help-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  font-size: 10px;
  color: #00ffff;
  background: rgba(0, 255, 255, 0.12);
  border: 1px solid rgba(0, 255, 255, 0.45);
  transition: all 0.2s;
}
.help-tip:hover .help-badge { background: rgba(0, 255, 255, 0.25); box-shadow: 0 0 8px rgba(0, 255, 255, 0.35); }
.help-tip::after {
  content: attr(data-tip);
  position: absolute;
  left: 0;
  bottom: calc(100% + 8px);
  /* 限定在滚动面板内（.ins-body 为 inline-size 容器），避免比面板宽被 overflow 裁掉 */
  width: min(240px, calc(100cqw - 12px));
  padding: 8px 10px;
  border-radius: 8px;
  background: rgba(12, 14, 22, 0.97);
  border: 1px solid rgba(0, 255, 255, 0.32);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.5);
  color: rgba(255, 255, 255, 0.82);
  font-size: 12px;
  line-height: 1.6;
  white-space: normal;
  opacity: 0;
  visibility: hidden;
  transform: translateY(4px);
  transition: opacity 0.15s, transform 0.15s, visibility 0s linear 0.15s;
  z-index: 30;
  pointer-events: none;
}
.help-tip:hover::after { opacity: 1; visibility: visible; transform: translateY(0); transition-delay: 0s; }

/* 顶层类型内联胶囊 */
.need-pick { color: #ffcf4d; margin-left: 4px; }
.type-grid { display: flex; flex-wrap: wrap; gap: 6px; }
.type-chip {
  position: relative;
  padding: 5px 10px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.3);
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.type-chip:hover { border-color: rgba(0, 255, 255, 0.5); color: #00ffff; }
.type-chip.active { border-color: #00ffff; color: #00ffff; background: rgba(0, 255, 255, 0.12); }
.type-chip.used { opacity: 0.62; }
.chip-used { margin-left: 4px; font-size: 10px; color: rgba(255, 207, 77, 0.85); }

/* 关联素材：缩略图选择 */
.field-head { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.asset-add {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 6px;
  border: 1px solid rgba(0, 255, 255, 0.3);
  background: transparent;
  color: #00ffff;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}
.asset-add:hover:not(:disabled) { background: rgba(0, 255, 255, 0.1); }
.asset-add:disabled { opacity: 0.4; cursor: not-allowed; }
.asset-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  min-height: 56px;
  border: 1px dashed rgba(0, 255, 255, 0.22);
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.35);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}
.asset-empty:hover { border-color: rgba(0, 255, 255, 0.5); color: rgba(0, 255, 255, 0.75); }
.asset-empty i { font-size: 14px; }
.asset-pick { display: flex; flex-wrap: wrap; gap: 8px; }
.asset-opt {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  width: 64px;
  padding: 4px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(0, 0, 0, 0.3);
  cursor: pointer;
  transition: all 0.2s;
}
/* 缩略图容器：loading / 错误 / hover 删除蒙层的定位参照 */
.asset-opt-thumb {
  position: relative;
  width: 54px;
  height: 54px;
  border-radius: 6px;
  overflow: hidden;
  background: #0a0a0f;
}
.asset-opt img { width: 100%; height: 100%; object-fit: contain; display: block; }
.asset-opt-ph { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: rgba(255, 255, 255, 0.15); }
/* 上传中 / 失败蒙层（与外部素材托盘一致） */
.asset-opt-mask {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
}
.asset-opt-mask--err { color: #ff4757; }
.mini-btn {
  padding: 3px 10px;
  border-radius: 6px;
  border: 1px solid #00ffff;
  background: rgba(0, 0, 0, 0.5);
  color: #00ffff;
  font-size: 11px;
  cursor: pointer;
}
/* hover 蒙层 + 中央删除：与外部素材托盘一致 */
.asset-opt-hover {
  position: absolute;
  inset: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(5, 8, 12, 0.55);
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.18s ease;
}
.asset-opt:hover .asset-opt-hover { opacity: 1; }
.asset-opt-del {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  border: 1px solid rgba(255, 71, 87, 0.55);
  background: rgba(20, 10, 14, 0.75);
  color: #ff6b7a;
  font-size: 13px;
  cursor: pointer;
  pointer-events: auto;
  transition: all 0.18s ease;
}
.asset-opt-del:hover { border-color: #ff4757; background: rgba(255, 71, 87, 0.22); color: #fff; }
.asset-opt-name {
  max-width: 100%;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.asset-opt:hover { border-color: rgba(0, 255, 255, 0.45); }
.asset-opt.active { border-color: #00ffff; box-shadow: 0 0 10px rgba(0, 255, 255, 0.22); }
.asset-opt.active .asset-opt-name { color: #00ffff; }
.asset-opt--err { border-color: rgba(255, 71, 87, 0.45); }
.asset-opt-owner {
  position: absolute;
  top: 2px;
  left: 2px;
  z-index: 3;
  padding: 0 4px;
  border-radius: 4px;
  background: rgba(255, 207, 77, 0.15);
  border: 1px solid rgba(255, 207, 77, 0.4);
  color: rgba(255, 207, 77, 0.9);
  font-size: 9px;
  line-height: 1.4;
}
.asset-opt-check {
  position: absolute;
  top: 2px;
  right: 2px;
  z-index: 3;
  width: 15px;
  height: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: #00ffff;
  color: #04110f;
  font-size: 8px;
}
.asset-link-tip { font-size: 11px; color: rgba(0, 255, 255, 0.55); }
.asset-clear {
  margin-left: 6px;
  padding: 0;
  border: none;
  background: none;
  color: #ff6b7a;
  font-size: 11px;
  cursor: pointer;
}
.asset-clear:hover { text-decoration: underline; }
</style>
