package generator

import (
	"strings"
	"testing"
)

func TestJoinPrompt(t *testing.T) {
	cases := []struct {
		name      string
		base      string
		framework string
		want      string
	}{
		{"both", "BASE", "FW", "BASE\n\nFW"},
		{"base only", "BASE", "", "BASE"},
		{"framework only", "", "FW", "FW"},
		{"none", "", "", ""},
		{"trim", "  BASE \n", "\n FW  ", "BASE\n\nFW"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := joinPrompt(c.base, c.framework); got != c.want {
				t.Fatalf("joinPrompt(%q, %q) = %q, want %q", c.base, c.framework, got, c.want)
			}
		})
	}
}

// 依赖真实文件：go test 的 cwd 是包目录，readPromptFile 会向上回退到 frontend_api/prompts
func TestLoadPromptRealFiles(t *testing.T) {
	for _, fw := range []string{"vue3", "react", "flutter"} {
		t.Run(fw, func(t *testing.T) {
			p := LoadPrompt(fw)
			if p == "" {
				t.Fatalf("LoadPrompt(%s) = empty", fw)
			}
			if !strings.Contains(p, "R1") {
				t.Fatalf("missing base rules: %s", p[:min(len(p), 80)])
			}
			if !strings.Contains(p, "## 框架：") {
				t.Fatalf("missing framework section")
			}
			// 公共段必须在框架段之前
			if strings.Index(p, "## 输出契约") > strings.Index(p, "## 框架：") {
				t.Fatalf("base section should come before framework section")
			}
		})
	}
}

func TestLoadPromptUnknownFrameworkFallsBackToBase(t *testing.T) {
	p := LoadPrompt("unknown")
	if p == "" || !strings.Contains(p, "R1") {
		t.Fatalf("unknown framework should still load base, got %q", p[:min(len(p), 80)])
	}
	if strings.Contains(p, "## 框架：") {
		t.Fatalf("unknown framework should not contain framework section")
	}
}
