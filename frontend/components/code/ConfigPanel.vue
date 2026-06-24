<template>
  <div class="config-panel">
    <div class="section-label"><span class="label-num">2</span> <span class="section-label-text">生成配置</span></div>
    <!-- 目标框架 -->
    <div class="config-group">
      <div class="config-label-row">
        <span class="config-label">目标框架</span>
      </div>
      <div class="config-chips">
        <button
          v-for="fw in frameworks"
          :key="fw.value"
          class="chip"
          :class="modelValue.framework === fw.value ? 'chip--active' : ''"
          @click="$emit('update:modelValue', { ...modelValue, framework: fw.value })"
        >
          {{ fw.label }}
        </button>
      </div>
    </div>

    <!-- 目标平台 -->
    <div class="config-group">
      <div class="config-label-row">
        <span class="config-label">目标平台</span>
      </div>
      <div class="config-chips">
        <button
          v-for="pl in platforms"
          :key="pl.value"
          class="chip"
          :class="modelValue.platform === pl.value ? 'chip--active' : ''"
          @click="$emit('update:modelValue', { ...modelValue, platform: pl.value })"
        >
          <i :class="pl.icon" class="chip-icon"></i>
          {{ pl.label }}
        </button>
      </div>
    </div>

    <!-- 生成选项 -->
    <div class="config-group">
      <div class="config-label-row">
        <span class="config-label">生成选项</span>
      </div>
      <div class="config-grid">
        <label
          v-for="opt in options"
          :key="opt.value"
          class="option-item"
          :class="modelValue.options.includes(opt.value) ? 'option-item--active' : ''"
          @click="toggleOption(opt.value)"
        >
          <i
            :class="modelValue.options.includes(opt.value) ? 'fa-regular fa-square-check' : 'fa-regular fa-square'"
            class="option-check"
          ></i>
          {{ opt.label }}
        </label>
      </div>
    </div>

    <!-- 组件库选择（勾选"使用组件库"后出现） -->
    <div v-if="modelValue.options.includes('component')" class="config-group">
      <div class="config-label-row">
        <span class="config-label lib-sub-label">选择组件库</span>
      </div>
      <div class="config-chips">
        <button
          v-for="lib in currentLibs"
          :key="lib.value"
          class="chip"
          :class="(modelValue.componentLib || defaultLib) === lib.value ? 'chip--active' : ''"
          @click="$emit('update:modelValue', { ...modelValue, componentLib: lib.value })"
        >
          {{ lib.label }}
        </button>
      </div>
    </div>

    <!-- 高级选项 -->
    <div class="config-group">
      <div class="config-label-row">
        <span class="config-label">高级选项</span>
      </div>
      <div class="config-chips">
        <button
          v-for="adv in advanced"
          :key="adv.value"
          class="chip"
          :class="modelValue.advanced.includes(adv.value) ? 'chip--active' : ''"
          @click="toggleAdvanced(adv.value)"
        >
          {{ adv.label }}
        </button>
      </div>
    </div>

    <!-- 开始生成 -->
    <button class="generate-btn" @click="$emit('generate')">
      <i class="fa-solid fa-wand-magic-sparkles"></i>
      开始生成代码
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      framework: 'Flutter',
      platform: 'mobile',
      options: ['responsive'],
      advanced: [],
      componentLib: 'material',
    }),
  },
})
const emit = defineEmits(['update:modelValue', 'generate'])

const frameworks = [
  { label: 'Flutter', value: 'Flutter' },
  { label: 'React', value: 'React' },
  { label: 'Vue', value: 'Vue' },
]
const platforms = [
  { label: '手机端', value: 'mobile', icon: 'fa-solid fa-mobile-screen-button' },
  { label: '桌面端', value: 'desktop', icon: 'fa-solid fa-desktop' },
]
const options = [
  { label: '生成响应式布局', value: 'responsive' },
  { label: '添加注释', value: 'comment' },
  { label: '使用组件库', value: 'component' },
]
const advanced = [
  { label: '优化性能', value: 'perf' },
  { label: '生成文档', value: 'docs' },
]

const toggleOption = (val) => {
  const opts = [...props.modelValue.options]
  const idx = opts.indexOf(val)
  idx === -1 ? opts.push(val) : opts.splice(idx, 1)
  emit('update:modelValue', { ...props.modelValue, options: opts })
}
const toggleAdvanced = (val) => {
  const adv = [...props.modelValue.advanced]
  const idx = adv.indexOf(val)
  idx === -1 ? adv.push(val) : adv.splice(idx, 1)
  emit('update:modelValue', { ...props.modelValue, advanced: adv })
}

// 各框架对应的组件库选项
const componentLibs = {
  Flutter: [
    { label: 'Material Design', value: 'material' },
  ],
  React: [
    { label: 'Ant Design', value: 'antd' },
    { label: 'Material UI', value: 'mui' },
  ],
  Vue: [
    { label: 'Element Plus', value: 'element-plus' },
    { label: 'Ant Design Vue', value: 'antdv' },
    { label: 'Naive UI', value: 'naive-ui' },
  ],
}

const currentLibs = computed(() => {
  return componentLibs[props.modelValue.framework] || componentLibs.Vue
})

const defaultLib = computed(() => currentLibs.value[0]?.value)
</script>

<style scoped>
.config-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-label {
  font-size: 14px;
  font-weight: 600;
  color: #60a5fa;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  margin-bottom: 16px;
  position: relative;
  display: inline-block;
  padding-left: 14px;
}
.section-label::before {
  content: '';
  position: absolute;
  left: 0;
  top: 2px;
  bottom: 2px;
  width: 2px;
  background: linear-gradient(135deg, #60a5fa, #818cf8);
  border-radius: 1px;
}
.section-label-text {
  font-weight: 700;
  color: #ffffff;
}

.label-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 8px;
  background: #ffffff;
  color: #60a5fa;
  font-size: 12px;
  font-weight: 700;
  margin-right: 6px;
}

.config-group {
  /* no extra styles needed */
}

.config-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}
.config-label {
  font-size: 13px;
  font-weight: 500;
  color: #ffffff;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.lib-sub-label {
  font-size: 12px;
  color: #ffffff;
  text-transform: none;
  letter-spacing: normal;
}

.label-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 8px;
  background: #60a5fa;
  color: #ffffff;
  font-size: 12px;
  font-weight: 700;
}

.config-chips {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.chip {
  padding: 6px 16px;
  border-radius: 20px;
  border: 1px solid rgba(96, 165, 250, 0.12);
  background: rgba(255, 255, 255, 0.02);
  color: rgba(255, 255, 255, 0.7);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}
.chip:hover {
  border-color: rgba(96, 165, 250, 0.25);
  color: rgba(255, 255, 255, 0.85);
}
.chip--active {
  background: rgba(96, 165, 250, 0.1);
  color: #60a5fa;
  border-color: rgba(96, 165, 250, 0.25);
}
.chip-icon {
  margin-right: 4px;
}

.config-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid rgba(96, 165, 250, 0.10);
  background: rgba(255, 255, 255, 0.02);
  color: rgba(255, 255, 255, 0.65);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}
.option-item:hover {
  border-color: rgba(96, 165, 250, 0.20);
}
.option-item--active {
  background: rgba(96, 165, 250, 0.08);
  color: rgba(255, 255, 255, 0.85);
  border-color: rgba(96, 165, 250, 0.20);
}
.option-check {
  color: rgba(255, 255, 255, 0.5);
  font-size: 13px;
}
.option-item--active .option-check {
  color: #60a5fa;
}

.generate-btn {
  width: 100%;
  padding: 12px 0;
  border-radius: 8px;
  border: none;
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: linear-gradient(135deg, #60a5fa, #818cf8);
  box-shadow: 0 4px 20px rgba(96, 165, 250, 0.25);
  transition: all 0.3s;
}
.generate-btn:hover {
  box-shadow: 0 6px 30px rgba(96, 165, 250, 0.4);
  transform: translateY(-1px);
}
</style>
