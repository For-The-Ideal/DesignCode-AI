<!-- pages/compare.vue -->
<template>
  <div class="page-container">
    <AppHeader />
    
    <main class="main-content">
      <div class="page-header">
        <h1>模型对比</h1>
        <p>对比不同 AI 模型生成的代码质量，选择最佳方案</p>
      </div>

      <!-- 模型卡片网格 -->
      <div class="models-grid">
        <div class="model-card" v-for="model in models" :key="model.name">
          <div class="model-header">
            <span class="model-name">{{ model.name }}</span>
            <span class="model-score">{{ model.score }}分</span>
          </div>
          <div class="model-icon">
            <i :class="model.icon"></i>
          </div>
          <div class="model-stats">
            <div class="stat-item">
              <span>视觉还原</span>
              <span>{{ model.visual }}%</span>
            </div>
            <div class="stat-item">
              <span>代码质量</span>
              <span>{{ model.codeQuality }}%</span>
            </div>
            <div class="stat-item">
              <span>响应式</span>
              <span>{{ model.responsive }}%</span>
            </div>
            <div class="stat-item">
              <span>性能</span>
              <span>{{ model.performance }}%</span>
            </div>
          </div>
          <div class="model-code-preview">
            <pre>{{ model.codePreview }}</pre>
          </div>
          <button class="select-btn">选择此方案</button>
        </div>
      </div>

      <!-- 差异对比区域 -->
      <div class="diff-section">
        <h3 class="diff-title">
          <i class="fas fa-code-branch"></i> 代码差异对比
        </h3>
        <div class="diff-viewer">
          <div class="diff-column">
            <div class="diff-header">GPT-4 Vision</div>
            <div class="diff-code">import 'package:flutter/material.dart';
            
class GPT4Widget extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.all(16),
      child: Text('GPT-4 Generated'),
    );
  }
}</div>
          </div>
          <div class="diff-column">
            <div class="diff-header">最佳方案 (GPT-4 Vision)</div>
            <div class="diff-code diff-added">import 'package:flutter/material.dart';
            
class BestWidget extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.all(16),
      child: Text('Optimized Code'),
    );
  }
}</div>
          </div>
        </div>
      </div>

      <!-- 悬浮按钮 -->
      <button class="floating-btn">
        ✨ 使用最佳方案
      </button>
    </main>
  </div>
</template>

<script setup>
// 模拟数据
const models = [
  {
    name: 'GPT-4 Vision',
    icon: 'fab fa-openai',
    score: 92,
    visual: 94,
    codeQuality: 91,
    responsive: 90,
    performance: 93,
    codePreview: 'import \'package:flutter/material.dart\';\n\nclass GPT4Widget...'
  },
  {
    name: 'Claude 3.5',
    icon: 'fas fa-brain',
    score: 88,
    visual: 89,
    codeQuality: 90,
    responsive: 87,
    performance: 86,
    codePreview: 'import \'package:flutter/material.dart\';\n\nclass ClaudeWidget...'
  },
  {
    name: 'Gemini Pro',
    icon: 'fas fa-star-of-life',
    score: 85,
    visual: 86,
    codeQuality: 84,
    responsive: 85,
    performance: 85,
    codePreview: 'import \'package:flutter/material.dart\';\n\nclass GeminiWidget...'
  }
]
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0a0f 0%, #0f1a1f 100%);
}

.main-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 24px;
}

.page-header {
  text-align: center;
  margin-bottom: 48px;
}

.page-header h1 {
  font-size: 40px;
  font-weight: 800;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 12px;
}

.page-header p {
  color: #888;
  font-size: 16px;
}

.models-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-bottom: 48px;
}

.model-card {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 24px;
  padding: 24px;
  transition: all 0.3s;
}

.model-card:hover {
  border-color: #00ffff;
  transform: translateY(-4px);
  box-shadow: 0 0 30px rgba(0, 255, 255, 0.1);
}

.model-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.model-name {
  font-size: 20px;
  font-weight: 600;
  color: #00ffff;
}

.model-score {
  background: rgba(0, 255, 255, 0.2);
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 14px;
}

.model-icon {
  text-align: center;
  font-size: 48px;
  color: #00ffff;
  margin: 20px 0;
}

.model-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  color: #aaa;
  padding: 6px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.stat-item span:last-child {
  color: #00ffff;
}

.model-code-preview {
  background: rgba(0, 0, 0, 0.4);
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 20px;
}

.model-code-preview pre {
  font-family: monospace;
  font-size: 11px;
  color: #888;
  overflow-x: auto;
  white-space: pre-wrap;
}

.select-btn {
  width: 100%;
  padding: 10px;
  background: rgba(0, 255, 255, 0.1);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 12px;
  color: #00ffff;
  cursor: pointer;
  transition: all 0.3s;
}

.select-btn:hover {
  background: rgba(0, 255, 255, 0.2);
}

.diff-section {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 24px;
  padding: 24px;
}

.diff-title {
  font-size: 18px;
  color: #00ffff;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.diff-viewer {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.diff-column {
  background: rgba(0, 0, 0, 0.3);
  border-radius: 16px;
  overflow: hidden;
}

.diff-header {
  background: rgba(0, 0, 0, 0.5);
  padding: 12px 16px;
  font-size: 14px;
  font-weight: 600;
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
}

.diff-code {
  padding: 16px;
  font-family: monospace;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  color: #ccc;
}

.diff-added {
  background: rgba(0, 255, 0, 0.1);
  border-left: 3px solid #0f0;
}

.floating-btn {
  position: fixed;
  bottom: 30px;
  right: 30px;
  padding: 14px 28px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  border-radius: 40px;
  color: white;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 20px rgba(0, 255, 255, 0.3);
}

.floating-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 255, 255, 0.4);
}

@media (max-width: 968px) {
  .models-grid {
    grid-template-columns: 1fr;
  }
  .diff-viewer {
    grid-template-columns: 1fr;
  }
  .floating-btn {
    bottom: 20px;
    right: 20px;
    padding: 10px 20px;
    font-size: 14px;
  }
}
</style>