<!-- pages/history.vue -->
<template>
  <div class="page-container">
    <AppHeader />
    
    <main class="main-content">
      <div class="page-header">
        <div>
          <h1>生成历史</h1>
          <p>查看和管理所有历史生成的代码</p>
        </div>
        <button class="clear-all-btn">清空历史</button>
      </div>

      <!-- 筛选栏 -->
      <div class="filters-bar">
        <input type="text" class="search-input" placeholder="搜索项目...">
        <select class="filter-select">
          <option value="">全部框架</option>
          <option value="flutter">Flutter</option>
          <option value="react">React</option>
          <option value="vue">Vue</option>
        </select>
      </div>

      <!-- 历史列表 -->
      <div class="history-list">
        <div class="history-item" v-for="item in historyData" :key="item.id">
          <div class="history-header">
            <div>
              <h3 class="history-name">{{ item.name }}</h3>
              <div class="history-meta">
                <span><i class="fas fa-code"></i> {{ item.framework }}</span>
                <span><i class="far fa-calendar"></i> {{ item.date }}</span>
                <span class="history-score"><i class="fas fa-star"></i> {{ item.score }}分</span>
              </div>
            </div>
            <div class="history-actions">
              <button class="action-icon view" title="查看代码">
                <i class="far fa-eye"></i>
              </button>
              <button class="action-icon download" title="下载">
                <i class="fas fa-download"></i>
              </button>
              <button class="action-icon delete" title="删除">
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>
          <div class="history-preview">
            <pre>{{ item.codePreview }}</pre>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div class="empty-state" v-if="historyData.length === 0">
        <i class="fas fa-inbox"></i>
        <p>暂无历史记录</p>
        <span>快去生成你的第一个设计稿吧~</span>
      </div>
    </main>
  </div>
</template>

<script setup>
// 模拟历史数据
const historyData = [
  {
    id: 1,
    name: '电商首页设计',
    framework: 'Flutter',
    date: '2024-01-15 14:30',
    score: 92,
    codePreview: 'import \'package:flutter/material.dart\';\n\nclass HomePage extends StatelessWidget {\n  @override\n  Widget build(BuildContext context) {\n    return Scaffold(...'
  },
  {
    id: 2,
    name: '登录页面设计',
    framework: 'React',
    date: '2024-01-14 10:15',
    score: 88,
    codePreview: 'import React from \'react\';\n\nfunction LoginPage() {\n  return (\n    <div className="login-container">\n      ...'
  },
  {
    id: 3,
    name: '个人中心设计',
    framework: 'Vue',
    date: '2024-01-13 09:45',
    score: 85,
    codePreview: '<template>\n  <div class="profile-container">\n    <h1>个人中心</h1>\n    ...'
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 36px;
  font-weight: 800;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 8px;
}

.page-header p {
  color: #888;
}

.clear-all-btn {
  padding: 8px 20px;
  background: rgba(255, 68, 68, 0.1);
  border: 1px solid rgba(255, 68, 68, 0.3);
  border-radius: 8px;
  color: #ff6666;
  cursor: pointer;
  transition: all 0.3s;
}

.clear-all-btn:hover {
  background: rgba(255, 68, 68, 0.2);
}

.filters-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 32px;
}

.search-input {
  flex: 1;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 12px;
  color: white;
  font-size: 14px;
}

.search-input::placeholder {
  color: #666;
}

.filter-select {
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 12px;
  color: white;
  cursor: pointer;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.history-item {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 20px;
  padding: 20px;
  transition: all 0.3s;
}

.history-item:hover {
  border-color: #00ffff;
  transform: translateX(4px);
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.history-name {
  font-size: 18px;
  font-weight: 600;
  color: #00ffff;
  margin-bottom: 8px;
}

.history-meta {
  display: flex;
  gap: 20px;
  font-size: 13px;
  color: #888;
}

.history-meta i {
  margin-right: 4px;
}

.history-score {
  color: #00ffff;
}

.history-actions {
  display: flex;
  gap: 8px;
}

.action-icon {
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  color: #aaa;
  cursor: pointer;
  transition: all 0.3s;
}

.action-icon:hover {
  background: rgba(0, 255, 255, 0.2);
  color: #00ffff;
}

.action-icon.delete:hover {
  background: rgba(255, 68, 68, 0.2);
  color: #ff6666;
}

.history-preview {
  background: rgba(0, 0, 0, 0.3);
  border-radius: 12px;
  padding: 12px;
}

.history-preview pre {
  font-family: monospace;
  font-size: 12px;
  color: #888;
  white-space: pre-wrap;
  line-height: 1.5;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 24px;
}

.empty-state i {
  font-size: 64px;
  color: #444;
  margin-bottom: 20px;
}

.empty-state p {
  font-size: 18px;
  color: #888;
  margin-bottom: 8px;
}

.empty-state span {
  font-size: 14px;
  color: #666;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  .filters-bar {
    flex-direction: column;
  }
  .history-header {
    flex-direction: column;
    gap: 12px;
  }
  .history-actions {
    align-self: flex-end;
  }
}
</style>