<!-- pages/code.vue -->
<template>
  <div class="page-container">
    <main class="main-content">
      <div class="page-header">
        <h1>AI 代码生成</h1>
        <p>上传设计稿 · AI 智能识别 · 秒级生成高质量代码</p>
      </div>

      <!-- ═══ 上传区 ═══ -->
      <UploaderImage @generated="onGenerated" />

      <!-- ═══ 未生成时的占位 ═══ -->
      <div v-if="!hasGenerated" class="empty-state">
        <div class="empty-icon">
          <i class="fas fa-code"></i>
        </div>
        <h3>等待生成代码</h3>
        <p>上传设计稿并点击「开始生成代码」，AI 将自动为你生成高质量代码</p>
        <div class="feature-hints">
          <div class="hint-item">
            <i class="fas fa-image"></i>
            <span>支持多张设计稿上传</span>
          </div>
          <div class="hint-item">
            <i class="fas fa-edit"></i>
            <span>为每张图片添加描述让 AI 更懂你</span>
          </div>
          <div class="hint-item">
            <i class="fas fa-mobile-alt"></i>
            <span>一键生成 Flutter / React / Vue 代码</span>
          </div>
        </div>
      </div>

      <!-- ═══ 生成结果展示 ═══ -->
      <template v-if="hasGenerated">
        <div class="core-layout">
          <div class="panel editor-panel">
            <div class="panel-header">
              <div class="panel-title">
                <i class="fas fa-code"></i>
                <span>{{ generatedLang }} 代码</span>
              </div>
              <div class="panel-actions">
                <button class="act-btn" @click="handleCopy"><i class="fas fa-copy"></i><span>复制</span></button>
                <button class="act-btn" @click="handleFormat"><i class="fas fa-magic"></i><span>格式化</span></button>
                <button class="act-btn ghost" @click="handleDownload"><i class="fas fa-download"></i><span>下载</span></button>
              </div>
            </div>
            <div class="panel-body editor-body">
              <CodeEditor
                v-model="generatedCode"
                :language="codeLanguage"
                :readonly="false"
                height="100%"
                placeholder="// AI 生成的代码将在这里展示..."
              />
            </div>
          </div>

          <div class="panel preview-panel">
            <div class="panel-header">
              <div class="panel-title">
                <i class="fas fa-mobile-alt"></i>
                <span>实时预览</span>
              </div>
              <span class="device-badge">iPhone 15 Pro</span>
            </div>
            <div class="panel-body preview-body">
              <PhoneMockup :html="phoneHtml" />
            </div>
          </div>
        </div>

        <div class="optimization-section">
          <OptimizationPanel :original-code="generatedCode" @accepted="onOptimized" />
        </div>
      </template>
    </main>
  </div>
</template>

<script setup>
import CodeEditor from '~/components/code/CodeEditor.vue'
import PhoneMockup from '~/components/screen/PhoneMockup.vue'
import OptimizationPanel from '~/components/code/OptimizationPanel.vue'
import UploaderImage from '~/components/upload/UploaderImage.vue'

// ═══ 生成状态 ═══
const hasGenerated = ref(false)
const generatedCode = ref('')
const generatedLang = ref('Dart')
const codeLanguage = ref('dart')

const langMap = {
  flutter: { label: 'Dart', lang: 'dart' },
  react: { label: 'TypeScript', lang: 'typescript' },
  vue: { label: 'Vue', lang: 'html' }
}

const onGenerated = (result) => {
  hasGenerated.value = true
  generatedCode.value = result.code
  if (result.framework && langMap[result.framework]) {
    generatedLang.value = langMap[result.framework].label
    codeLanguage.value = langMap[result.framework].lang
  }
  updatePreview()
}

// ═══ 手机预览 ═══
const phoneStyles = `<style>
.phone-header{background:#fff;padding:10px 14px 8px;border-bottom:1px solid #e5e5ea;flex-shrink:0}
.header-top{display:flex;justify-content:space-between;align-items:center;margin-bottom:8px}
.logo-text{font-size:18px;font-weight:800;color:#1c1c1e;letter-spacing:-0.3px}
.header-actions{display:flex;gap:14px}
.header-actions i{font-size:16px;color:#1c1c1e}
.search-bar{background:#f2f2f7;border-radius:10px;padding:8px 14px;display:flex;align-items:center;gap:8px;color:#8e8e93;font-size:13px}
.search-bar i{color:#8e8e93;font-size:14px}
.category-tabs{display:flex;gap:4px;margin-top:8px;overflow-x:auto;scrollbar-width:none}
.category-tabs::-webkit-scrollbar{display:none}
.cat-item{padding:5px 14px;border-radius:16px;font-size:12px;color:#8e8e93;white-space:nowrap;transition:0.2s}
.cat-item.active{background:#1c1c1e;color:#fff;font-weight:600}
.phone-body{flex:1;overflow-y:auto;padding:12px;background:#f2f2f7;scrollbar-width:none}
.phone-body::-webkit-scrollbar{display:none}
.hero-banner{background:linear-gradient(135deg,#ff3b6e,#ff6b3d);border-radius:14px;padding:16px;display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;overflow:hidden;position:relative}
.hero-banner::after{content:'';position:absolute;right:-30px;top:-30px;width:100px;height:100px;border-radius:50%;background:rgba(255,255,255,0.12)}
.banner-content{position:relative;z-index:1}
.banner-tag{display:inline-block;background:rgba(255,255,255,0.25);border-radius:4px;padding:2px 8px;font-size:10px;font-weight:700;color:#fff;margin-bottom:6px}
.hero-banner h4{font-size:18px;font-weight:800;color:#fff;margin:0 0 2px}
.hero-banner p{font-size:11px;color:rgba(255,255,255,0.85);margin:0 0 10px}
.banner-btn{background:#fff;border:none;border-radius:14px;padding:5px 14px;font-size:11px;font-weight:700;color:#ff3b6e;cursor:default}
.banner-visual{position:relative;z-index:1;width:64px;height:64px;border-radius:10px;overflow:hidden;background:rgba(255,255,255,0.18);flex-shrink:0}
.banner-visual img{width:100%;height:100%;object-fit:cover}
.section-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:10px}
.section-title{font-size:15px;font-weight:700;color:#1c1c1e}
.section-more{font-size:12px;color:#8e8e93}
.product-grid{display:grid;grid-template-columns:1fr 1fr;gap:8px}
.product-card{background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,0.06)}
.product-img{width:100%;aspect-ratio:4/5;overflow:hidden;background:#f2f2f7}
.product-img img{width:100%;height:100%;object-fit:cover}
.product-info{padding:8px 10px 10px}
.product-name{font-size:12px;color:#1c1c1e;line-height:1.3;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.product-price-row{display:flex;justify-content:space-between;align-items:baseline;margin-top:4px}
.product-price{font-size:15px;font-weight:800;color:#ff3b6e}
.product-sold{font-size:10px;color:#c7c7cc}
<\\/style>`

const products = [
  { name: 'Air Max 270', price: '¥899', img: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=200&h=250&fit=crop' },
  { name: '简约白T恤', price: '¥159', img: 'https://images.unsplash.com/photo-1583743814966-8936f5b7be1a?w=200&h=250&fit=crop' },
  { name: '智能手表', price: '¥1299', img: 'https://images.unsplash.com/photo-1524592094714-0f0654e20314?w=200&h=250&fit=crop' },
  { name: '降噪耳机', price: '¥599', img: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=200&h=250&fit=crop' },
]

const buildProductsHtml = () => products.map(p => `
  <div class="product-card">
    <div class="product-img"><img src="${p.img}" alt="${p.name}" loading="lazy" /></div>
    <div class="product-info">
      <div class="product-name">${p.name}</div>
      <div class="product-price-row">
        <span class="product-price">${p.price}</span>
        <span class="product-sold">已售1.2k</span>
      </div>
    </div>
  </div>
`).join('')

const phoneHtml = ref('')

const updatePreview = () => {
  phoneHtml.value = `${phoneStyles}<div class="phone-header"><div class="header-top"><div class="logo-text">StyleHub</div><div class="header-actions"><i class="fas fa-bell"></i><i class="fas fa-shopping-bag"></i></div></div><div class="search-bar"><i class="fas fa-search"></i><span>搜索潮流好物...</span></div><div class="category-tabs"><span class="cat-item active">推荐</span><span class="cat-item">服饰</span><span class="cat-item">鞋靴</span><span class="cat-item">数码</span><span class="cat-item">家居</span></div></div><div class="phone-body"><div class="hero-banner"><div class="banner-content"><div class="banner-tag">限时特惠</div><h4>夏日焕新季</h4><p>全场低至 5 折 · 满299包邮</p><button class="banner-btn">立即抢购 →</button></div><div class="banner-visual"><img src="https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da?w=140&h=140&fit=crop" alt="" /></div></div><div class="section-header"><span class="section-title">🔥 为你推荐</span><span class="section-more">查看更多 →</span></div><div class="product-grid">${buildProductsHtml()}</div></div>`
}

const handleCopy = async () => {
  try { await navigator.clipboard.writeText(generatedCode.value); ElMessage.success('已复制到剪贴板') }
  catch { ElMessage.error('复制失败') }
}

const handleFormat = () => {
  let indent = 0
  generatedCode.value = generatedCode.value.split('\n').map(line => {
    const t = line.trim()
    if (!t) return ''
    if (/^[})]/.test(t)) indent = Math.max(0, indent - 1)
    const r = '  '.repeat(indent) + t
    if (/[({]$/.test(t)) indent++
    return r
  }).join('\n')
  ElMessage.success('代码已格式化')
}

const handleDownload = () => {
  const ext = codeLanguage.value === 'dart' ? 'dart' : codeLanguage.value === 'typescript' ? 'tsx' : 'vue'
  const blob = new Blob([generatedCode.value], { type: 'text/plain' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `generated.${ext}`
  a.click()
  URL.revokeObjectURL(a.href)
  ElMessage.success('已下载')
}

const onOptimized = (code) => { generatedCode.value = code }
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0a0f 0%, #0f1a1f 100%);
}
.main-content {
  max-width: 1600px;
  margin: 0 auto;
  padding: 40px 24px 60px;
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
}
.page-header h1 {
  font-size: 40px;
  font-weight: 800;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 10px;
}
.page-header p {
  color: #6b7280;
  font-size: 15px;
}

/* ═══ 空状态 ═══ */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
  background: rgba(15, 20, 30, 0.35);
  border: 1px solid rgba(0, 255, 255, 0.1);
  border-radius: 24px;
  margin-top: 24px;
}
.empty-icon {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(0,255,255,0.12), rgba(255,0,255,0.12));
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
}
.empty-icon i {
  font-size: 36px;
  color: #00ffff;
}
.empty-state h3 {
  font-size: 22px;
  color: #ccc;
  margin-bottom: 10px;
}
.empty-state > p {
  color: #6b7280;
  font-size: 14px;
  margin-bottom: 36px;
}
.feature-hints {
  display: flex;
  gap: 32px;
  flex-wrap: wrap;
  justify-content: center;
}
.hint-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background: rgba(0,255,255,0.05);
  border: 1px solid rgba(0,255,255,0.12);
  border-radius: 14px;
  color: #888;
  font-size: 13px;
}
.hint-item i {
  color: #00cfff;
  font-size: 16px;
}

/* ═══ 核心区 ═══ */
.core-layout {
  display: flex;
  gap: 24px;
  height: 620px;
  margin-top: 24px;
  margin-bottom: 32px;
}

.panel {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.18);
  border-radius: 18px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 18px;
  background: rgba(0, 0, 0, 0.3);
  border-bottom: 1px solid rgba(0, 255, 255, 0.1);
  flex-shrink: 0;
}
.panel-title {
  display: flex;
  align-items: center;
  gap: 9px;
  font-size: 14px;
  font-weight: 600;
  color: #00cfff;
}
.panel-title i { font-size: 15px; }
.panel-actions { display: flex; gap: 6px; }

.act-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 12px;
  background: rgba(0, 255, 255, 0.07);
  border: 1px solid rgba(0, 255, 255, 0.18);
  border-radius: 18px;
  color: #00cfff;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}
.act-btn i { font-size: 12px; }
.act-btn:hover { background: rgba(0, 255, 255, 0.14); border-color: #00ffff; }
.act-btn.ghost { background: rgba(255,255,255,0.04); border-color: rgba(255,255,255,0.12); color: #aaa; }
.act-btn.ghost:hover { background: rgba(255,255,255,0.08); color: #fff; }

.device-badge {
  font-size: 11px;
  color: #6b7280;
  letter-spacing: 0.3px;
  background: rgba(255,255,255,0.04);
  padding: 3px 10px;
  border-radius: 20px;
}

.panel-body { overflow: hidden; flex: 1; min-height: 0; }

.editor-panel { flex: 1.55; min-width: 0; }
.editor-body { flex: 1; min-height: 0; }

.preview-panel { flex: 1; min-width: 0; }
.preview-body {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 16px;
  background: rgba(0, 0, 0, 0.15);
  overflow-y: auto;
}

.optimization-section { margin-top: 8px; }

@media (max-width: 1100px) {
  .core-layout { flex-direction: column; height: auto; }
  .editor-panel { height: 480px; }
  .preview-panel { height: auto; }
  .preview-body { justify-content: center; }
  .feature-hints { flex-direction: column; align-items: center; }
}

@media (max-width: 768px) {
  .main-content { padding: 20px 14px 40px; }
  .page-header h1 { font-size: 28px; }
  .editor-panel { height: 380px; }
  .act-btn span { display: none; }
}
</style>
