<template>
  <!-- 右栏：一列可折叠分区（区域 = 树 + 属性 / 生成配置），不切 tab；底部常驻开始生成 -->
  <div class="wb-inspector">
    <div class="insp-scroll">
      <!-- 区域：树在上、选中节点的属性紧接在下，一段内完成「选节点 → 改属性」 -->
      <section class="insp-sec">
        <button class="sec-head" @click="toggle('nodes')">
          <i class="fas" :class="open.nodes ? 'fa-chevron-down' : 'fa-chevron-right'"></i>
          <i class="fas fa-sitemap sec-icon"></i>
          <span class="sec-title">区域结构</span>
          <span class="sec-badge">{{ store.nodeList.length }}</span>
        </button>
        <div v-show="open.nodes">
          <div class="sec-body sec-body--tree">
            <HierarchyTree
              :hide-head="true"
              @addChild="(key) => emit('addChild', key)"
              @deleteNode="(key) => emit('deleteNode', key)"
              @select="(key) => emit('select', key)"
            />
          </div>
          <!-- 没有区域时提示已在树区统一展示，属性区不再重复空态 -->
          <div v-if="store.hasNodes" class="sec-body sec-body--form">
            <NodeInspector />
          </div>
        </div>
      </section>

      <!-- 页面关系（多图）：以当前页为主体；P1 是入口页没有上层，从 P2 起才显示本段 -->
      <section v-if="isMulti && !isEntryPage" class="insp-sec">
        <button class="sec-head" @click="toggle('relations')">
          <i class="fas" :class="open.relations ? 'fa-chevron-down' : 'fa-chevron-right'"></i>
          <i class="fas fa-diagram-project sec-icon"></i>
          <span class="sec-title">页面关系</span>
          <span class="sec-badge">{{ relCount }}</span>
        </button>
        <div v-show="open.relations" class="sec-body sec-body--rel">
          <PageRelationSection />
        </div>
      </section>

      <!-- 生成配置 -->
      <section class="insp-sec">
        <button class="sec-head" @click="toggle('config')">
          <i class="fas" :class="open.config ? 'fa-chevron-down' : 'fa-chevron-right'"></i>
          <i class="fas fa-sliders sec-icon"></i>
          <span class="sec-title">生成配置</span>
          <span v-if="configSummary" class="sec-badge">{{ configSummary }}</span>
        </button>
        <div v-show="open.config" class="sec-body sec-body--config">
          <!-- 配置是任务级全局项：多页共一份，避免各页各选一套导致冲突 -->
          <p v-if="isMulti" class="cfg-shared">
            <i class="fas fa-circle-info"></i>
            以下配置对整个任务页面统一生效，不随当前页切换
          </p>
          <ConfigPanel :show-generate="false" />
        </div>
      </section>
    </div>

    <!-- 底部常驻：折叠/滚动都不影响点击 -->
    <div class="insp-foot">
      <GenerateButton />
    </div>
  </div>
</template>

<script setup>
import { reactive, computed } from 'vue'
import HierarchyTree from '~/components/uiProject/HierarchyTree.client.vue'
import NodeInspector from '~/components/uiProject/NodeInspector.client.vue'
import ConfigPanel from '~/components/code/ConfigPanel.vue'
import GenerateButton from '~/components/code/GenerateButton.vue'
import PageRelationSection from './PageRelationSection.vue'
import { useUiProjectStore } from '~/stores/uiProject'
import { useCodeStore } from '~/stores/code'
import { useConfigStore } from '~/stores/config'

const emit = defineEmits(['addChild', 'deleteNode', 'select'])

const store = useUiProjectStore()
const codeStore = useCodeStore()
const configStore = useConfigStore()
const isMulti = computed(() => codeStore.mode === 'multi')

// 三段折叠：区域（树 + 属性）/ 页面关系（多图）/ 生成配置；默认全展开
const open = reactive({ nodes: true, relations: true, config: true })
const toggle = (key) => { open[key] = !open[key] }

// 配置段徽标：目标框架 · 目标平台（配置字典未加载完成时为空）
const configSummary = computed(() => {
  const fw = (configStore.frameworks || []).find((f) => f.value === codeStore.config.framework)?.label
  const pl = (configStore.platforms || []).find((p) => p.value === codeStore.config.platform)?.label
  return [fw, pl].filter(Boolean).join(' · ')
})

// 当前页是否为入口页：列表第 1 页（无图时也按入口处理，不显示关系段）
const isEntryPage = computed(() => {
  const idx = codeStore.images.findIndex((i) => i.pageId === codeStore.currentImage?.pageId)
  return idx <= 0
})

// 关系段徽标：当前页的上层关系条数（入口页恒为 0）
const relCount = computed(() => {
  const cur = codeStore.currentImage?.pageId
  return cur ? codeStore.pageLinks.filter((l) => l.toPageId === cur).length : 0
})
</script>

<style scoped>
.wb-inspector {
  display: flex;
  flex-direction: column;
  min-height: 0;
  height: 100%;
  border-radius: 14px;
  border: 1px solid rgba(0, 255, 255, 0.08);
  background: rgba(10, 14, 23, 0.7);
  overflow: hidden;
}

/* 唯一滚动容器：三段一起滚，折叠后自动收起高度；
   overflow-x: hidden 防任何子元素隐形溢出在各面板底部冒出横向滚动条 */
.insp-scroll {
  flex: 1;
  min-height: 0;
  overflow-x: hidden;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 10px;
}

.insp-sec {
  flex-shrink: 0;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.07);
  background: rgba(255, 255, 255, 0.02);
  overflow: hidden;
}

.sec-head {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 7px;
  height: 36px;
  padding: 0 10px;
  border: none;
  background: none;
  color: rgba(255, 255, 255, 0.72);
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}
.sec-head:hover { background: rgba(0, 255, 255, 0.05); color: #fff; }
.sec-head > .fa-chevron-down,
.sec-head > .fa-chevron-right {
  width: 10px;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
}
.sec-icon { color: #00ffff; font-size: 12px; }
.sec-title { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* 折叠头右侧摘要：保持收起时也能看到关键状态 */
.sec-badge {
  max-width: 46%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 1px 7px;
  border-radius: 9px;
  font-size: 10px;
  font-weight: 500;
  background: rgba(0, 255, 255, 0.1);
  color: rgba(0, 255, 255, 0.85);
}

.sec-body {
  padding: 8px 10px 10px;
}
/* 树区限高：节点多时内部滚动，属性区紧接在下方 */
.sec-body--tree {
  height: 30vh;
  min-height: 160px;
  padding: 2px 6px 6px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}
.sec-body--form { padding: 8px 10px 10px; border-top: 1px solid rgba(255, 255, 255, 0.06); }
/* 属性区内的 NodeInspector 自带边框，段体已有容器，去掉双层边框 */
.sec-body :deep(.inspector) {
  border: none;
  border-radius: 0;
  background: transparent;
}
.sec-body--config { padding: 12px 10px 10px; border-top: 1px solid rgba(255, 255, 255, 0.06); }
/* 配置段头已显示「生成配置」，隐藏 ConfigPanel 内部重复标签 */
.sec-body--config :deep(.config-panel > .section-label) { display: none; }
/* 全局配置提示：多图时说明一份配置对全部页面生效 */
.cfg-shared {
  display: flex;
  align-items: flex-start;
  gap: 7px;
  margin: 0 0 10px;
  padding: 7px 9px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.14);
  background: rgba(0, 255, 255, 0.04);
  font-size: 11px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.45);
}
.cfg-shared i { margin-top: 2px; color: rgba(0, 255, 255, 0.55); }
.sec-body--rel { padding: 10px; border-top: 1px solid rgba(255, 255, 255, 0.06); }

/* 底部常驻生成按钮 */
.insp-foot {
  flex-shrink: 0;
  padding: 10px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(10, 14, 23, 0.9);
}
</style>
