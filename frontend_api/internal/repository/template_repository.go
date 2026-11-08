package repository

import (
	"errors"
	"strings"

	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"

	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
//  Template Repository — 预置模板（首页落地页演示用）
// ═══════════════════════════════════════════════

// 默认模板固定 id=1：前端首页第二屏按 /api/v1/template/1 拉取并做打字机流式展示
const defaultTemplateID = 1

const defaultTemplateName = "Luxe 电商首页"
const defaultTemplateFramework = "flutter"
const defaultTemplateDesc = "Flutter 电商首页模板：问候头部 / 搜索栏 / 分类快捷入口 / 活动 Banner / 双列商品瀑布流 / 底部导航"

// TemplateRepository 预置模板存储
type TemplateRepository struct{}

// NewTemplateRepository 创建预置模板存储
func NewTemplateRepository() *TemplateRepository {
	return &TemplateRepository{}
}

// SeedDefault 幂等写入默认模板（id=1）：
//   - 行不存在 → 插入完整默认模板；
//   - 行存在但 template_code/preview_code 为空（历史脏数据）→ 补齐并修正框架/描述；
//   - 已有代码内容 → 完全跳过，不覆盖后台维护过的模板。
func (r *TemplateRepository) SeedDefault() error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}

	var existing model.TaskTemplate
	err := db.Where("id = ?", defaultTemplateID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(&model.TaskTemplate{
			ID:           defaultTemplateID,
			Name:         defaultTemplateName,
			Framework:    defaultTemplateFramework,
			TemplateCode: defaultTemplateCode,
			PreviewCode:  defaultTemplatePreview,
			Description:  defaultTemplateDesc,
		}).Error
	}
	if err != nil {
		return err
	}

	// 已有完整代码内容则视为后台维护过的模板，不覆盖
	if strings.TrimSpace(existing.TemplateCode) != "" && strings.TrimSpace(existing.PreviewCode) != "" {
		return nil
	}

	return db.Model(&model.TaskTemplate{}).Where("id = ?", defaultTemplateID).
		Updates(map[string]interface{}{
			"name":          defaultTemplateName,
			"framework":     defaultTemplateFramework,
			"template_code": defaultTemplateCode,
			"preview_code":  defaultTemplatePreview,
			"description":   defaultTemplateDesc,
		}).Error
}

// defaultTemplateCode 默认模板源码（Flutter 单文件，首页代码区逐字流式展示）
const defaultTemplateCode = `import 'package:flutter/material.dart';

void main() {
  runApp(const LuxeApp());
}

// 品牌主色：活力橙
const Color kBrand = Color(0xFFFF6A3D);

class LuxeApp extends StatelessWidget {
  const LuxeApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'LUXE 精选',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        useMaterial3: true,
        scaffoldBackgroundColor: const Color(0xFFF7F7FA),
        colorScheme: ColorScheme.fromSeed(seedColor: kBrand),
      ),
      home: const HomePage(),
    );
  }
}

// 商品数据模型
class Product {
  final String name;
  final double price;
  final String emoji;
  final String tag;

  const Product(this.name, this.price, this.emoji, {this.tag = ''});
}

// 首页
class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  int _tabIndex = 0;

  static const List<Map<String, String>> _categories = [
    {'icon': '👗', 'label': '女装'},
    {'icon': '👟', 'label': '运动'},
    {'icon': '👜', 'label': '箱包'},
    {'icon': '💄', 'label': '美妆'},
    {'icon': '🎧', 'label': '数码'},
  ];

  static const List<Product> _products = [
    Product('轻盈通勤连衣裙', 299, '👗', tag: '新品'),
    Product('缓震跑步鞋 Air', 459, '👟', tag: '热卖'),
    Product('极简通勤托特包', 389, '👜'),
    Product('降噪蓝牙耳机 Pro', 699, '🎧', tag: '直降'),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 12, 16, 24),
          children: [
            _buildHeader(),
            const SizedBox(height: 16),
            _buildSearchBar(),
            const SizedBox(height: 20),
            _buildCategories(),
            const SizedBox(height: 20),
            _buildBanner(),
            const SizedBox(height: 24),
            _buildSectionTitle(),
            const SizedBox(height: 12),
            _buildProductGrid(),
          ],
        ),
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _tabIndex,
        onTap: (i) => setState(() => _tabIndex = i),
        selectedItemColor: kBrand,
        unselectedItemColor: const Color(0xFF9A9AA5),
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home_rounded), label: '首页'),
          BottomNavigationBarItem(icon: Icon(Icons.explore_rounded), label: '发现'),
          BottomNavigationBarItem(icon: Icon(Icons.shopping_bag_rounded), label: '购物袋'),
          BottomNavigationBarItem(icon: Icon(Icons.person_rounded), label: '我的'),
        ],
      ),
    );
  }

  Widget _buildHeader() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Hi，欢迎回来 👋',
                style: TextStyle(fontSize: 12, color: Colors.grey.shade500)),
            const SizedBox(height: 2),
            const Text(
              'LUXE 精选',
              style: TextStyle(fontSize: 21, fontWeight: FontWeight.w800),
            ),
          ],
        ),
        Container(
          width: 36,
          height: 36,
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(12),
            boxShadow: [
              BoxShadow(color: Colors.black.withOpacity(0.06), blurRadius: 12),
            ],
          ),
          child: const Icon(Icons.notifications_none_rounded, size: 20),
        ),
      ],
    );
  }

  Widget _buildSearchBar() {
    return Container(
      height: 42,
      padding: const EdgeInsets.symmetric(horizontal: 14),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(14),
        boxShadow: [
          BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 14),
        ],
      ),
      child: Row(
        children: [
          const Icon(Icons.search_rounded, size: 18, color: Color(0xFFA7A7B4)),
          const SizedBox(width: 8),
          Text('搜索「连衣裙 / 运动鞋」',
              style: TextStyle(fontSize: 13, color: Colors.grey.shade400)),
        ],
      ),
    );
  }

  Widget _buildCategories() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: _categories.map((c) {
        return Column(
          children: [
            Container(
              width: 48,
              height: 48,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(16),
                boxShadow: [
                  BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 12),
                ],
              ),
              child: Text(c['icon']!, style: const TextStyle(fontSize: 22)),
            ),
            const SizedBox(height: 6),
            Text(c['label']!,
                style: const TextStyle(fontSize: 11, color: Color(0xFF55555F))),
          ],
        );
      }).toList(),
    );
  }

  Widget _buildBanner() {
    return Container(
      padding: const EdgeInsets.all(18),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(18),
        gradient: const LinearGradient(
          colors: [Color(0xFFFF8A5C), Color(0xFFFF5F3D), Color(0xFFFF3D6E)],
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.22),
              borderRadius: BorderRadius.circular(99),
            ),
            child: const Text('618 年中盛典',
                style: TextStyle(fontSize: 10, color: Colors.white, letterSpacing: 1)),
          ),
          const SizedBox(height: 8),
          const Text('爆款限时 5 折起',
              style: TextStyle(
                  fontSize: 19, fontWeight: FontWeight.w800, color: Colors.white)),
          const SizedBox(height: 4),
          const Text('满 300 减 50 · 会员再享 9 折',
              style: TextStyle(fontSize: 11, color: Colors.white70)),
          const SizedBox(height: 12),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 7),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(99),
            ),
            child: const Text('立即抢购 →',
                style: TextStyle(
                    fontSize: 12, fontWeight: FontWeight.w700, color: kBrand)),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionTitle() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        const Text('为你精选',
            style: TextStyle(fontSize: 16, fontWeight: FontWeight.w800)),
        Text('查看全部 ›',
            style: TextStyle(fontSize: 11, color: Colors.grey.shade500)),
      ],
    );
  }

  Widget _buildProductGrid() {
    return GridView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        mainAxisSpacing: 12,
        crossAxisSpacing: 12,
        childAspectRatio: 0.72,
      ),
      itemCount: _products.length,
      itemBuilder: (context, index) {
        final p = _products[index];
        return Container(
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(16),
            boxShadow: [
              BoxShadow(color: Colors.black.withOpacity(0.05), blurRadius: 16),
            ],
          ),
          clipBehavior: Clip.antiAlias,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Container(
                  width: double.infinity,
                  alignment: Alignment.center,
                  color: const Color(0xFFFFF1EA),
                  child: Text(p.emoji, style: const TextStyle(fontSize: 40)),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(10),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(p.name,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: const TextStyle(
                            fontSize: 12, fontWeight: FontWeight.w600)),
                    const SizedBox(height: 5),
                    RichText(
                      text: TextSpan(
                        text: '¥${p.price.toStringAsFixed(0)}',
                        style: const TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.w800,
                            color: kBrand),
                        children: [
                          if (p.tag.isNotEmpty)
                            TextSpan(
                              text: '  ${p.tag}',
                              style: TextStyle(
                                  fontSize: 9, color: Colors.grey.shade500),
                            ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
`

// defaultTemplatePreview 默认模板预览（手机屏内容区 HTML：内嵌样式、cth- 前缀、无外部资源）
const defaultTemplatePreview = `<style>
  .cth-page { flex: 1; min-height: 0; overflow-y: auto; background: #f7f7fa; font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', sans-serif; color: #23232b; }
  .cth-inner { padding: 14px 16px 22px; }
  .cth-head { display: flex; align-items: center; justify-content: space-between; }
  .cth-hi { font-size: 12px; color: #8f8f9d; }
  .cth-brand { margin-top: 2px; font-size: 21px; font-weight: 800; letter-spacing: 0.5px; }
  .cth-bell { width: 36px; height: 36px; border-radius: 12px; display: flex; align-items: center; justify-content: center; background: #fff; box-shadow: 0 4px 12px rgba(30,30,60,.08); font-size: 16px; }
  .cth-search { margin-top: 14px; display: flex; align-items: center; gap: 8px; height: 42px; padding: 0 14px; border-radius: 14px; background: #fff; box-shadow: 0 4px 14px rgba(30,30,60,.06); font-size: 13px; color: #a7a7b4; }
  .cth-cats { margin-top: 18px; display: flex; justify-content: space-between; }
  .cth-cat { display: flex; flex-direction: column; align-items: center; gap: 6px; font-size: 11px; color: #55555f; }
  .cth-cat-ico { width: 48px; height: 48px; border-radius: 16px; display: flex; align-items: center; justify-content: center; font-size: 22px; background: #fff; box-shadow: 0 4px 12px rgba(30,30,60,.06); }
  .cth-banner { margin-top: 20px; border-radius: 18px; padding: 18px; background: linear-gradient(120deg, #ff8a5c 0%, #ff5f3d 55%, #ff3d6e 100%); color: #fff; position: relative; overflow: hidden; }
  .cth-banner::after { content: ''; position: absolute; right: -30px; top: -30px; width: 120px; height: 120px; border-radius: 50%; background: rgba(255,255,255,.18); }
  .cth-banner-tag { display: inline-block; padding: 2px 8px; border-radius: 99px; background: rgba(255,255,255,.22); font-size: 10px; letter-spacing: 1px; }
  .cth-banner-title { margin-top: 8px; font-size: 19px; font-weight: 800; }
  .cth-banner-sub { margin-top: 4px; font-size: 11px; opacity: .85; }
  .cth-banner-btn { margin-top: 12px; display: inline-block; padding: 7px 14px; border-radius: 99px; background: #fff; color: #ff5f3d; font-size: 12px; font-weight: 700; }
  .cth-sec { margin-top: 22px; display: flex; align-items: baseline; justify-content: space-between; }
  .cth-sec-title { font-size: 16px; font-weight: 800; }
  .cth-sec-more { font-size: 11px; color: #8f8f9d; }
  .cth-grid { margin-top: 12px; display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
  .cth-card { background: #fff; border-radius: 16px; overflow: hidden; box-shadow: 0 6px 16px rgba(30,30,60,.06); }
  .cth-thumb { height: 118px; display: flex; align-items: center; justify-content: center; font-size: 40px; background: linear-gradient(150deg, #fff1ea, #ffe3d6); position: relative; }
  .cth-thumb--b { background: linear-gradient(150deg, #eaf1ff, #dbe7ff); }
  .cth-thumb--c { background: linear-gradient(150deg, #f1eaff, #e6dbff); }
  .cth-thumb--d { background: linear-gradient(150deg, #eafff3, #d8ffe9); }
  .cth-flag { position: absolute; left: 8px; top: 8px; padding: 1px 7px; border-radius: 99px; background: #ff5f3d; color: #fff; font-size: 9px; font-weight: 700; }
  .cth-info { padding: 9px 10px 11px; }
  .cth-name { font-size: 12px; font-weight: 600; color: #33333d; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .cth-price { margin-top: 5px; font-size: 14px; font-weight: 800; color: #ff5f3d; }
  .cth-price small { font-size: 10px; font-weight: 600; }
  .cth-sold { margin-left: 6px; font-size: 9px; color: #9a9aa5; font-weight: 400; }
</style>
<div class="cth-page">
  <div class="cth-inner">
    <div class="cth-head">
      <div>
        <div class="cth-hi">Hi，欢迎回来 👋</div>
        <div class="cth-brand">LUXE 精选</div>
      </div>
      <div class="cth-bell">🔔</div>
    </div>
    <div class="cth-search"><span>🔍</span><span>搜索「连衣裙 / 运动鞋」</span></div>
    <div class="cth-cats">
      <div class="cth-cat"><div class="cth-cat-ico">👗</div><span>女装</span></div>
      <div class="cth-cat"><div class="cth-cat-ico">👟</div><span>运动</span></div>
      <div class="cth-cat"><div class="cth-cat-ico">👜</div><span>箱包</span></div>
      <div class="cth-cat"><div class="cth-cat-ico">💄</div><span>美妆</span></div>
      <div class="cth-cat"><div class="cth-cat-ico">🎧</div><span>数码</span></div>
    </div>
    <div class="cth-banner">
      <span class="cth-banner-tag">618 年中盛典</span>
      <div class="cth-banner-title">爆款限时 5 折起</div>
      <div class="cth-banner-sub">满 300 减 50 · 会员再享 9 折</div>
      <span class="cth-banner-btn">立即抢购 →</span>
    </div>
    <div class="cth-sec">
      <span class="cth-sec-title">为你精选</span>
      <span class="cth-sec-more">查看全部 ›</span>
    </div>
    <div class="cth-grid">
      <div class="cth-card">
        <div class="cth-thumb"><span class="cth-flag">新品</span>👗</div>
        <div class="cth-info"><div class="cth-name">轻盈通勤连衣裙</div><div class="cth-price"><small>¥</small>299<span class="cth-sold">已售 1.2k</span></div></div>
      </div>
      <div class="cth-card">
        <div class="cth-thumb cth-thumb--b"><span class="cth-flag">热卖</span>👟</div>
        <div class="cth-info"><div class="cth-name">缓震跑步鞋 Air</div><div class="cth-price"><small>¥</small>459<span class="cth-sold">已售 3.4k</span></div></div>
      </div>
      <div class="cth-card">
        <div class="cth-thumb cth-thumb--c">👜</div>
        <div class="cth-info"><div class="cth-name">极简通勤托特包</div><div class="cth-price"><small>¥</small>389<span class="cth-sold">已售 860</span></div></div>
      </div>
      <div class="cth-card">
        <div class="cth-thumb cth-thumb--d"><span class="cth-flag">直降</span>🎧</div>
        <div class="cth-info"><div class="cth-name">降噪蓝牙耳机 Pro</div><div class="cth-price"><small>¥</small>699<span class="cth-sold">已售 2.1k</span></div></div>
      </div>
    </div>
  </div>
</div>
`
