package controllers

import (
	"frontend_api/utils"

	"github.com/gin-gonic/gin"
)

// GenerateController AI 代码生成控制器
type GenerateController struct{}

// DesignItem 设计稿项
type DesignItem struct {
	Image       string `json:"image" binding:"required"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// GenerateRequest AI 代码生成请求
type GenerateRequest struct {
	Designs   []DesignItem `json:"designs" binding:"required,min=1"`
	Framework string       `json:"framework" binding:"required"`
	Quality   int          `json:"quality" binding:"required,min=1,max=100"`
}

// ScoreDimension 评分维度
type ScoreDimension struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Icon  string `json:"icon"`
}

// GenerateResponse AI 代码生成响应
type GenerateResponse struct {
	Code       string           `json:"code"`
	Score      int              `json:"score"`
	Dimensions []ScoreDimension `json:"dimensions"`
}

// Generate 处理代码生成请求
func (g *GenerateController) Generate(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	// TODO: 接入 AI 服务 (OpenAI / Claude / 自定义模型)
	// aiService := services.NewAIService()
	// result, err := aiService.GenerateCode(req.Designs, req.Framework, req.Quality)
	// if err != nil {
	//     utils.InternalError(c, "AI 代码生成失败: "+err.Error())
	//     return
	// }

	resp := GenerateResponse{
		Code:  getMockCode(req.Framework),
		Score: getMockScore(req.Quality),
		Dimensions: []ScoreDimension{
			{Name: "视觉还原度", Score: calcDimScore(req.Quality, 0), Icon: "fas fa-palette"},
			{Name: "代码质量", Score: calcDimScore(req.Quality, -3), Icon: "fas fa-code"},
			{Name: "响应式设计", Score: calcDimScore(req.Quality, -5), Icon: "fas fa-mobile-alt"},
			{Name: "性能优化", Score: calcDimScore(req.Quality, -2), Icon: "fas fa-tachometer-alt"},
		},
	}

	utils.Success(c, resp, "代码生成成功")
}

// getMockCode 生成示例代码（实际应替换为 AI 返回结果）
func getMockCode(framework string) string {
	switch framework {
	case "flutter":
		return `import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'DesignCode AI',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
        useMaterial3: true,
      ),
      home: const HomePage(),
    );
  }
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('DesignCode AI'),
        centerTitle: true,
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.auto_awesome,
              size: 80,
              color: Theme.of(context).colorScheme.primary,
            ),
            const SizedBox(height: 24),
            Text(
              'AI 生成的 Flutter 代码',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
          ],
        ),
      ),
    );
  }
}`
	case "react":
		return `import React, { useState } from 'react';
import './App.css';

const App: React.FC = () => {
  const [count, setCount] = useState(0);

  return (
    <div className="app">
      <header className="app-header">
        <h1>DesignCode AI</h1>
        <p className="subtitle">AI 生成的 React 代码</p>
      </header>
      <main className="app-main">
        <div className="card">
          <button onClick={() => setCount(count + 1)}>
            点击次数: {count}
          </button>
          <p className="hint">
            编辑 <code>src/App.tsx</code> 并保存以热更新
          </p>
        </div>
      </main>
    </div>
  );
};

export default App;`
	case "vue":
		return `<template>
  <div class="app">
    <header class="app-header">
      <h1>DesignCode AI</h1>
      <p class="subtitle">AI 生成的 Vue 代码</p>
    </header>
    <main class="app-main">
      <div class="card">
        <button @click="count++">
          点击次数: {{ count }}
        </button>
        <p class="hint">
          编辑 <code>App.vue</code> 并保存以热更新
        </p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const count = ref(0);
</script>

<style scoped>
.app { text-align: center; }
.app-header { padding: 40px 0; }
.card { margin: 20px; padding: 30px; border-radius: 12px; border: 1px solid #e0e0e0; }
button { padding: 12px 24px; font-size: 16px; cursor: pointer; border-radius: 8px; }
</style>`
	default:
		return "// 请选择框架: flutter, react, vue"
	}
}

// getMockScore 根据质量参数计算总分
func getMockScore(quality int) int {
	base := 75
	bonus := (quality - 50) * 30 / 50
	score := base + bonus
	if score > 98 {
		return 98
	}
	if score < 60 {
		return 60
	}
	return score
}

// calcDimScore 计算单项维度分数
func calcDimScore(quality int, offset int) int {
	score := getMockScore(quality) + offset
	if score > 98 {
		return 98
	}
	if score < 50 {
		return 50
	}
	return score
}
