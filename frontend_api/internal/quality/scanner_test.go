package quality

import (
	"strings"
	"testing"

	"frontend_api/internal/model"
)

const testDSL = `{"globalColors":{"primary":"#1677FF","text":"#333333"},"pageBackground":"#F5F5F5"}`

// 一段能通过全部规则的 Vue 代码（长度 > 100，含 template/script，使用 token 色与 el- 组件）
const cleanVue = `<template>
  <div class="el-card">
    <span class="title">标题</span>
    <p class="desc">这是一段足够长的示例内容，用于保证代码长度超过完整性检查的下限。</p>
  </div>
</template>
<script setup lang="ts">
const title = '示例'
</script>
<style scoped>
.el-card { background: #FFFFFF; color: #333333; }
.title { color: #1677FF; }
</style>`

func TestScanPass(t *testing.T) {
	r := Scan(cleanVue, "vue3", testDSL, "element-plus")
	if r.Status != "pass" {
		t.Fatalf("status = %s, want pass; issues=%v", r.Status, r.Issues)
	}
	if r.Passed != ruleTotal || r.Total != ruleTotal {
		t.Fatalf("passed/total = %d/%d, want %d/%d", r.Passed, r.Total, ruleTotal, ruleTotal)
	}
}

func TestScanExternalCDN(t *testing.T) {
	code := cleanVue + `<link href="https://cdn.jsdelivr.net/npm/x.css">`
	r := Scan(code, "vue3", testDSL, "element-plus")
	if r.Status != "fail" {
		t.Fatalf("status = %s, want fail", r.Status)
	}
	if !hasIssue(r, RuleExternalResource) {
		t.Fatalf("missing %s issue: %v", RuleExternalResource, r.Issues)
	}
}

func TestScanTodoPlaceholder(t *testing.T) {
	code := strings.Replace(cleanVue, "const title = '示例'", "// TODO: 待补充\nconst title = '示例'", 1)
	r := Scan(code, "vue3", testDSL, "element-plus")
	if r.Status != "fail" || !hasIssue(r, RuleTodoPlaceholder) {
		t.Fatalf("status=%s issues=%v, want fail with %s", r.Status, r.Issues, RuleTodoPlaceholder)
	}
}

func TestScanIncompleteCode(t *testing.T) {
	r := Scan("const a = 1", "vue3", testDSL, "")
	if r.Status != "fail" || !hasIssue(r, RuleCompleteness) {
		t.Fatalf("status=%s issues=%v, want fail with %s", r.Status, r.Issues, RuleCompleteness)
	}
}

func TestScanColorOutOfToken(t *testing.T) {
	code := strings.Replace(cleanVue, "color: #1677FF", "color: #FF0000", 1)
	r := Scan(code, "vue3", testDSL, "element-plus")
	if r.Status != "warning" || !hasIssue(r, RuleColorToken) {
		t.Fatalf("status=%s issues=%v, want warning with %s", r.Status, r.Issues, RuleColorToken)
	}
}

func TestScanNeutralColorsAllowed(t *testing.T) {
	code := strings.Replace(cleanVue, "color: #1677FF", "color: #888888", 1)
	r := Scan(code, "vue3", testDSL, "element-plus")
	if hasIssue(r, RuleColorToken) {
		t.Fatalf("neutral gray should pass R4: %v", r.Issues)
	}
}

func TestScanWithoutDSLColorsSkipsR4(t *testing.T) {
	code := strings.Replace(cleanVue, "color: #1677FF", "color: #ABCDEF", 1)
	r := Scan(code, "vue3", "", "element-plus")
	if hasIssue(r, RuleColorToken) {
		t.Fatalf("empty DSL should skip R4: %v", r.Issues)
	}
}

func TestScanComponentLibMissing(t *testing.T) {
	code := strings.ReplaceAll(cleanVue, "el-card", "my-card")
	r := Scan(code, "vue3", testDSL, "element-plus")
	if r.Status != "warning" || !hasIssue(r, RuleComponentLib) {
		t.Fatalf("status=%s issues=%v, want warning with %s", r.Status, r.Issues, RuleComponentLib)
	}
}

// Vant 4 名称含 “ant”，不能被当成 Ant Design 误报（R5 回归）
func TestScanVantComponentLib(t *testing.T) {
	code := strings.ReplaceAll(cleanVue, "el-card", "van-card")
	r := Scan(code, "vue3", testDSL, "Vant 4")
	if hasIssue(r, RuleComponentLib) {
		t.Fatalf("Vant 4 不应误判为 Ant Design: %v", r.Issues)
	}
}

func TestScanVantMissing(t *testing.T) {
	code := strings.ReplaceAll(cleanVue, "el-card", "my-card")
	r := Scan(code, "vue3", testDSL, "Vant 4")
	if !hasIssue(r, RuleComponentLib) {
		t.Fatalf("未使用 van- 组件应报 R5: %v", r.Issues)
	}
}

// React Vant 同样含 “ant”，且代码中无 van- 前缀，应通过（引用 react-vant 即可）
func TestScanReactVantComponentLib(t *testing.T) {
	code := cleanVue + `\nimport { Button } from 'react-vant';`
	r := Scan(code, "react", testDSL, "React Vant")
	if hasIssue(r, RuleComponentLib) {
		t.Fatalf("React Vant 不应误报: %v", r.Issues)
	}
}

// Flutter Element / Material UI (MUI) 不应被 element / material 分支误伤
func TestScanNoFalsePositiveOnOtherLibs(t *testing.T) {
	flutterCode := `import 'package:flutter/material.dart';
class Page extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return const Container(child: Text('示例文本，保证代码长度超过完整性下限'));
  }
}`
	if hasIssue(Scan(flutterCode, "flutter", testDSL, "Flutter Element"), RuleComponentLib) {
		t.Fatal("Flutter Element 不应被 Element Plus 分支误判")
	}
	reactCode := cleanVue + `\nimport { Button } from '@mui/material';`
	if hasIssue(Scan(reactCode, "react", testDSL, "Material UI (MUI)"), RuleComponentLib) {
		t.Fatal("Material UI (MUI) 不应被 Material.dart 分支误判")
	}
}

func TestScanFlutterColorLiteral(t *testing.T) {
	code := `import 'package:flutter/material.dart';
class Page extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(color: Color(0xFF00FF00), child: const Text('示例'));
  }
}`
	r := Scan(code, "flutter", testDSL, "material")
	if r.Status != "warning" || !hasIssue(r, RuleColorToken) {
		t.Fatalf("status=%s issues=%v, want warning with %s", r.Status, r.Issues, RuleColorToken)
	}
}

func TestNormalizeHex(t *testing.T) {
	cases := map[string]string{"#abc": "#AABBCC", "#1677ff": "#1677FF", " #FFF ": "#FFFFFF"}
	for in, want := range cases {
		if got := normalizeHex(in); got != want {
			t.Fatalf("normalizeHex(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsNeutral(t *testing.T) {
	for _, c := range []string{"#FFFFFF", "#000000", "#333333", "#F5F5F5"} {
		if !isNeutral(c) {
			t.Fatalf("%s should be neutral", c)
		}
	}
	if isNeutral("#1677FF") {
		t.Fatal("#1677FF should not be neutral")
	}
}

func hasIssue(r model.QualityReport, rule string) bool {
	for _, is := range r.Issues {
		if is.Rule == rule {
			return true
		}
	}
	return false
}
