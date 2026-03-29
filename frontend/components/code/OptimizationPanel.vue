<!-- components/optimization/OptimizationPanel.vue -->
<template>
  <div class="optimization-panel">
    <!-- 面板头部 -->
    <div class="panel-header">
      <div class="header-left">
        <span class="header-title">AI 智能优化</span>
        <span class="badge" v-if="selectedCount > 0">{{ selectedCount }}个优化项</span>
      </div>
      <button class="btn-optimize" @click.stop="startOptimization" :disabled="optimizing">
        <i class="fas fa-rocket"></i>
        {{ optimizing ? '优化中...' : '一键优化' }}
      </button>
    </div>

    <!-- 面板内容 -->
    <div class="panel-content">
      <!-- 快速优化选项 -->
      <div class="quick-options">
        <div
          v-for="opt in quickOptions"
          :key="opt.id"
          :class="['quick-option', { active: selectedOpt === opt.id }]"
          @click="selectQuickOpt(opt.id)"
        >
          <i :class="opt.icon"></i>
          <span>{{ opt.name }}</span>
        </div>
      </div>

      <!-- 自定义优化输入 -->
      <div class="custom-input-area">
        <div class="input-header">
          <label><i class="fas fa-edit"></i> 自定义优化要求</label>
          <span class="char-count">{{ customText.length }} / 500</span>
        </div>
        <div class="textarea-wrapper">
          <textarea
            ref="textareaRef"
            v-model="customText"
            class="custom-input"
            rows="3"
            maxlength="600"
            placeholder="💡 例如：添加暗色模式支持、优化按钮动画效果..."
          ></textarea>
          <div class="resize-handle" @mousedown="startResize">
            <i class="fas fa-arrows-alt-v"></i>
            <span>拖动调整高度</span>
          </div>
        </div>
      </div>

      <!-- AI 检测的问题列表 -->
      <div class="suggestions-area">
        <div class="suggestions-header">
          <div class="header-left-part">
            <div class="icon-wrapper"><i class="fas fa-robot"></i></div>
            <div>
              <h4>AI 智能诊断</h4>
              <p class="header-desc">发现 {{ issues.length }} 个可优化项</p>
            </div>
          </div>
          <button class="select-all-btn" @click="toggleAll">
            <i :class="allChecked ? 'fas fa-check-circle' : 'far fa-circle'"></i>
            <span>{{ allChecked ? '取消全选' : '全选' }}</span>
          </button>
        </div>

        <div class="suggestions-list">
          <div
            v-for="(item, idx) in issues"
            :key="idx"
            class="suggestion-card"
            :class="{ checked: checkedIssues.includes(item.text) }"
            @click="toggleIssue(item.text)"
          >
            <div class="card-check">
              <i :class="checkedIssues.includes(item.text) ? 'fas fa-check-circle' : 'far fa-circle'"></i>
            </div>
            <div class="card-content">
              <div class="card-text">{{ item.text }}</div>
              <div class="card-meta">
                <span class="suggestion-tag" :class="item.type">
                  <i :class="getTagIcon(item.type)"></i> {{ item.tag }}
                </span>
                <span class="suggestion-impact">
                  <i class="fas fa-chart-line"></i> {{ getImpactText(item.type) }}
                </span>
              </div>
            </div>
            <!-- <div class="card-action"><i class="fas fa-chevron-right"></i></div> -->
          </div>
        </div>
      </div>
    </div>

    <!-- 进度弹窗 -->
    <ProgressModal
      :visible="showProgressModal"
      :progress="progress"
      :progress-msg="progressMsg"
      :optimizing="isOptimizing"
      @close="closeProgressModal"
      @complete="onProgressComplete"
    />

    <!-- 对比弹窗 -->
    <CompareModal
      :visible="showCompareModal"
      :original-code="originalCode"
      :optimized-code="optimizedCode"
      :original-score="originalScore"
      :optimized-score="optimizedScore"
      :diff-stats="diffStats"
      @close="closeCompareModal"
      @accept="onAcceptOptimize"
      @reject="onRejectOptimize"
    />

    <!-- 提示组件 -->
    <ToastNotification ref="toastRef" />
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import CompareModal from "/components/optimization/CompareModal.vue";
import ProgressModal from "/components/optimization/ProgressModal.vue";
import ToastNotification from "/components/common/ToastNotification.vue";

const props = defineProps({
  originalCode: { type: String, default: "" }
});
const emit = defineEmits(["accepted"]);

// ========== 状态 ==========
const optimizing = ref(false);
const selectedOpt = ref(null);
const customText = ref("");
const checkedIssues = ref([]);
const optimizedCode = ref("");

const showProgressModal = ref(false);
const showCompareModal = ref(false);
const isOptimizing = ref(false);
const progress = ref(0);
const progressMsg = ref("");

const originalScore = ref(78);
const optimizedScore = ref(92);
const diffStats = ref({ added: 12, removed: 5 });

const textareaRef = ref(null);
let isResizing = false;
let startY = 0;
let startHeight = 0;

const quickOptions = [
  { id: "performance", name: "性能优化", icon: "fas fa-tachometer-alt" },
  { id: "readability", name: "可读性优化", icon: "fas fa-book" },
  { id: "responsive", name: "响应式优化", icon: "fas fa-mobile-alt" },
  { id: "clean", name: "代码清理", icon: "fas fa-broom" }
];

const issues = ref([
  { text: "使用 const 替代 final 提升性能", tag: "性能", type: "performance" },
  { text: "添加缺失的类型注解", tag: "类型", type: "readability" },
  { text: "提取重复代码为独立方法", tag: "重构", type: "clean" },
  { text: "添加响应式布局适配", tag: "响应式", type: "responsive" }
]);

const selectedCount = computed(() => {
  let count = 0;
  if (selectedOpt.value) count++;
  if (customText.value.trim()) count++;
  count += checkedIssues.value.length;
  return count;
});

const allChecked = computed(() => checkedIssues.value.length === issues.value.length);

// 拉伸功能
const startResize = (e) => {
  e.preventDefault();
  isResizing = true;
  startY = e.clientY;
  startHeight = textareaRef.value?.offsetHeight || 100;
  document.addEventListener("mousemove", onResize);
  document.addEventListener("mouseup", stopResize);
  document.body.style.userSelect = "none";
  document.body.style.cursor = "ns-resize";
};

const onResize = (e) => {
  if (!isResizing) return;
  const deltaY = e.clientY - startY;
  let newHeight = startHeight + deltaY;
  newHeight = Math.min(Math.max(newHeight, 100), 500);
  if (textareaRef.value) {
    textareaRef.value.style.height = `${newHeight}px`;
    if (deltaY > 0 && newHeight > startHeight) {
      const rect = textareaRef.value.getBoundingClientRect();
      const bottomOffset = rect.bottom - window.innerHeight;
      if (bottomOffset > 0) {
        window.scrollBy({ top: bottomOffset + 20, behavior: "smooth" });
      }
    }
  }
};

const stopResize = () => {
  isResizing = false;
  document.removeEventListener("mousemove", onResize);
  document.removeEventListener("mouseup", stopResize);
  document.body.style.userSelect = "";
  document.body.style.cursor = "";
};

const getTagIcon = (type) => {
  const icons = {
    performance: "fas fa-tachometer-alt",
    readability: "fas fa-book",
    responsive: "fas fa-mobile-alt",
    clean: "fas fa-broom"
  };
  return icons[type] || "fas fa-tag";
};

const getImpactText = (type) => {
  const impacts = {
    performance: "提升 15% 性能",
    readability: "提高可维护性",
    responsive: "适配更多设备",
    clean: "减少代码量"
  };
  return impacts[type] || "优化代码质量";
};

const selectQuickOpt = (id) => { selectedOpt.value = selectedOpt.value === id ? null : id; };
const toggleAll = () => {
  if (allChecked.value) checkedIssues.value = [];
  else checkedIssues.value = issues.value.map(i => i.text);
};
const toggleIssue = (text) => {
  const index = checkedIssues.value.indexOf(text);
  if (index > -1) checkedIssues.value.splice(index, 1);
  else checkedIssues.value.push(text);
};

const startOptimization = async () => {
  showProgressModal.value = true;
  isOptimizing.value = true;
  progress.value = 0;
  progressMsg.value = "正在分析代码...";

  const steps = [
    { p: 20, msg: "分析代码结构..." },
    { p: 40, msg: "识别可优化点..." },
    { p: 60, msg: "应用优化策略..." },
    { p: 80, msg: "验证优化结果..." },
    { p: 100, msg: "优化完成！" }
  ];

  for (const step of steps) {
    await new Promise(r => setTimeout(r, 600));
    progress.value = step.p;
    progressMsg.value = step.msg;
  }

  let newCode = props.originalCode;
  if (selectedOpt.value === "performance") newCode = newCode.replace(/final /g, "const ");
  if (selectedOpt.value === "readability") newCode = "/// AI 优化后的代码\n" + newCode;
  if (selectedOpt.value === "responsive") newCode = newCode + "\n\n// 响应式适配\nfinal screenWidth = MediaQuery.of(context).size.width;";
  if (customText.value.trim()) newCode = `// 优化要求: ${customText.value}\n${newCode}`;
  if (checkedIssues.value.length > 0) newCode = `// 已应用优化:\n// ${checkedIssues.value.join("\n// ")}\n\n${newCode}`;

  optimizedCode.value = newCode;
  isOptimizing.value = false;
};

const onProgressComplete = () => {
  showProgressModal.value = false;
  showCompareModal.value = true;
};

const closeProgressModal = () => {
  showProgressModal.value = false;
  isOptimizing.value = false;
};

const closeCompareModal = () => {
  showCompareModal.value = false;
};

const onAcceptOptimize = (code) => {
  showCompareModal.value = false;
  if (toastRef.value) {
    toastRef.value.show({
      type: 'success',
      title: '优化成功',
      message: '优化代码已应用，可在生成历史中查看',
      duration: 3000
    });
  }
  emit("accepted", code);
  selectedOpt.value = null;
  customText.value = "";
  checkedIssues.value = [];
};

const onRejectOptimize = () => {
  showCompareModal.value = false;
  if (toastRef.value) {
    toastRef.value.show({
      type: 'info',
      title: '已取消',
      message: '优化已取消，可随时重新优化',
      duration: 2000
    });
  }
  optimizedCode.value = "";
};

const toastRef = ref();
</script>

<style scoped>
.optimization-panel {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 20px;
  margin-top: 20px;
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: rgba(0, 0, 0, 0.3);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
}

.badge {
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.2), rgba(255, 0, 255, 0.2));
  padding: 2px 8px;
  border-radius: 20px;
  font-size: 12px;
  color: #00ffff;
}

.btn-optimize {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 20px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  border-radius: 40px;
  color: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-optimize:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 20px rgba(0, 255, 255, 0.3);
}

.btn-optimize:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.panel-content {
  padding: 20px;
  border-top: 1px solid rgba(0, 255, 255, 0.1);
}

/* 快速选项 */
.quick-options {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.quick-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
  text-align: center;
}

.quick-option i {
  font-size: 20px;
  color: #888;
}

.quick-option span {
  font-size: 12px;
  color: #aaa;
}

.quick-option:hover {
  border-color: #00ffff;
  background: rgba(0, 255, 255, 0.05);
}

.quick-option.active {
  border-color: #00ffff;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.15), rgba(255, 0, 255, 0.15));
}

.quick-option.active i,
.quick-option.active span {
  color: #00ffff;
}

/* 自定义输入 */
.custom-input-area {
  margin-bottom: 24px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 16px;
  padding: 16px;
}

.input-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.input-header label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #00ffff;
}

.char-count {
  font-size: 11px;
  color: #888;
  background: rgba(0, 0, 0, 0.3);
  padding: 2px 8px;
  border-radius: 20px;
}

.textarea-wrapper {
  position: relative;
}

.custom-input {
  width: 100%;
  padding: 14px 16px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(0, 255, 255, 0.25);
  border-radius: 12px;
  color: #fff;
  font-size: 13px;
  resize: vertical;
  min-height: 150px;
  font-family: inherit;
}

.custom-input:focus {
  outline: none;
  border-color: #00ffff;
}

.resize-handle {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 40px;
  height: 40px;
  cursor: ns-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 12px 0 12px 0;
  font-size: 10px;
  color: #00ffff;
  opacity: 0.6;
  z-index: 5;
}

.resize-handle:hover {
  opacity: 1;
  width: 90px;
}

.resize-handle span {
  display: none;
}

.resize-handle:hover span {
  display: inline;
}

/* 问题列表 */
.suggestions-area {
  margin-bottom: 24px;
}

.suggestions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(0, 255, 255, 0.15);
}

.header-left-part {
  display: flex;
  align-items: center;
  gap: 14px;
}

.icon-wrapper {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.15), rgba(255, 0, 255, 0.15));
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(0, 255, 255, 0.3);
}

.icon-wrapper i {
  font-size: 22px;
  color: #00ffff;
}

.header-left-part h4 {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 4px;
}

.header-desc {
  font-size: 12px;
  color: #888;
}

.select-all-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(0, 255, 255, 0.08);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 40px;
  color: #00ffff;
  font-size: 13px;
  cursor: pointer;
}

.select-all-btn:hover {
  background: rgba(0, 255, 255, 0.15);
}

.suggestions-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 400px;
  overflow-y: auto;
}

.suggestion-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  background: rgba(0, 0, 0, 0.25);
  border: 1px solid rgba(0, 255, 255, 0.1);
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.suggestion-card:hover {
  background: rgba(0, 255, 255, 0.05);
  border-color: rgba(0, 255, 255, 0.3);
}

.suggestion-card.checked {
  background: rgba(0, 255, 255, 0.08);
  border-color: rgba(0, 255, 255, 0.4);
}

.card-check i {
  font-size: 20px;
  color: #00ffff;
}

.card-content {
  flex: 1;
}

.card-text {
  font-size: 13px;
  color: #e0e0e0;
  margin-bottom: 6px;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.suggestion-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 10px;
  font-weight: 500;
}

.suggestion-tag.performance {
  background: rgba(0, 255, 255, 0.2);
  color: #00ffff;
}

.suggestion-tag.readability {
  background: rgba(255, 0, 255, 0.2);
  color: #ff00ff;
}

.suggestion-tag.responsive {
  background: rgba(0, 255, 0, 0.2);
  color: #0f0;
}

.suggestion-tag.clean {
  background: rgba(255, 255, 0, 0.2);
  color: #ff0;
}

.suggestion-impact {
  font-size: 10px;
  color: #888;
  display: flex;
  align-items: center;
  gap: 4px;
}

.card-action i {
  color: #444;
  font-size: 12px;
  transition: all 0.3s;
}

.suggestion-card:hover .card-action i {
  color: #00ffff;
  transform: translateX(2px);
}

/* 滚动条 */
.suggestions-list::-webkit-scrollbar {
  width: 4px;
}

.suggestions-list::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.suggestions-list::-webkit-scrollbar-thumb {
  background: rgba(0, 255, 255, 0.3);
  border-radius: 4px;
}

/* 响应式 */
@media (max-width: 768px) {
  .quick-options {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>