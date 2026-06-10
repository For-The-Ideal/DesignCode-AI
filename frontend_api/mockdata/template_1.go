package mockdata

import (
	"log"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════
//  模板1 — 电商商品数据（模拟）
// ═══════════════════════════════════════════════

// Product 商品结构
type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Emoji string `json:"emoji"`
	Img   string `json:"img"`
}

// Template1Data 模板1完整数据
type Template1Data struct {
	TemplateCode string `json:"template_code"`
	PreviewCode  string `json:"preview_code"`
}

var shopProducts = []Product{
	{Name: "Air Max 270", Price: "¥899", Emoji: "👟", Img: "https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=200&h=250&fit=crop"},
	{Name: "简约白T恤", Price: "¥159", Emoji: "👕", Img: "https://images.unsplash.com/photo-1583743814966-8936f5b7be1a?w=200&h=250&fit=crop"},
	{Name: "智能运动手表", Price: "¥1299", Emoji: "⌚", Img: "https://images.unsplash.com/photo-1524592094714-0f0654e20314?w=200&h=250&fit=crop"},
	{Name: "降噪耳机", Price: "¥599", Emoji: "🎧", Img: "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=200&h=250&fit=crop"},
	{Name: "休闲双肩包", Price: "¥299", Emoji: "🎒", Img: "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=200&h=250&fit=crop"},
	{Name: "复古太阳镜", Price: "¥189", Emoji: "🕶️", Img: "https://images.unsplash.com/photo-1572635196237-14b3f281503f?w=200&h=250&fit=crop"},
	{Name: "运动水壶", Price: "¥79", Emoji: "🧴", Img: "https://images.unsplash.com/photo-1548839140-29a749e1cf4d?w=200&h=250&fit=crop"},
	{Name: "蓝牙音箱", Price: "¥399", Emoji: "🔊", Img: "https://images.unsplash.com/photo-1589003077984-894e133dabab?w=200&h=250&fit=crop"},
	{Name: "经典帆布鞋", Price: "¥259", Emoji: "👟", Img: "https://images.unsplash.com/photo-1549298916-b41d501d3772?w=200&h=250&fit=crop"},
	{Name: "潮流棒球帽", Price: "¥129", Emoji: "🧢", Img: "https://images.unsplash.com/photo-1588850561407-ed78c282e36b?w=200&h=250&fit=crop"},
}

func init() {
	var lines string
	for _, p := range shopProducts {
		lines += "      Product(\n        name: '" + p.Name + "',\n        price: '" + p.Price + "',\n        imageUrl: '" + p.Img + "',\n      ),\n"
	}
	productsDartList = lines
	fullFlutterCode = fullFlutterCodeHead + productsDartList + fullFlutterCodeTail
}

var productsDartList string
var fullFlutterCode string

var fullFlutterCodeHead = `import 'package:flutter/material.dart';

// Product Model
class Product {
  final String name;
  final String price;
  final String imageUrl;

  const Product({
    required this.name,
    required this.price,
    required this.imageUrl,
  });
}

class PhonePreview extends StatelessWidget {
  const PhonePreview({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF2F2F7),
      body: Column(
        children: [
          // Phone Header
          _buildPhoneHeader(),
          // Phone Body
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(12.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _buildHeroBanner(),
                  const SizedBox(height: 16),
                  _buildSectionHeader('🔥 为你推荐', '查看更多 →'),
                  const SizedBox(height: 10),
                  _buildProductGrid(),
                ],
              ),
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

  // Phone Header
  Widget _buildPhoneHeader() {
    return Container(
      color: Colors.white,
      padding: const EdgeInsets.fromLTRB(14, 10, 14, 8),
      child: Column(
        children: [
          // Header Top Row
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'StyleHub',
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w800,
                  color: Color(0xFF1C1C1E),
                  letterSpacing: -0.3,
                ),
              ),
              Row(
                children: const [
                  Icon(Icons.notifications_none, size: 16, color: Color(0xFF1C1C1E)),
                  SizedBox(width: 14),
                  Icon(Icons.shopping_bag_outlined, size: 16, color: Color(0xFF1C1C1E)),
                ],
              ),
            ],
          ),
          const SizedBox(height: 8),
          // Search Bar
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
            decoration: BoxDecoration(
              color: const Color(0xFFF2F2F7),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Row(
              children: const [
                Icon(Icons.search, size: 14, color: Color(0xFF8E8E93)),
                SizedBox(width: 8),
                Text('搜索潮流好物...', style: TextStyle(fontSize: 13, color: Color(0xFF8E8E93))),
              ],
            ),
          ),
          const SizedBox(height: 8),
          // Category Tabs
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Row(
              children: const [
                _CategoryItem(label: '推荐', isActive: true),
                _CategoryItem(label: '服饰', isActive: false),
                _CategoryItem(label: '鞋靴', isActive: false),
                _CategoryItem(label: '数码', isActive: false),
                _CategoryItem(label: '家居', isActive: false),
              ],
            ),
          ),
        ],
      ),
    );
  }

  // Hero Banner
  Widget _buildHeroBanner() {
    return Container(
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFFFF3B6E), Color(0xFFFF6B3D)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(14),
      ),
      child: Stack(
        children: [
          // Background decoration circle
          Positioned(
            right: -30,
            top: -30,
            child: Container(
              width: 100,
              height: 100,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: Colors.white.withOpacity(0.12),
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                // Banner Content
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: Colors.white.withOpacity(0.25),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: const Text(
                          '限时特惠',
                          style: TextStyle(
                            fontSize: 10,
                            fontWeight: FontWeight.w700,
                            color: Colors.white,
                          ),
                        ),
                      ),
                      const SizedBox(height: 6),
                      const Text(
                        '夏日焕新季',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w800,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(height: 2),
                      const Text(
                        '全场低至 5 折 · 满299包邮',
                        style: TextStyle(
                          fontSize: 11,
                          color: Colors.white70,
                        ),
                      ),
                      const SizedBox(height: 10),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 5),
                        decoration: BoxDecoration(
                          color: Colors.white,
                          borderRadius: BorderRadius.circular(14),
                        ),
                        child: const Text(
                          '立即抢购 →',
                          style: TextStyle(
                            fontSize: 11,
                            fontWeight: FontWeight.w700,
                            color: Color(0xFFFF3B6E),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 12),
                // Banner Visual
                Container(
                  width: 64,
                  height: 64,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(10),
                    color: Colors.white.withOpacity(0.18),
                  ),
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(10),
                    child: Image.network(
                      'https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da?w=140&h=140&fit=crop',
                      fit: BoxFit.cover,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  // Section Header
  Widget _buildSectionHeader(String title, String moreText) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          title,
          style: const TextStyle(
            fontSize: 15,
            fontWeight: FontWeight.w700,
            color: Color(0xFF1C1C1E),
          ),
        ),
        Text(
          moreText,
          style: const TextStyle(
            fontSize: 12,
            color: Color(0xFF8E8E93),
          ),
        ),
      ],
    );
  }

  // Product Grid
  Widget _buildProductGrid() {
    final products = _getProducts();
    return GridView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        crossAxisSpacing: 8,
        mainAxisSpacing: 8,
        childAspectRatio: 0.75,
      ),
      itemCount: products.length,
      itemBuilder: (context, index) {
        return _buildProductCard(products[index]);
      },
    );
  }

  // Product Card
  Widget _buildProductCard(Product product) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.06),
            blurRadius: 3,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Product Image
          ClipRRect(
            borderRadius: const BorderRadius.only(
              topLeft: Radius.circular(12),
              topRight: Radius.circular(12),
            ),
            child: Image.network(
              product.imageUrl,
              width: double.infinity,
              height: 150,
              fit: BoxFit.cover,
              loadingBuilder: (context, child, loadingProgress) {
                if (loadingProgress == null) return child;
                return Container(
                  height: 150,
                  color: const Color(0xFFF2F2F7),
                  child: const Center(
                    child: CircularProgressIndicator(strokeWidth: 2),
                  ),
                );
              },
            ),
          ),
          // Product Info
          Padding(
            padding: const EdgeInsets.fromLTRB(10, 8, 10, 10),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  product.name,
                  style: const TextStyle(
                    fontSize: 12,
                    color: Color(0xFF1C1C1E),
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 4),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      product.price,
                      style: const TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w800,
                        color: Color(0xFFFF3B6E),
                      ),
                    ),
                    const Text(
                      '已售1.2k',
                      style: TextStyle(
                        fontSize: 10,
                        color: Color(0xFFC7C7CC),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  // Product Data
  List<Product> _getProducts() {
    return const [
`

var fullFlutterCodeTail = `    ];
  }
}

// Category Item Widget
class _CategoryItem extends StatelessWidget {
  final String label;
  final bool isActive;

  const _CategoryItem({
    required this.label,
    required this.isActive,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(right: 4),
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 5),
      decoration: BoxDecoration(
        color: isActive ? const Color(0xFF1C1C1E) : Colors.transparent,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 12,
          color: isActive ? Colors.white : const Color(0xFF8E8E93),
          fontWeight: isActive ? FontWeight.w600 : FontWeight.normal,
        ),
      ),
    );
  }
}
`

var phoneStyles = `<style>
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
</style>`

// buildProductsHtml 在服务端拼接商品卡片 HTML
func buildProductsHtml() string {
	var cards []string
	for _, p := range shopProducts {
		card := `
    <div class="product-card">
      <div class="product-img"><img src="` + p.Img + `" alt="` + p.Name + `" loading="lazy" /></div>
      <div class="product-info">
        <div class="product-name">` + p.Name + `</div>
        <div class="product-price-row">
          <span class="product-price">` + p.Price + `</span>
          <span class="product-sold">已售1.2k</span>
        </div>
      </div>
    </div>`
		cards = append(cards, card)
	}
	return strings.Join(cards, "")
}

// buildPreviewHtml 在服务端一次性拼接完整手机预览 HTML
func buildPreviewHtml() string {
	return phoneStyles + `<div class="phone-header"><div class="header-top"><div class="logo-text">StyleHub</div><div class="header-actions"><i class="fas fa-bell"></i><i class="fas fa-shopping-bag"></i></div></div><div class="search-bar"><i class="fas fa-search"></i><span>搜索潮流好物...</span></div><div class="category-tabs"><span class="cat-item active">推荐</span><span class="cat-item">服饰</span><span class="cat-item">鞋靴</span><span class="cat-item">数码</span><span class="cat-item">家居</span></div></div><div class="phone-body"><div class="hero-banner"><div class="banner-content"><div class="banner-tag">限时特惠</div><h4>夏日焕新季</h4><p>全场低至 5 折 · 满299包邮</p><button class="banner-btn">立即抢购 →</button></div><div class="banner-visual"><img src="https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da?w=140&h=140&fit=crop" alt="" /></div></div><div class="section-header"><span class="section-title">🔥 为你推荐</span><span class="section-more">查看更多 →</span></div><div class="product-grid">` + buildProductsHtml() + `</div></div>`
}

// getTemplate1 返回模板1数据
// 🔌 后期数据库接入：将此处替换为数据库查询
func getTemplate1() *Template1Data {
	return &Template1Data{
		TemplateCode: fullFlutterCode,
		PreviewCode:  buildPreviewHtml(),
	}
}

// ═══════════════════════════════════════════════
//  流式代码生成通道（替代 generate.go 中的 GetMockCode）
// ═══════════════════════════════════════════════

// StreamTemplateCode 返回一个缓冲通道，逐块产出 Flutter 代码段
//
// 每块 5 行代码，块间隔 60ms，模拟 AI 逐字输出效果。
// 全部产出后 close(ch)，调用方可 range 读取。
//
// 用法：
//
//	for chunk := range mockdata.StreamTemplateCode() {
//	    broker.Publish(SSEEvent{Event: "message", Data: chunk})
//	}
func StreamTemplateCode() <-chan string {
	ch := make(chan string, 8)

	go func() {
		defer close(ch)
		lines := strings.Split(fullFlutterCode, "\n")
		chunkSize := 5

		log.Printf("[Template1] 开始流式推送，总行数: %d，每块 %d 行", len(lines), chunkSize)

		for i := 0; i < len(lines); i += chunkSize {
			end := i + chunkSize
			if end > len(lines) {
				end = len(lines)
			}
			chunk := strings.Join(lines[i:end], "\n")
			ch <- chunk
			time.Sleep(60 * time.Millisecond)
		}

		log.Printf("[Template1] 代码推送完成，总块数: %d", (len(lines)+chunkSize-1)/chunkSize)
	}()

	return ch
}

// GetTemplatePreviewHTML 返回模板1的完整预览 HTML
//
// 包含 phoneStyles + header + body + heroBanner + productGrid，
// 替代 generate.go 中的 BuildPreviewHTML()。
func GetTemplatePreviewHTML() string {
	return buildPreviewHtml()
}
