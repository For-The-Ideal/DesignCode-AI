<template>
  <div v-if="node" class="inspector">

    <!-- 焦点离开输入区即落库：input 的 300ms debounce 未触发时点「开始生成」也不会漏掉最新区域名 -->
    <div class="ins-body" @focusout="flushNow">
      <!-- 顶层：从枚举内联选类型（取代旧的图上浮层） -->
      <div v-if="isTop" class="field">
        <span class="field-label"><span class="req">*</span>区域类型<span v-if="!form.name" class="need-pick">（必选其一）</span></span>
        <div class="type-grid">
          <button
            v-for="t in TOP_REGION_ENUM"
            :key="t"
            class="type-chip"
            :class="{ active: form.name === t, used: otherUsed[t] }"
            :title="`${t}｜${TOP_REGION_TIPS[t] || ''}`"
            @click="pickType(t)"
          >{{ regionLabel(t) }}<span v-if="otherUsed[t]" class="chip-used">已用</span></button>
        </div>
      </div>

      <!-- ≥2 级：先选元素类型（容器/叶子同表），区域名只写业务含义；
           已关联附图时名称随素材自动生成，禁止手动输入，取消关联后恢复可编辑 -->
      <div v-if="!isTop" class="field">
        <span class="field-label"><span class="req">*</span>元素类型<span v-if="!form.component" class="need-pick">（必选其一）</span></span>
        <div v-for="group in ELEMENT_TYPE_GROUPS" :key="group.label" class="type-row">
          <span class="type-row-label">{{ group.label }}</span>
          <div class="type-grid type-grid--dense">
            <button
              v-for="t in group.types"
              :key="t"
              class="type-chip type-chip--sm"
              :class="{ active: form.component === t }"
              :title="`${t}｜${ELEMENT_TYPE_TIPS[t] || ''}${form.component === t ? '（再次点击取消）' : ''}`"
              @click="pickElementType(t)"
            >{{ elementLabel(t) }}</button>
          </div>
        </div>
      </div>

      <div v-if="isContainer" class="field field--row">
        <span class="field-label">排列方式
          <span v-if="form.layout" class="layout-src">{{ LAYOUT_SRC_TEXT[form.layoutSource] || '已设置' }}</span>
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
        <div v-if="childCount" class="layout-sub">决定这块区域内 {{ childCount }} 个子区域怎么排；子区域自己内部再各看各的排列方式</div>
      </div>

      <label v-if="!isTop" class="field">
        <span class="field-label">
          <span class="req">*</span>区域名
          <span class="help-tip" :data-tip="NAME_TIP"><i class="fas fa-circle-question help-badge"></i></span>
          <span v-if="form.assetId" class="name-locked"><i class="fas fa-lock"></i> 随素材自动生成</span>
        </span>
        <input
          v-model="form.name"
          class="inp"
          maxlength="24"
          :disabled="!!form.assetId"
          :placeholder="form.assetId ? '已关联附图，名称随素材自动生成' : '写这块框的业务含义'"
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
          <span>点击上传素材 Logo / 图标 / 插画</span>
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
    <i class="fas fa-hand-pointer"></i>
    <span>选中区域后可编辑描述</span>
  </div>
</template>

<script setup>
import { reactive, computed, watch, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUiProjectStore } from '~/stores/uiProject'
import { useCodeStore } from '~/stores/code'
import { MAX_ASSETS } from '~/stores/code/upload'
import {
  TOP_REGION_ENUM,
  TOP_REGION_TIPS,
  regionLabel,
  ELEMENT_CONTAINER_TYPES,
  ELEMENT_LEAF_TYPES,
  ELEMENT_TYPE_TIPS,
  elementLabel,
  isContainerElement,
} from '~/utils/uiTree.js'

const store = useUiProjectStore()
const codeStore = useCodeStore()

const STATUS_COLOR = {
  pending: 'rgba(255,255,255,0.3)',
  ai_draft: '#ffcf4d',
  confirmed: '#37e08d',
  empty: '#ff4757',
}
const STATUS_TEXT = { pending: '待描述', ai_draft: 'AI 初稿', confirmed: '已确认', empty: '待补全' }
// 排列方式：row/column/grid 为流式排列；stack 为叠放（子元素共用同一块区域、按坐标叠加，如背景图 + 浮层）
const LAYOUT_OPTIONS = [
  { value: 'row', label: '横排 ⇄' },
  { value: 'column', label: '竖排 ⇅' },
  { value: 'grid', label: '网格 ⊞' },
  { value: 'stack', label: '叠放 ⧉' },
]
// 元素类型分两行：容器型（可继续下钻）/ 叶子型
const ELEMENT_TYPE_GROUPS = [
  { label: '容器', types: ELEMENT_CONTAINER_TYPES },
  { label: '元素', types: ELEMENT_LEAF_TYPES },
]

// 排列方式的优先级说明：用 data-tip 交给纯 CSS 气泡渲染（不用 v-title，避免与原生 title 气泡叠层）
const LAYOUT_TIP = '排列方式取值优先级：手动 ＞ 自动推断 ＞ AI 建议。画完框选了类型时，若顶层 AI 给了排列方向会自动选中；AI 没给就留空不选。增删子区域后会按子框位置重新自动推断，手动选定过就不再被覆盖。叠放＝子元素不排队，共用这块区域按原图坐标叠加（典型：底图铺满 + 文案/按钮/卡片浮在上面）。再次点击已选中的排列方式可取消。'
const LAYOUT_SRC_TEXT = { manual: '手动', ai: 'AI 建议', auto: '自动推断' }
// 关联素材说明：（只在当前主图有可用附图时显示该字段）
const ASSET_TIP = '选择上传区附带的真实素材（Logo / 图标 / 插画 / 照片），生成时该区域直接使用这张素材，不再由 AI 重绘。'
// 区域名说明：名称只承载业务语义（元素类型已在上方单独选），生成前必填
const NAME_TIP = '只写业务含义（这块框承载的文案或功能），不要重复上方的元素类型；生成前必填，会显示在层级树与框线标签上。'

const node = computed(() => store.selectedNode)
const isTop = computed(() => !!node.value && !node.value.parentId)
// 可引用/可管理的附图：只列当前主图名下的素材（多图切页后自动跟着变），
// 上传中与失败也入列表（loading/重试/删除与外部托盘一致）
const assetOptions = computed(() => {
  const all = codeStore.assets
  const scope = store.assetScopeId
  return scope ? all.filter((a) => !a.ownerId || a.ownerId === scope) : all
})
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
  !!(
    node.value &&
    (!node.value.parentId ||
      node.value.childrenIds?.length ||
      node.value.nodeType === 'container' ||
      isContainerElement(node.value.component))
  ),
)


// 表单本地态，改动 debounce 回写 store（触发 markEdited / manual layout）
const form = reactive({ name: '', component: '', summary: '', styleHints: '', behavior: '', layout: '', layoutSource: '', assetId: '' })
// 该区域内直接子区域数量：排列方式描述的是它们之间的排布
const childCount = computed(() => node.value?.childrenIds?.length || 0)
let timer = null

function pullFromNode(n) {
  form.name = n?.name || ''
  form.component = n?.component || ''
  form.summary = n?.summary || ''
  form.styleHints = n?.styleHints || ''
  form.behavior = n?.behavior || ''
  form.layout = n?.layout || ''
  form.layoutSource = n?.layoutSource || ''
  form.assetId = n?.assetId || ''
}

// 切换节点前先落库上一个节点的输入，再载入新节点内容
watch(node, (n, old) => {
  if (old) flushTo(old)
  if (n) pullFromNode(n)
}, { immediate: true })
// 外部改动（如顶部工具栏「AI 描述」批量回填）时，同步当前节点的描述到表单
watch(() => node.value?.summary, (s) => {
  if (node.value && form.summary !== (s || '')) form.summary = s || ''
})

// 立即把表单写回指定节点（不依赖当前选中）：失焦、切换节点前调用，
// 避免 debounce 未落库时生成校验读到旧的区域名
function flushTo(target) {
  if (timer) { clearTimeout(timer); timer = null }
  if (!target) return
  store.updateNode(target.nodeKey, {
    name: form.name,
    summary: form.summary,
    styleHints: form.styleHints,
    behavior: form.behavior,
  })
}
const flushNow = () => flushTo(node.value)

// name/layout 单独处理：layout 走 manual，name 走 updateNode
watch(() => [form.name, form.summary, form.styleHints, form.behavior], () => {
  clearTimeout(timer)
  timer = setTimeout(() => flushTo(node.value), 300)
})

// 点击选中；再次点击同一个即取消（显式清空，交给 layoutSource=manual 阻止自动推断覆盖）
function setLayout(val) {
  if (!node.value) return
  const next = form.layout === val ? '' : val
  form.layout = next
  form.layoutSource = 'manual'
  store.updateNode(node.value.nodeKey, { layout: next })
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
  if (res && res.ok && res.hints && res.hints.length) ElMessage.warning(`「${regionLabel(name)}」${res.hints[0]}`)
}

// 非顶层选元素类型：再次点击已选中的类型 = 取消选择（回到未选，生成侧按容器推断）
function pickElementType(type) {
  if (!node.value) return
  const next = node.value.component === type ? '' : type
  const res = store.assignElementType(node.value.nodeKey, next)
  pullFromNode(node.value)
  if (res && res.hints && res.hints.length) ElMessage.warning(res.hints[0])
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
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 24px 20px;
  text-align: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.3);
}
.inspector--empty i { font-size: 18px; color: rgba(0, 255, 255, 0.35); }
.ins-body { padding: 12px; display: flex; flex-direction: column; gap: 12px; overflow-x: hidden; overflow-y: auto; container-type: inline-size; }
.field { position: relative; display: flex; flex-direction: column; gap: 5px; }
.field--row { flex-direction: row; flex-wrap: wrap; align-items: center; gap: 8px; }
.field-label { font-size: 12px; color: rgba(255, 255, 255, 0.45); }
/* 必填项标记：红色星号统一放在标签前 */
.req { color: #ff4d5e; margin-right: 3px; }
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
/* 排列方式来源：让用户知道已有值且是系统按子框推断的，不必手填 */
.layout-src { margin-left: 5px; font-size: 10px; color: rgba(0, 255, 255, 0.55); }
.layout-sub { flex-basis: 100%; font-size: 10px; line-height: 1.5; color: rgba(255, 255, 255, 0.32); }
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
/* 气泡相对整个字段定位（.field 为定位祖先），而不是相对问号小图标：
   否则图标靠左 + 气泡 240px 宽会把面板撑出横向溢出，触发横向滚动条 */
.help-tip { position: static; display: inline-flex; vertical-align: middle; margin-left: 5px; cursor: pointer; }
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

/* 元素类型（≥2 级）：容器/元素两行，胶囊更紧凑 */
.type-row { display: flex; gap: 6px; }
.type-row + .type-row { margin-top: 6px; }
.type-row-label { flex: none; width: 24px; padding-top: 4px; font-size: 10px; color: rgba(255, 255, 255, 0.35); }
.type-grid--dense { flex: 1; gap: 4px; }
.type-chip--sm { padding: 3px 7px; border-radius: 6px; font-size: 11px; }

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
/* 删除：与选中互斥（选中时模板不渲染），未选中时仅 hover 在右上角浮现 */
.asset-opt-hover {
  position: absolute;
  top: 2px;
  right: 2px;
  z-index: 2;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.18s ease;
}
.asset-opt:hover .asset-opt-hover { opacity: 1; }
.asset-opt-del {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  border: 1px solid rgba(255, 71, 87, 0.55);
  background: rgba(20, 10, 14, 0.8);
  color: #ff6b7a;
  font-size: 10px;
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
