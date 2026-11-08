package generator

import (
	"strings"
	"testing"
)

func TestCodeLooksComplete(t *testing.T) {
	cases := []struct {
		name   string
		target string
		code   string
		want   bool
	}{
		{"vue3 complete", "vue3", "<template><div /></template>\n<script setup></script>", true},
		{"vue3 truncated template", "vue3", "<template><div>", false},
		{"vue3 no script", "vue3", "<template><div /></template>", false},
		{"vue3 style unclosed", "vue3", "<template></template><script></script><style>.a{", false},
		{"react complete", "react", "export default function A() { return <div>{1}</div> }", true},
		{"react unbalanced", "react", "export default function A() { return <div>", false},
		{"react no export", "react", "function A() { return null }", false},
		{"flutter complete", "flutter", "class A { Widget build(BuildContext c) { return Text('x'); } }", true},
		{"flutter truncated", "flutter", "class A { Widget build(BuildContext c) { return Text(", false},
		{"empty", "vue3", "   ", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := codeLooksComplete(c.target, c.code); got != c.want {
				t.Fatalf("codeLooksComplete(%s) = %v, want %v (code=%q)", c.target, got, c.want, c.code)
			}
		})
	}
}

func TestMergeContinuation(t *testing.T) {
	cases := []struct {
		name string
		base string
		next string
		want string
	}{
		{"no overlap", "a\nb\nc", "d\ne", "a\nb\nc\nd\ne"},
		{"tail overlap", "1\n2\n3\n4\n5", "3\n4\n5\n6", "1\n2\n3\n4\n5\n6"},
		{"rewrite longer", "A\nB\nC\nD", "A\nB\nC\nD\nE", "A\nB\nC\nD\nE"},
		{"rewrite shorter", "A\nB\nC\nD", "A\nB\nC", "A\nB\nC\nD"},
		{"empty piece", "a\nb", "  ", "a\nb"},
		{"duplicate only", "1\n2\n3\n4", "2\n3\n4", "1\n2\n3\n4"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := mergeContinuation(c.base, c.next); got != c.want {
				t.Fatalf("mergeContinuation() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestStripCodeFence(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"both fences", "```html\n<div />\n```", "<div />"},
		{"open fence only", "```html\n<div />", "<div />"},
		{"close fence only", "<div />\n```", "<div />"},
		{"plain", "<div />", "<div />"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := stripCodeFence(c.raw); got != c.want {
				t.Fatalf("stripCodeFence(%q) = %q, want %q", c.raw, got, c.want)
			}
		})
	}
}

// 拼接后的代码必须能通过闭合检查（续写的核心目标）
func TestContinuationMakesCodeComplete(t *testing.T) {
	base := "<template>\n<div class=\"page\">\n<section class=\"head\">\n<h1>邀请好友</h1>\n<p class=\"sub\">分享给朋友</p>"
	next := "<section class=\"head\">\n<h1>邀请好友</h1>\n<p class=\"sub\">分享给朋友</p>\n<van-button>立即邀请</van-button>\n</div>\n</template>\n<script setup>\nconst a = 1\n</script>"
	merged := mergeContinuation(base, next)
	if !codeLooksComplete("vue3", merged) {
		t.Fatalf("merged vue3 code should be complete, got:\n%s", merged)
	}
	if strings.Count(merged, "<h1>邀请好友</h1>") != 1 {
		t.Fatalf("overlap should be deduped, got:\n%s", merged)
	}
}
