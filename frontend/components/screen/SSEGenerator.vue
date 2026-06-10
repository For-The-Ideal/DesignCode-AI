<template>
  <div class="screen" ref="screenRef">
    <div class="content">
      <div class="header-title">
        <h2>AI实时代码生成</h2>
      </div>
      <div class="sse-layout">
        <!-- 左侧代码区 -->
        <div class="code-area">
          <CodeEditor
            v-model="generatedCode"
            language="dart"
            :readonly="true"
            height="700px"
            placeholder="// 滚动到此处将自动开始流式生成 Flutter 代码..."
          />
        </div>
        <!-- 右侧手机模拟器 -->
        <div class="preview-area">
          <FlutterTemplate :html="phoneHtml" />
        </div>
      </div>
      <div class="scroll-hint"><span>▼ 继续滚动查看 AI 智能诊断 ▼</span></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import CodeEditor from '@/components/code/CodeEditor.vue';
import FlutterTemplate from '@/components/common/FlutterTemplate.vue';

const screenRef = ref(null);
const generatedCode = ref('');
// 电商商品数据
const shopProducts = [
  { name: 'Air Max 270', price: '¥899', emoji: '👟', img: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=200&h=250&fit=crop' },
  { name: '简约白T恤', price: '¥159', emoji: '👕', img: 'https://images.unsplash.com/photo-1583743814966-8936f5b7be1a?w=200&h=250&fit=crop' },
  { name: '智能运动手表', price: '¥1299', emoji: '⌚', img: 'https://images.unsplash.com/photo-1524592094714-0f0654e20314?w=200&h=250&fit=crop' },
  { name: '降噪耳机', price: '¥599', emoji: '🎧', img: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=200&h=250&fit=crop' },
  { name: '休闲双肩包', price: '¥299', emoji: '🎒', img: 'https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=200&h=250&fit=crop' },
  { name: '复古太阳镜', price: '¥189', emoji: '🕶️', img: 'https://images.unsplash.com/photo-1572635196237-14b3f281503f?w=200&h=250&fit=crop' },
  { name: '运动水壶', price: '¥79', emoji: '🧴', img: 'https://images.unsplash.com/photo-1548839140-29a749e1cf4d?w=200&h=250&fit=crop' },
  { name: '蓝牙音箱', price: '¥399', emoji: '🔊', img: 'https://images.unsplash.com/photo-1589003077984-894e133dabab?w=200&h=250&fit=crop' },
  { name: '经典帆布鞋', price: '¥259', emoji: '👟', img: 'https://images.unsplash.com/photo-1549298916-b41d501d3772?w=200&h=250&fit=crop' },
  { name: '潮流棒球帽', price: '¥129', emoji: '🧢', img: 'https://images.unsplash.com/photo-1588850561407-ed78c282e36b?w=200&h=250&fit=crop' },
];

const productsDartList = shopProducts.map(p =>
  `    Product('${p.name}', '${p.price}', '${p.emoji}'),`
).join('\n');

const fullFlutterCode = `import 'package:flutter/material.dart';

void main() => runApp(const MyApp());

// ── 数据模型 ──
class Product {
  final String name;
  final String price;
  final String emoji;
  const Product(this.name, this.price, this.emoji);
}

// ── 根组件 ──
class MyApp extends StatelessWidget {
  const MyApp({super.key});
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '潮流精选',
      debugShowCheckedModeBanner: false,
      theme: ThemeData.dark().copyWith(
        primaryColor: Colors.cyan,
        scaffoldBackgroundColor: const Color(0xFF0C1222),
      ),
      home: const ShopHomePage(),
    );
  }
}

// ── 商品卡片组件 ──
class ProductCard extends StatelessWidget {
  final Product product;
  const ProductCard({super.key, required this.product});
  @override
  Widget build(BuildContext context) {
    return Card(
      color: const Color(0xFF141E32),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Center(
              child: Text(product.emoji, style: const TextStyle(fontSize: 44)),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(10),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(product.name,
                  style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w500, color: Colors.white),
                  maxLines: 1, overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 4),
                Text(product.price,
                  style: const TextStyle(fontSize: 15, fontWeight: FontWeight.bold, color: Colors.cyan),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// ── Banner 组件 ──
class PromoBanner extends StatelessWidget {
  const PromoBanner({super.key});
  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 12),
      padding: const EdgeInsets.all(18),
      decoration: BoxDecoration(
        gradient: const LinearGradient(colors: [Colors.pinkAccent, Colors.deepOrange]),
        borderRadius: BorderRadius.circular(14),
      ),
      child: Row(children: [
        const Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('🔥 夏日焕新季',
              style: TextStyle(fontSize: 17, fontWeight: FontWeight.bold, color: Colors.white),
            ),
            SizedBox(height: 4),
            Text('全场低至 5 折', style: TextStyle(color: Colors.white70, fontSize: 13)),
          ],
        ),
        const Spacer(),
        Container(
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: Colors.white24, borderRadius: BorderRadius.circular(8),
          ),
          child: const Icon(Icons.local_offer, size: 28, color: Colors.white),
        ),
      ]),
    );
  }
}

// ── 首页 ──
class ShopHomePage extends StatefulWidget {
  const ShopHomePage({super.key});
  @override
  State<ShopHomePage> createState() => _ShopHomePageState();
}

class _ShopHomePageState extends State<ShopHomePage> {
  int _currentTab = 0;

  final List<Product> _products = const [
${productsDartList}
  ];

  static const List<String> _categories = ['推荐', '服饰', '鞋靴', '数码', '家居'];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        children: [
          const SizedBox(height: 8),
          // 搜索栏
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12),
            child: TextField(
              decoration: InputDecoration(
                hintText: '搜索潮流好物...',
                prefixIcon: const Icon(Icons.search, color: Colors.cyan),
                filled: true,
                fillColor: Colors.white10,
                contentPadding: const EdgeInsets.symmetric(vertical: 10),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                  borderSide: BorderSide.none,
                ),
              ),
            ),
          ),
          const SizedBox(height: 10),
          // 分类标签
          SizedBox(
            height: 36,
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 12),
              itemCount: _categories.length,
              itemBuilder: (_, i) => Padding(
                padding: const EdgeInsets.only(right: 8),
                child: ChoiceChip(
                  label: Text(_categories[i]),
                  selected: false,
                  selectedColor: Colors.cyan,
                  backgroundColor: Colors.white10,
                  labelStyle: TextStyle(
                    color: false ? Colors.black : Colors.white70,
                    fontSize: 12,
                  ),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                  onSelected: (_) {},
                ),
              ),
            ),
          ),
          const SizedBox(height: 10),
          const PromoBanner(),
          const SizedBox(height: 6),
          // 推荐标题
          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 14),
            child: Align(
              alignment: Alignment.centerLeft,
              child: Text('🔥 为你推荐',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.white),
              ),
            ),
          ),
          const SizedBox(height: 8),
          // 商品网格
          Expanded(
            child: GridView.builder(
              padding: const EdgeInsets.symmetric(horizontal: 12),
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2,
                childAspectRatio: 0.72,
                crossAxisSpacing: 10,
                mainAxisSpacing: 10,
              ),
              itemCount: _products.length,
              itemBuilder: (_, i) => ProductCard(product: _products[i]),
            ),
          ),
        ],
      ),
      bottomNavigationBar: BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        currentIndex: _currentTab,
        selectedItemColor: Colors.cyan,
        unselectedItemColor: Colors.white38,
        backgroundColor: const Color(0xFF0A0F1E),
        onTap: (i) => setState(() => _currentTab = i),
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home), label: '首页'),
          BottomNavigationBarItem(icon: Icon(Icons.explore), label: '发现'),
          BottomNavigationBarItem(icon: Icon(Icons.shopping_bag), label: '购物袋'),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: '我的'),
        ],
      ),
    );
  }
}`;

let streamInterval = null;
let currentIndex = 0;
let generationStarted = false;
let isGenerating = false;
let isPaused = false;
let productsShown = false; // 商品是否已提前展示



// 构建完整手机 HTML
const phoneHtml = ref('');

const buildProductsHtml = (isGenerating) => {
  if (isGenerating) {
    const skeletonCard = (i) => `
      <div class="skeleton-card" style="animation-delay:${i * 0.1}s">
        <div class="skeleton-img"></div>
        <div class="skeleton-info">
          <div class="skeleton-line w80"></div>
          <div class="skeleton-line w50"></div>
        </div>
      </div>`;
    return [0,1,2,3,4,5].map(skeletonCard).join('');
  }
  return shopProducts.map(p => `
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
  `).join('');
};

// ── 动态内容样式（header + body + banner + 商品 + 骨架）──
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
.skeleton-card{background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,0.06);animation:skeletonPulse 1.5s ease-in-out infinite}
.skeleton-img{width:100%;aspect-ratio:4/5;background:linear-gradient(90deg,#e5e5ea 25%,#f0f0f5 50%,#e5e5ea 75%);background-size:200% 100%;animation:shimmer 1.5s infinite}
.skeleton-info{padding:10px}
.skeleton-line{height:12px;border-radius:4px;background:linear-gradient(90deg,#e5e5ea 25%,#f0f0f5 50%,#e5e5ea 75%);background-size:200% 100%;animation:shimmer 1.5s infinite}
.skeleton-line+.skeleton-line{margin-top:8px}
.w80{width:80%}
.w50{width:50%}
@keyframes shimmer{0%{background-position:200% 0}100%{background-position:-200% 0}}
@keyframes skeletonPulse{0%,100%{opacity:1}50%{opacity:0.8}}
<\/style>`;

const updatePhoneHtml = (showSkeleton) => {
  const productsHtml = buildProductsHtml(showSkeleton);
  phoneHtml.value = `${phoneStyles}<div class="phone-header"><div class="header-top"><div class="logo-text">StyleHub</div><div class="header-actions"><i class="fas fa-bell"></i><i class="fas fa-shopping-bag"></i></div></div><div class="search-bar"><i class="fas fa-search"></i><span>搜索潮流好物...</span></div><div class="category-tabs"><span class="cat-item active">推荐</span><span class="cat-item">服饰</span><span class="cat-item">鞋靴</span><span class="cat-item">数码</span><span class="cat-item">家居</span></div></div><div class="phone-body"><div class="hero-banner"><div class="banner-content"><div class="banner-tag">限时特惠</div><h4>夏日焕新季</h4><p>全场低至 5 折 · 满299包邮</p><button class="banner-btn">立即抢购 →</button></div><div class="banner-visual"><img src="https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da?w=140&h=140&fit=crop" alt="" /></div></div><div class="section-header"><span class="section-title">🔥 为你推荐</span><span class="section-more">查看更多 →</span></div><div class="product-grid">${productsHtml}</div></div>`;
};

const startStreaming = (fromIndex) => {
  isGenerating = true;
  isPaused = false;
  productsShown = false;
  updatePhoneHtml(true); // 生成中显示骨架屏
  let index = fromIndex;
  const earlyRevealAt = Math.floor(fullFlutterCode.length * 0.15); // 流到 15% 提前展示商品
  streamInterval = setInterval(() => {
    if (index < fullFlutterCode.length) {
      const chunk = fullFlutterCode.slice(index, index + 12);
      generatedCode.value += chunk;
      currentIndex = index + 12;
      index += 12;
      // 代码流到 15% 时提前渲染商品，不等全部完成
      if (!productsShown && currentIndex >= earlyRevealAt) {
        productsShown = true;
        updatePhoneHtml(false);
      }
    } else {
      clearInterval(streamInterval);
      streamInterval = null;
      isGenerating = false;
      isPaused = false;
      if (!productsShown) updatePhoneHtml(false); // 兜底
    }
  }, 25);
};

const pauseStreaming = () => {
  if (streamInterval) {
    clearInterval(streamInterval);
    streamInterval = null;
    isPaused = true;
  }
};

const resumeStreaming = () => {
  if (isPaused && currentIndex < fullFlutterCode.length) {
    startStreaming(currentIndex);
  }
};

const startGeneration = () => {
  if (generationStarted && !isPaused) return;
  if (!generationStarted) {
    generationStarted = true;
    generatedCode.value = '';
    currentIndex = 0;
    updatePhoneHtml(true);
  }
  if (isPaused) {
    resumeStreaming();
  } else {
    startStreaming(0);
  }
};

let observer = null;
onMounted(() => {
  observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        startGeneration();
      } else {
        pauseStreaming();
      }
    });
  }, { threshold: 0.1 });
  if (screenRef.value) observer.observe(screenRef.value);
  // 初始化预览
  updatePhoneHtml(false);
});
onUnmounted(() => {
  if (streamInterval) clearInterval(streamInterval);
  if (observer) observer.disconnect();
});
</script>

<style scoped>
.screen {
  scroll-snap-align: start;
  min-height: 93vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
}
.content {
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
}
.header-title {
  text-align: center;
  margin-bottom: 2rem;
}
.header-title h2 {
  font-family: 'Orbitron', monospace;
  font-size: 2rem;
}
.sse-layout {
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
}
.code-area {
  flex: 1.2;
  max-height: 720px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}
.preview-area {
  flex: 1;
  max-height: 752px;
  backdrop-filter: blur(16px);
  border-radius: 40px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow: hidden;
}
.code-display {
  background: #01050e;
  border-radius: 40px;
  padding: 1rem;
  margin-top: 1rem;
  max-height: 752px;
  overflow-y: auto;
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  flex: 1;
}

.scroll-hint {
  position: absolute;
  bottom: 1.2rem;
  left: 50%;
  transform: translateX(-50%);
  text-align: center;
  font-size: 0.8rem;
  color: #6f9eff;
  z-index: 1;
  animation: bounce 1.8s infinite;
}
@keyframes bounce {
  0%, 100% { transform: translateX(-50%) translateY(0); }
  50% { transform: translateX(-50%) translateY(8px); }
}
@media (max-width: 968px) {
  .sse-layout { flex-direction: column; }
}
</style>