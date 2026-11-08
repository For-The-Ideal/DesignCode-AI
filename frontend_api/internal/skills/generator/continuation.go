package generator

import (
	"context"
	"strings"

	"frontend_api/pkg/ai"
)

// ═══════════════════════════════════════════════
//  截断续写：AI 输出被 max_tokens 截断时，把已生成部分回传，
//  要求模型从断点接着写，最后拼接去重，尽量还原成一份完整代码。
// ═══════════════════════════════════════════════

// maxContinuationRounds 单次生成最多续写轮数：
// 首轮 + 2 轮续写 = 最多 3 次请求，避免模型持续截断时无限消耗
const maxContinuationRounds = 2

// continuationInstruction 续写指令：只补剩余部分，不重复、不解释、不带 markdown 围栏
const continuationInstruction = "你上一次的输出因长度限制被截断。请从截断处继续补全剩余代码，要求：\n" +
	"1. 紧接着已输出的最后一个字符往下写，不要重写或重复已输出的内容；\n" +
	"2. 不要输出任何解释、说明，也不要用 ``` 代码围栏包裹；\n" +
	"3. 保持相同的缩进与代码风格，直到完整闭合全部结构（标签 / 括号 / 组件段）。"

// generateWithContinuation 调用生成模型；输出被截断时自动续写并拼接。
//
// 截断判定（任一命中即续写）：
//   - finish_reason=length（被 max_tokens 截断）
//   - 代码未通过框架闭合检查（SFC 缺结束标签 / 括号不配平 / 缺导出结构）
//
// 续写失败不丢弃已生成内容，保留现状交给质量检查兜底。
func generateWithContinuation(
	ctx context.Context,
	client ai.Client,
	target string,
	messages []ai.Message,
	logf func(format string, args ...interface{}),
) (string, error) {
	if logf == nil {
		logf = func(string, ...interface{}) {}
	}
	opts := []ai.Option{ai.WithDeepSeekThinking("high"), ai.WithMaxTokens(16384)}

	res, err := client.ChatWithMeta(ctx, messages, opts...)
	if err != nil {
		return "", err
	}
	code := stripCodeFence(res.Content)
	incomplete := res.Truncated() || !codeLooksComplete(target, code)

	for round := 1; round <= maxContinuationRounds && incomplete; round++ {
		logf("[续写] %s 输出不完整（finish_reason=%q, code=%d chars），第 %d/%d 轮续写",
			target, res.FinishReason, len(code), round, maxContinuationRounds)

		// 已输出内容作为 assistant 消息回传，让模型"接着说"；续写不再启用 thinking，聚焦补全
		contMessages := make([]ai.Message, 0, len(messages)+2)
		contMessages = append(contMessages, messages...)
		contMessages = append(contMessages,
			ai.Message{Role: "assistant", Content: res.Content},
			ai.Message{Role: "user", Content: continuationInstruction},
		)

		next, nextErr := client.ChatWithMeta(ctx, contMessages, ai.WithMaxTokens(16384))
		if nextErr != nil {
			logf("[续写] %s 第 %d 轮失败: %v，保留已生成 %d chars", target, round, nextErr, len(code))
			break
		}
		piece := stripCodeFence(next.Content)
		if piece == "" {
			logf("[续写] %s 第 %d 轮返回空内容，停止续写", target, round)
			break
		}
		code = mergeContinuation(code, piece)
		res = next
		incomplete = next.Truncated() || !codeLooksComplete(target, code)
		if !incomplete {
			logf("[续写] %s 第 %d 轮补全成功，最终 %d chars", target, round, len(code))
		}
	}
	if incomplete {
		// 续写次数用尽仍不完整：保留拼接结果，交给质量检查 R3 标记，日志供排查
		logf("[续写] %s 续写 %d 轮后仍不完整（code=%d chars），保留当前结果", target, maxContinuationRounds, len(code))
	}
	return code, nil
}

// codeLooksComplete 框架闭合检查：判断代码是否已完整结束
//
//	vue3    - 出现开标签就必须有对应闭标签，且必须具备 <template> / <script>
//	react   - 括号配平且存在组件导出
//	flutter - 括号配平且存在 Widget 构建结构
func codeLooksComplete(target, code string) bool {
	if strings.TrimSpace(code) == "" {
		return false
	}
	switch target {
	case "vue3":
		pairs := [][2]string{
			{"<template", "</template>"},
			{"<script", "</script>"},
			{"<style", "</style>"},
		}
		for _, p := range pairs {
			if strings.Contains(code, p[0]) && !strings.Contains(code, p[1]) {
				return false
			}
		}
		return strings.Contains(code, "<template") && strings.Contains(code, "<script")
	case "react":
		return hasBalancedBrackets(code) &&
			(strings.Contains(code, "export default") || strings.Contains(code, "export function"))
	case "flutter":
		return hasBalancedBrackets(code) && strings.Contains(code, "build(")
	default:
		return hasBalancedBrackets(code)
	}
}

// hasBalancedBrackets 括号配平检查（启发式：不跳过字符串/注释中的括号，
// 误差只会导致多一次续写，不影响最终正确性）
func hasBalancedBrackets(code string) bool {
	var braces, parens, brackets int
	for _, r := range code {
		switch r {
		case '{':
			braces++
		case '}':
			braces--
		case '(':
			parens++
		case ')':
			parens--
		case '[':
			brackets++
		case ']':
			brackets--
		}
	}
	return braces == 0 && parens == 0 && brackets == 0
}

// stripCodeFence 去掉代码块标记。
// 与 workflow.cleanCodeBlock 不同：续写片段只有开头或只有结尾带 ``` 时也要处理。
func stripCodeFence(raw string) string {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "```") {
		if i := strings.Index(s, "\n"); i != -1 {
			s = s[i+1:]
		} else {
			s = strings.TrimPrefix(s, "```")
		}
	}
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "```") {
		s = strings.TrimSuffix(s, "```")
	}
	return strings.TrimSpace(s)
}

// mergeContinuation 把续写片段拼到已生成代码后面，处理两类重叠：
//
//  1. 模型重写了整份代码（开头与已生成代码一致）→ 取更完整的一份；
//  2. 模型接着断点写但重复了尾部若干行 → 按最大重叠行数去重后拼接。
func mergeContinuation(base, next string) string {
	next = strings.TrimSpace(next)
	if next == "" {
		return base
	}
	baseLines := strings.Split(base, "\n")
	nextLines := strings.Split(next, "\n")

	if len(baseLines) >= 3 && len(nextLines) >= 3 && sameLines(baseLines[:3], nextLines[:3]) {
		if len(next) >= len(base) {
			return next
		}
		return base
	}

	max := 40
	if len(baseLines) < max {
		max = len(baseLines)
	}
	if len(nextLines) < max {
		max = len(nextLines)
	}
	for n := max; n >= 3; n-- {
		if sameLines(baseLines[len(baseLines)-n:], nextLines[:n]) {
			nextLines = nextLines[n:]
			break
		}
	}

	rest := strings.Join(nextLines, "\n")
	if strings.TrimSpace(rest) == "" {
		return base
	}
	if strings.HasSuffix(base, "\n") {
		return base + rest
	}
	return base + "\n" + rest
}

// sameLines 逐行比较（忽略行首行尾空白，容忍续写缩进漂移）
func sameLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if strings.TrimSpace(a[i]) != strings.TrimSpace(b[i]) {
			return false
		}
	}
	return true
}
