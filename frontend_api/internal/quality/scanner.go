package quality

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"frontend_api/internal/model"
)

// ═══════════════════════════════════════════════
//  生成代码的确定性质量检查
//
//  与模型自评分不同：全部规则由代码扫描判定，同一份代码结果恒定。
//  规则编号与 prompts/{flutter,vue3,react}/generate.txt 的「硬性规则」一一对应，
//  改动规则时两处必须同步。
// ═══════════════════════════════════════════════

// 规则编号
const (
	RuleExternalResource = "R1" // 禁止外部 CDN / 字体服务
	RuleTodoPlaceholder  = "R2" // 禁止 TODO / 省略占位
	RuleCompleteness     = "R3" // 代码必须完整且具备框架必需结构
	RuleColorToken       = "R4" // 颜色必须取自 DSL 设计 token
	RuleComponentLib     = "R5" // 必须使用任务指定的组件库
)

const ruleTotal = 5

var (
	externalCDNs = []string{
		"cdn.", "unpkg.com", "jsdelivr.net", "cdnjs.cloudflare.com",
		"fonts.googleapis.com", "fonts.gstatic.com",
	}
	todoMarkers = []string{"TODO", "FIXME", "待实现", "// ...", "/* ... */", "省略后续"}
	hexColorRe  = regexp.MustCompile(`#[0-9a-fA-F]{6}\b|#[0-9a-fA-F]{3}\b`)
	// Flutter 颜色字面量：Color(0xFF1677FF)，按 ARGB 取后 6 位
	dartColorRe = regexp.MustCompile(`0x[0-9a-fA-F]{8}`)
)

// Scan 执行全部规则并返回质量报告。
//
//	code         生成器产出的源码
//	framework    flutter | vue3 | react
//	dsl          该页的 DSL（用于提取设计 token 颜色）
//	componentLib 任务指定的组件库（可为空）
func Scan(code, framework, dsl, componentLib string) model.QualityReport {
	report := model.QualityReport{Total: ruleTotal}
	issues := make([]model.QualityIssue, 0, ruleTotal)

	add := func(rule, severity, detail string) {
		issues = append(issues, model.QualityIssue{Rule: rule, Severity: severity, Detail: detail})
	}

	if hit := findExternalResource(code); hit != "" {
		add(RuleExternalResource, "critical", "引用了外部 CDN/字体服务："+hit)
	}
	if hit := findTodoMarker(code); hit != "" {
		add(RuleTodoPlaceholder, "critical", "存在未完成标记："+hit)
	}
	if detail := checkCompleteness(code, framework); detail != "" {
		add(RuleCompleteness, "critical", detail)
	}
	if detail := checkColors(code, dsl); detail != "" {
		add(RuleColorToken, "warning", detail)
	}
	if detail := checkComponentLib(code, componentLib); detail != "" {
		add(RuleComponentLib, "warning", detail)
	}

	report.Passed = ruleTotal - len(issues)
	report.Issues = issues
	switch {
	case hasCritical(issues):
		report.Status = "fail"
	case len(issues) > 0:
		report.Status = "warning"
	default:
		report.Status = "pass"
	}
	return report
}

func hasCritical(issues []model.QualityIssue) bool {
	for _, is := range issues {
		if is.Severity == "critical" {
			return true
		}
	}
	return false
}

func findExternalResource(code string) string {
	lower := strings.ToLower(code)
	for _, d := range externalCDNs {
		if strings.Contains(lower, d) {
			return d
		}
	}
	return ""
}

func findTodoMarker(code string) string {
	for _, m := range todoMarkers {
		if strings.Contains(code, m) {
			return m
		}
	}
	return ""
}

// checkCompleteness 检查代码非空且具备框架必需结构，避免把截断/半成品判为成功
func checkCompleteness(code, framework string) string {
	trimmed := strings.TrimSpace(code)
	if len(trimmed) < 100 {
		return fmt.Sprintf("代码为空或过短（%d 字符）", len(trimmed))
	}
	switch framework {
	case "vue3":
		if !strings.Contains(trimmed, "<template") || !strings.Contains(trimmed, "<script") {
			return "缺少 Vue SFC 必需结构（<template> / <script>）"
		}
	case "react":
		if !strings.Contains(trimmed, "export default") && !strings.Contains(trimmed, "export function") {
			return "缺少 React 组件导出（export default / export function）"
		}
	case "flutter":
		if !strings.Contains(trimmed, "Widget") || !strings.Contains(trimmed, "build(") {
			return "缺少 Flutter Widget 必需结构（Widget / build()）"
		}
	}
	return ""
}

// checkColors 比对代码中出现的十六进制颜色与 DSL 内的设计 token。
// DSL 无颜色信息时跳过（不算违规），黑白灰（R==G==B）允许自由使用。
func checkColors(code, dsl string) string {
	tokens := collectDSLColors(dsl)
	if len(tokens) == 0 {
		return ""
	}
	allowed := make(map[string]bool, len(tokens))
	for _, c := range tokens {
		allowed[normalizeHex(c)] = true
	}

	var bad []string
	seen := make(map[string]bool)
	for _, c := range extractColors(code) {
		n := normalizeHex(c)
		if seen[n] || allowed[n] || isNeutral(n) {
			continue
		}
		seen[n] = true
		bad = append(bad, n)
	}
	if len(bad) == 0 {
		return ""
	}
	sort.Strings(bad)
	if len(bad) > 3 {
		return fmt.Sprintf("有 %d 个颜色不在 DSL 设计 token 内（如 %s）", len(bad), strings.Join(bad[:3], "、"))
	}
	return "颜色不在 DSL 设计 token 内：" + strings.Join(bad, "、")
}

// extractColors 提取代码中的颜色：#RRGGBB / #RGB / Flutter 0xAARRGGBB
func extractColors(code string) []string {
	out := hexColorRe.FindAllString(code, -1)
	for _, m := range dartColorRe.FindAllString(code, -1) {
		out = append(out, "#"+m[4:]) // 0xAARRGGBB → #RRGGBB
	}
	return out
}

// collectDSLColors 递归收集 DSL 中所有 #hex（不依赖 globalColors 的具体结构）
func collectDSLColors(dsl string) []string {
	var v interface{}
	if err := json.Unmarshal([]byte(dsl), &v); err != nil {
		return nil
	}
	out := make([]string, 0, 8)
	var walk func(interface{})
	walk = func(x interface{}) {
		switch t := x.(type) {
		case string:
			out = append(out, extractColors(t)...)
		case map[string]interface{}:
			for _, vv := range t {
				walk(vv)
			}
		case []interface{}:
			for _, vv := range t {
				walk(vv)
			}
		}
	}
	walk(v)
	return out
}

// normalizeHex 统一为大写 #RRGGBB
func normalizeHex(c string) string {
	c = strings.ToUpper(strings.TrimSpace(c))
	if len(c) == 4 {
		return string([]byte{'#', c[1], c[1], c[2], c[2], c[3], c[3]})
	}
	return c
}

// isNeutral 黑白灰（R==G==B）
func isNeutral(hex string) bool {
	if len(hex) != 7 {
		return false
	}
	return hex[1] == hex[3] && hex[3] == hex[5] && hex[2] == hex[4] && hex[4] == hex[6]
}

// checkComponentLib 检查任务指定的组件库是否真的被使用。
// 判定顺序有意义：组件库名之间可能互相包含（如 “Vant 4” 含 “ant”、
// “Material UI (MUI)” 含 “material”），必须先剔除子串误命中再按库检查。
func checkComponentLib(code, componentLib string) string {
	lib := strings.ToLower(strings.TrimSpace(componentLib))
	if lib == "" {
		return ""
	}
	lower := strings.ToLower(code)
	switch {
	case strings.Contains(lib, "vant"):
		// Vant / React Vant：名称含 “ant”，必须在 Ant Design 分支之前判定
		if !strings.Contains(lower, "van-") && !strings.Contains(lower, "vant") {
			return "指定 Vant，但代码中未使用（无 van- 组件/引用）"
		}
	case strings.Contains(lib, "element") && !strings.Contains(lib, "flutter"):
		// Flutter Element 没有 el- 前缀，不参与该检查
		if !strings.Contains(lower, "el-") && !strings.Contains(lower, "element-plus") {
			return "指定 Element Plus，但代码中未使用（无 el- 组件/引用）"
		}
	case strings.Contains(lib, "ant"):
		if !strings.Contains(lower, "ant-") && !strings.Contains(lower, "antd") {
			return "指定 Ant Design，但代码中未使用（无 ant- 组件/引用）"
		}
	case strings.Contains(lib, "material") && !strings.Contains(lib, "ui"):
		// Flutter Material Design；“Material UI (MUI)” 是 React 库，不查 material.dart
		if !strings.Contains(lower, "material.dart") {
			return "指定 Material，但代码未引用 material.dart"
		}
	}
	return ""
}
