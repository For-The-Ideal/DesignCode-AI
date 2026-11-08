package workflow

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"frontend_api/internal/model"
)

// ═══════════════════════════════════════════════
//  compose_preview：把逐页生成的 HTML 合成单文件 SPA
//  每页包成一个 <section class="app-page" id="page-<pageID>">，
//  通过 [data-nav] 属性 + 一段内联路由脚本实现页面切换，
//  支持浏览器返回（hash + history）。
//  单页任务直接返回该页 HTML，不做任何包装。
// ═══════════════════════════════════════════════

var (
	htmlBodyRe      = regexp.MustCompile(`(?is)<body[^>]*>(.*)</body>`)
	htmlHeadStyleRe = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	htmlHeadLinkRe  = regexp.MustCompile(`(?is)<link[^>]*rel=["']?stylesheet["']?[^>]*>`)
)

// extractPageContent 返回单页可嵌入内容：body 内容 + head 里的 <style>/<link>。
// 生成模型常把样式放在 head，直接取 body 会丢样式；body 内的样式已随内容保留。
// 多页类名冲突是单文件合成的固有限制，Phase 2 用隔离渲染解决。
func extractPageContent(doc string) string {
	body := htmlBodyRe.FindStringSubmatch(doc)
	if len(body) != 2 {
		return doc
	}
	content := body[1]
	// 仅从 body 外的部分收集文档级资源，避免与内容里的样式重复
	head := strings.Replace(doc, content, "", 1)
	var b strings.Builder
	for _, s := range htmlHeadStyleRe.FindAllString(head, -1) {
		b.WriteString(s)
	}
	for _, l := range htmlHeadLinkRe.FindAllString(head, -1) {
		b.WriteString(l)
	}
	b.WriteString(content)
	return b.String()
}

// BuildAppPreview 合成单文件 SPA 预览。
//
//	pages   - 逐页生成结果（入口页在前，由 OrderedPages 保证）
//	entryID - 入口页 page_id，首屏显示该页
func BuildAppPreview(pages []model.PageResult, entryID string) string {
	if len(pages) == 0 {
		return ""
	}
	if len(pages) == 1 {
		return pages[0].Preview
	}
	if entryID == "" {
		entryID = pages[0].PageID
	}

	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html><head><meta charset=\"utf-8\">")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no\">")
	b.WriteString("<style>html,body{margin:0;padding:0;height:100%;}")
	b.WriteString(".app-page{display:none;height:100%;overflow-y:auto;-webkit-overflow-scrolling:touch;}")
	b.WriteString(".app-page.active{display:block;}</style></head><body>")

	for _, p := range pages {
		cls := "app-page"
		if p.PageID == entryID {
			cls += " active"
		}
		route := p.Route
		if route == "" {
			route = "/" + p.PageID
		}
		fmt.Fprintf(&b, "<section class=\"%s\" id=\"page-%s\" data-route=\"%s\">", cls, html.EscapeString(p.PageID), html.EscapeString(route))
		b.WriteString(extractPageContent(p.Preview))
		b.WriteString("</section>")
	}

	b.WriteString(appRouterJS(entryID))
	b.WriteString("</body></html>")
	return b.String()
}

// appRouterJS 页面切换脚本：
//   - 点击任意 [data-nav="<pageID>"] 元素 → 切页并压入历史
//   - [data-nav="back"] → 浏览器返回
//   - 地址 hash 直接定位（刷新/分享保持当前页）
//   - 暴露 window.__APP_GO__ 供 iframe 外部的宿主控制（后续 Phase 2 用）
func appRouterJS(entryID string) string {
	return `<script>(function(){` +
		`function show(id){var hit=false;` +
		`Array.prototype.forEach.call(document.querySelectorAll('.app-page'),function(s){` +
		`var on=s.id==='page-'+id;if(on)hit=true;s.classList.toggle('active',on);});return hit;}` +
		`function go(id){if(!id){return;}if(id==='back'){history.back();return;}` +
		`if(!show(id)){return;}try{history.pushState({page:id},'','#'+id);}catch(e){}}` +
		`document.addEventListener('click',function(e){` +
		`var t=e.target.closest?e.target.closest('[data-nav]'):null;if(!t){return;}` +
		`e.preventDefault();go(t.getAttribute('data-nav'));},true);` +
		`window.addEventListener('popstate',function(e){` +
		`show((e.state&&e.state.page)||location.hash.replace('#','')||'` + entryID + `');});` +
		`show(location.hash.replace('#','')||'` + entryID + `');` +
		`window.__APP_GO__=go;` +
		`})();</script>`
}

// BuildNavMap 把任务提交的页面边表转成结果侧的实际跳转图。
// Route 由目标页补齐，前端可直接用于展示「A → B」及其路由。
func BuildNavMap(links []model.PageLink, pages []model.PageResult) []model.NavEdge {
	if len(links) == 0 {
		return nil
	}
	routeOf := make(map[string]string, len(pages))
	for _, p := range pages {
		routeOf[p.PageID] = p.Route
	}
	out := make([]model.NavEdge, 0, len(links))
	for _, l := range links {
		out = append(out, model.NavEdge{
			FromPageID: l.FromPageID,
			ToPageID:   l.ToPageID,
			Trigger:    l.Trigger,
			NodeKey:    l.NodeKey,
			Route:      routeOf[l.ToPageID],
		})
	}
	return out
}
