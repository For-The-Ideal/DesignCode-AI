package generator

import (
	"context"
	"fmt"
	"frontend_api/pkg/logger"
	"os"
	"path/filepath"
)

// ═══════════════════════════════════════════════
//  GeneratorSkill - React
//  职责：将 DSL 描述转换为 React + TypeScript 代码
// ═══════════════════════════════════════════════

// ReactSkill React 代码生成技能
type ReactSkill struct {
	logger *logger.Logger
}

// NewReactSkill 创建 React 生成技能
func NewReactSkill() *ReactSkill {
	return &ReactSkill{
		logger: logger.NewLogger("react-gen"),
	}
}

// Name 返回技能名称
func (s *ReactSkill) Name() string {
	return "ReactGeneratorSkill"
}

// Execute 执行代码生成
// TODO: 接入 AI 模型进行真实代码生成
func (s *ReactSkill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	gi, ok := input.(Input)
	if !ok {
		return nil, fmt.Errorf("react: 输入类型错误")
	}

	s.logger.Infof("[ReactSkill] 开始生成 React 代码, DSL length=%d", len(gi.DSL))

	prompt := s.loadPrompt()
	_ = prompt // TODO: 将 prompt + DSL 传给 AI 模型

	code := fmt.Sprintf(`// React 代码生成于 DesignCode AI
// DSL: %s

import React from 'react';

function App() {
  return (
    <div className="app">
      <header className="app-header">
        <h1>DesignCode AI</h1>
        <p className="subtitle">AI 生成的 React 代码</p>
      </header>
      <main className="app-main">
        <div className="card">
          <p>React 代码已生成</p>
        </div>
      </main>
    </div>
  );
}

export default App;
`, gi.DSL[:min(len(gi.DSL), 80)])

	preview := `<div class="phone-body"><div class="generated-badge">✨ React 代码已生成</div></div>`

	return Output{Code: code, Preview: preview, Score: 83}, nil
}

func (s *ReactSkill) loadPrompt() string {
	paths := []string{
		"prompts/react/generate.txt",
		"../prompts/react/generate.txt",
		filepath.Join("..", "prompts", "react", "generate.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
