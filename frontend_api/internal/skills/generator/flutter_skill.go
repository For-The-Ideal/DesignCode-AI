package generator

import (
	"context"
	"fmt"
	"frontend_api/pkg/logger"
	"os"
	"path/filepath"
)

// ═══════════════════════════════════════════════
//  GeneratorSkill - Flutter
//  职责：将 DSL 描述转换为 Flutter/Dart 代码
// ═══════════════════════════════════════════════

// FlutterSkill Flutter 代码生成技能
type FlutterSkill struct {
	logger *logger.Logger
}

// NewFlutterSkill 创建 Flutter 生成技能
func NewFlutterSkill() *FlutterSkill {
	return &FlutterSkill{
		logger: logger.NewLogger("flutter-gen"),
	}
}

// Name 返回技能名称
func (s *FlutterSkill) Name() string {
	return "FlutterGeneratorSkill"
}

// Execute 执行代码生成
// TODO: 接入 AI 模型进行真实代码生成
func (s *FlutterSkill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	gi, ok := input.(Input)
	if !ok {
		return nil, fmt.Errorf("flutter: 输入类型错误")
	}

	s.logger.Infof("[FlutterSkill] 开始生成 Flutter 代码, DSL length=%d", len(gi.DSL))

	prompt := s.loadPrompt()
	_ = prompt // TODO: 将 prompt + DSL 传给 AI 模型

	// Mock 返回（与旧 mockdata 保持一致）
	code := fmt.Sprintf(`// Flutter 代码生成于 DesignCode AI
// DSL: %s
// 当前为 Mock 实现，后续接入 AI 模型

import 'package:flutter/material.dart';

void main() => runApp(const MyApp());

class MyApp extends StatelessWidget {
  const MyApp({super.key});
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'DesignCode AI',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorSchemeSeed: const Color(0xFF007AFF),
        useMaterial3: true,
      ),
      home: const HomePage(),
    );
  }
}

class HomePage extends StatelessWidget {
  const HomePage({super.key});
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('DesignCode AI')),
      body: const Center(child: Text('Flutter 代码已生成')),
    );
  }
}
`, gi.DSL[:min(len(gi.DSL), 80)])

	preview := `<div class="phone-body"><div class="generated-badge">✨ Flutter 代码已生成</div></div>`

	return Output{Code: code, Preview: preview, Score: 85}, nil
}

func (s *FlutterSkill) loadPrompt() string {
	paths := []string{
		"prompts/flutter/generate.txt",
		"../prompts/flutter/generate.txt",
		filepath.Join("..", "prompts", "flutter", "generate.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
