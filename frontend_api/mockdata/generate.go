package mockdata

import "fmt"

// GetMockCode 根据框架生成示例代码
func GetMockCode(framework string) string {
	switch framework {
	case "flutter":
		return `import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp1());
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

// GetMockScore 根据质量参数计算总分
func GetMockScore(quality int) int {
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

// CalcDimScore 计算单项维度分数
func CalcDimScore(quality int, offset int) int {
	score := GetMockScore(quality) + offset
	if score > 98 {
		return 98
	}
	if score < 50 {
		return 50
	}
	return score
}

// GetFrameworkLabel 框架名称中文映射
func GetFrameworkLabel(framework string) string {
	switch framework {
	case "flutter":
		return "Flutter"
	case "react":
		return "React"
	case "vue":
		return "Vue"
	default:
		return framework
	}
}

// BuildPreviewHTML 构建预览 HTML
func BuildPreviewHTML(framework string) string {
	return fmt.Sprintf(`<style>
.phone-header{background:#fff;padding:10px 14px 8px;border-bottom:1px solid #e5e5ea}
.header-top{display:flex;justify-content:space-between;align-items:center;margin-bottom:8px}
.logo-text{font-size:18px;font-weight:800;color:#1c1c1e}
.phone-body{flex:1;overflow-y:auto;padding:16px;background:#f2f2f7}
.generated-badge{background:linear-gradient(135deg,#00ffff,#ff00ff);border-radius:14px;padding:20px;text-align:center;color:#fff;font-size:16px;font-weight:700}
</style><div class="phone-header"><div class="header-top"><div class="logo-text">DesignCode AI</div></div></div><div class="phone-body"><div class="generated-badge">✨ %s 代码已生成</div></div>`,
		GetFrameworkLabel(framework))
}
