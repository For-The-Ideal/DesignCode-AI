package generator

import (
	"context"
	"fmt"
	"frontend_api/pkg/logger"
	"os"
	"path/filepath"
)

// ═══════════════════════════════════════════════
//  GeneratorSkill - Vue3
//  职责：将 DSL 描述转换为 Vue 3 代码
// ═══════════════════════════════════════════════

// Vue3Skill Vue3 代码生成技能
type Vue3Skill struct {
	logger *logger.Logger
}

// NewVue3Skill 创建 Vue3 生成技能
func NewVue3Skill() *Vue3Skill {
	return &Vue3Skill{
		logger: logger.NewLogger("vue3-gen"),
	}
}

// Name 返回技能名称
func (s *Vue3Skill) Name() string {
	return "Vue3GeneratorSkill"
}

// Execute 执行代码生成
// TODO: 接入 AI 模型进行真实代码生成
func (s *Vue3Skill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	gi, ok := input.(Input)
	if !ok {
		return nil, fmt.Errorf("vue3: 输入类型错误")
	}

	s.logger.Infof("[Vue3Skill] 开始生成 Vue3 代码, DSL length=%d", len(gi.DSL))

	prompt := s.loadPrompt()
	_ = prompt // TODO: 将 prompt + DSL 传给 AI 模型

	code := fmt.Sprintf(`<!-- Vue3 代码生成于 DesignCode AI -->
<!-- DSL: %s -->

<template>
  <div class="app">
    <header class="app-header">
      <h1>DesignCode AI</h1>
      <p class="subtitle">AI 生成的 Vue 3 代码</p>
    </header>
    <main class="app-main">
      <div class="card">
        <p>Vue3 代码已生成</p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const count = ref(0)
</script>

<style scoped>
.app { text-align: center; padding: 20px; }
.app-header { padding: 40px 0; }
.card { padding: 30px; border-radius: 12px; border: 1px solid #e0e0e0; }
</style>`, gi.DSL[:min(len(gi.DSL), 80)])

	preview := `<div class="phone-body"><div class="generated-badge">✨ Vue3 代码已生成</div></div>`

	return Output{Code: code, Preview: preview, Score: 82}, nil
}

func (s *Vue3Skill) loadPrompt() string {
	paths := []string{
		"prompts/vue3/generate.txt",
		"../prompts/vue3/generate.txt",
		filepath.Join("..", "prompts", "vue3", "generate.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
