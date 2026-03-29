pages/index.vue
<template>
  <div class="page-container">
    <AppHeader />

    <main class="main-content">

      <HeroSection />
      <!-- 主内容网格 -->
      <div class="main-grid">
        <!-- 左侧：上传区 -->
        <div class="upload-section">
          <!-- 拖拽上传卡片 -->
          <div class="tech-card">
            <div class="card-title">
              <i class="fas fa-cloud-upload-alt"></i>
              <span>上传设计稿</span>
            </div>
            <div class="dropzone">
              <div class="dropzone-icon">
                <i class="fas fa-file-image"></i>
              </div>
              <p>点击或拖拽上传设计稿</p>
              <span>支持 PNG、JPG、JPEG，最多 20 张</span>
            </div>

            <!-- 预览区域（模拟数据） -->
            <div class="preview-area">
              <div class="preview-header">
                <span><i class="fas fa-images"></i> 已上传 (3)</span>
                <button class="clear-btn">清空全部</button>
              </div>
              <div class="preview-grid">
                <div class="preview-item">
                  <img
                    src="https://picsum.photos/200/200?random=1"
                    alt="预览"
                  />
                  <div class="preview-remove">×</div>
                </div>
                <div class="preview-item">
                  <img
                    src="https://picsum.photos/200/200?random=2"
                    alt="预览"
                  />
                  <div class="preview-remove">×</div>
                </div>
                <div class="preview-item">
                  <img
                    src="https://picsum.photos/200/200?random=3"
                    alt="预览"
                  />
                  <div class="preview-remove">×</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 配置卡片 -->
          <div class="tech-card config-card">
            <div class="card-title">
              <i class="fas fa-sliders-h"></i>
              <span>生成配置</span>
            </div>

            <div class="framework-selector">
              <div class="framework-btn active">
                <i class="fab fa-flutter"></i> Flutter
              </div>
              <div class="framework-btn">
                <i class="fab fa-react"></i> React
              </div>
              <div class="framework-btn"><i class="fab fa-vuejs"></i> Vue</div>
            </div>

            <div class="slider-container">
              <div class="slider-label">
                <span>质量要求</span>
                <span class="quality-value">90</span>
              </div>
              <input
                type="range"
                value="90"
                min="60"
                max="100"
                step="5"
                class="quality-slider"
              />
            </div>
          </div>

          <!-- 生成按钮 -->
          <button class="generate-btn">
            <i class="fas fa-play"></i> 开始生成代码
          </button>

          <div class="score-card">
            <div class="score-header">
              <div>
                <h3>质量评分</h3>
                <p class="score-sub">AI 多维评估</p>
              </div>
              <div class="score-value">85</div>
            </div>
            <div class="score-dimensions">
              <div class="dimension-item">
                <div class="dimension-label">
                  <span><i class="fas fa-palette"></i> 视觉还原度</span>
                  <span>85</span>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 85%"></div>
                </div>
              </div>
              <div class="dimension-item">
                <div class="dimension-label">
                  <span><i class="fas fa-code"></i> 代码质量</span>
                  <span>82</span>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 82%"></div>
                </div>
              </div>
              <div class="dimension-item">
                <div class="dimension-label">
                  <span><i class="fas fa-mobile-alt"></i> 响应式设计</span>
                  <span>78</span>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 78%"></div>
                </div>
              </div>
              <div class="dimension-item">
                <div class="dimension-label">
                  <span><i class="fas fa-tachometer-alt"></i> 性能优化</span>
                  <span>75</span>
                </div>
                <div class="progress-bar">
                  <div class="progress-fill" style="width: 75%"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：结果区 -->
        <div class="result-section">
          <CodeEditor
            v-model="code"
            language="dart"
            :readonly="true"
            :font-size="14"
            :word-wrap="true"
            :line-numbers="true"
            theme="vs-dark"
            height="600px"
            @change="handleCodeChange"
            @copy="handleCopy"
            @format="handleFormat"
          />
          <OptimizationPanel />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import CodeEditor from "~/components/code/CodeEditor.vue";
import OptimizationPanel from "~/components/code/OptimizationPanel.vue";
import HeroSection from "~/components/layout/HeroSection.vue";

// 代码数据
const code = ref(
  `import 'package:flutter/material.dart';

class CustomBottomNavBarComponents extends StatelessWidget {
  final int currentIndex;
  final ValueChanged<int> onTabChanged; // 现在这个是必须的

  final List<Map<String, dynamic>> _routeNames = [
    {
      'name': '发现',
      'icon': Icons.home,
    },
    {
      'name': '资讯',
      'icon': Icons.explore,
    },
    {
      'name': '学院',
      'icon': Icons.school,
    },
    {
      'name': '个人中心',
      'icon': Icons.person,
    },
  ];

  CustomBottomNavBarComponents({super.key, required this.currentIndex, required this.onTabChanged});

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 70, // 比默认更高
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: Colors.black12,
            blurRadius: 10,
            offset: Offset(0, -2),
          ),
        ],
        borderRadius: BorderRadius.only(
          topLeft: Radius.circular(20),
          topRight: Radius.circular(20),
        ),
      ),
      child: Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children:List.generate(_routeNames.length, (index){
        return _buildNavItem(index, _routeNames[index]);
      }) 
    ),
    );
  }

  Widget _buildNavItem(int index,Map<String, dynamic> item) {

    bool isSelected = currentIndex == index;

    return GestureDetector(
      onTap: () {
        if(index == currentIndex) return;
        onTabChanged(index); // 调用回调函数
      },
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(item['icon'],
            color: isSelected ? Colors.black : Colors.grey,
            size: 24,
          ),
          SizedBox(height: 4),
          Text(item['name'],style: TextStyle(
              color: isSelected ? Colors.black : Colors.grey,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }
}
`,
);

// 事件处理
const handleCodeChange = (newCode) => {
  console.log("代码已更新:", newCode);
};

const handleCopy = (code) => {
  console.log("代码已复制");
};

const handleFormat = (formattedCode) => {
  console.log("代码已格式化");
};

// 通过 ref 调用组件方法
const editorRef = ref();

const formatCode = () => {
  editorRef.value?.format();
};

const copyCode = () => {
  editorRef.value?.copy();
};

const getCode = () => {
  return editorRef.value?.getValue();
};
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0a0f 0%, #0f1a1f 100%);
}

.main-content {
  max-width: 90%;
  margin: 0 auto;
  padding: 0 24px;
}

.hero-section {
  text-align: center;
  padding: 60px 0 40px;
  position: relative;
  overflow: hidden;
}

/* 背景光晕效果 */
.hero-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80%;
  height: 80%;
  background: radial-gradient(circle, rgba(0, 255, 255, 0.1), transparent);
  filter: blur(60px);
  pointer-events: none;
  animation: glowPulse 3s ease-in-out infinite;
}

@keyframes glowPulse {
  0%, 100% {
    opacity: 0.3;
    transform: translate(-50%, -50%) scale(0.8);
  }
  50% {
    opacity: 0.8;
    transform: translate(-50%, -50%) scale(1.2);
  }
}

/* 徽章样式 */
.badge {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 8px 20px;
  background: rgba(0, 255, 255, 0.08);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 60px;
  font-size: 12px;
  color: #00ffff;
  margin-bottom: 28px;
  backdrop-filter: blur(8px);
  position: relative;
  z-index: 1;
}

/* 闪烁的圆点 */
.blinking-dot {
  font-size: 8px;
  animation: blink 1.2s ease-in-out infinite;
}

@keyframes blink {
  0%, 100% {
    opacity: 0.3;
    text-shadow: 0 0 0px #00ffff;
  }
  50% {
    opacity: 1;
    text-shadow: 0 0 8px #00ffff;
  }
}

/* 脉冲边框效果 */
.pulse-badge {
  animation: borderPulse 2s ease-in-out infinite;
}

@keyframes borderPulse {
  0%, 100% {
    border-color: rgba(0, 255, 255, 0.3);
    box-shadow: 0 0 0px rgba(0, 255, 255, 0);
  }
  50% {
    border-color: rgba(0, 255, 255, 0.8);
    box-shadow: 0 0 20px rgba(0, 255, 255, 0.3);
  }
}

/* 标题 - 渐变流动效果 */
.hero-title {
  font-size: 56px;
  font-weight: 800;
  background: linear-gradient(
    135deg,
    #00ffff 0%,
    #ff00ff 25%,
    #00ffff 50%,
    #ff00ff 75%,
    #00ffff 100%
  );
  background-size: 300% auto;
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 20px;
  position: relative;
  z-index: 1;
  letter-spacing: -0.02em;
}

.gradient-flow {
  animation: gradientFlow 3s linear infinite;
}

@keyframes gradientFlow {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 300% 50%;
  }
}

/* 文字发光效果 */
.glow-text {
  font-size: 18px;
  color: #ccc;
  max-width: 600px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
  animation: textGlow 2s ease-in-out infinite;
}

@keyframes textGlow {
  0%, 100% {
    text-shadow: 0 0 0px rgba(0, 255, 255, 0);
    color: #ccc;
  }
  50% {
    text-shadow: 0 0 20px rgba(0, 255, 255, 0.5);
    color: #fff;
  }
}

/* 添加粒子漂浮效果 */
.hero-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: radial-gradient(circle at 20% 40%, rgba(0, 255, 255, 0.1) 1px, transparent 1px);
  background-size: 40px 40px;
  pointer-events: none;
  animation: particleFloat 20s linear infinite;
}

@keyframes particleFloat {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 100px 100px;
  }
}

/* 响应式 */
@media (max-width: 768px) {
  .hero-title {
    font-size: 36px;
  }
  
  .hero-desc {
    font-size: 14px;
    padding: 0 20px;
  }
  
  .badge {
    font-size: 10px;
    padding: 6px 16px;
    gap: 8px;
  }
}

/* 主网格 */
.main-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
  margin-top: 20px;
}

/* 科技卡片 */
.tech-card {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 24px;
  padding: 24px;
  margin-bottom: 24px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #00ffff;
}

/* 拖拽上传区 */
.dropzone {
  border: 2px dashed rgba(0, 255, 255, 0.4);
  border-radius: 20px;
  padding: 40px;
  text-align: center;
  background: rgba(0, 0, 0, 0.3);
  cursor: pointer;
}

.dropzone-icon {
  width: 70px;
  height: 70px;
  background: linear-gradient(
    135deg,
    rgba(0, 255, 255, 0.2),
    rgba(255, 0, 255, 0.2)
  );
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.dropzone-icon i {
  font-size: 28px;
  color: #00ffff;
}

.dropzone p {
  margin-bottom: 8px;
}

.dropzone span {
  font-size: 12px;
  color: #666;
}

/* 预览区域 */
.preview-area {
  margin-top: 20px;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 13px;
  color: #00ffff;
}

.clear-btn {
  background: none;
  border: none;
  color: #ff4444;
  cursor: pointer;
  font-size: 12px;
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.preview-item {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  aspect-ratio: 1;
}

.preview-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-remove {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 22px;
  height: 22px;
  background: rgba(255, 0, 0, 0.8);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.3s;
}

.preview-item:hover .preview-remove {
  opacity: 1;
}

/* 框架选择器 */
.framework-selector {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.framework-btn {
  padding: 12px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 12px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.framework-btn.active {
  background: linear-gradient(
    135deg,
    rgba(0, 255, 255, 0.2),
    rgba(255, 0, 255, 0.2)
  );
  border-color: #00ffff;
  color: #00ffff;
}

/* 滑块 */
.slider-container {
  margin-top: 8px;
}

.slider-label {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 14px;
  color: #aaa;
}

.quality-value {
  color: #00ffff;
}

.quality-slider {
  width: 100%;
  height: 4px;
  -webkit-appearance: none;
  background: #333;
  border-radius: 2px;
}

.quality-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 16px;
  height: 16px;
  background: #00ffff;
  border-radius: 50%;
  cursor: pointer;
}

/* 生成按钮 */
.generate-btn {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  border-radius: 16px;
  font-size: 16px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.3s;
}

.generate-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(0, 255, 255, 0.3);
}

/* 代码编辑器 */
.code-editor {
  background: #0a0a0f;
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 24px;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.5);
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
}

.editor-dots {
  display: flex;
  gap: 8px;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.dot-red {
  background: #ff5f56;
}
.dot-yellow {
  background: #ffbd2e;
}
.dot-green {
  background: #27c93f;
}

.editor-lang {
  font-family: monospace;
  font-size: 12px;
  color: #00ffff;
}

.editor-copy {
  padding: 4px 12px;
  background: rgba(0, 255, 255, 0.1);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
}

.code-area {
  width: 100%;
  height: 320px;
  background: transparent;
  border: none;
  padding: 20px;
  color: #ccc;
  font-family: monospace;
  font-size: 13px;
  line-height: 1.6;
  resize: none;
  outline: none;
}

/* 评分卡片 */
.score-card {
  background: rgba(15, 20, 30, 0.6);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 24px;
  overflow: hidden;
  margin: 24px 0px;
}

.score-header {
  background: linear-gradient(
    135deg,
    rgba(0, 255, 255, 0.2),
    rgba(255, 0, 255, 0.2)
  );
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.score-header h3 {
  font-size: 18px;
}

.score-sub {
  font-size: 12px;
  opacity: 0.7;
  margin-top: 4px;
}

.score-value {
  font-size: 48px;
  font-weight: 700;
  color: #00ffff;
}

.score-dimensions {
  padding: 20px;
}

.dimension-item {
  margin-bottom: 16px;
}

.dimension-label {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  margin-bottom: 8px;
  color: #aaa;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: #333;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #00ffff, #ff00ff);
  border-radius: 3px;
}

.suggestions {
  padding: 20px;
  border-top: 1px solid rgba(0, 255, 255, 0.2);
}

.suggestions h4 {
  font-size: 14px;
  margin-bottom: 12px;
  color: #00ffff;
}

.suggestions ul {
  list-style: none;
}

.suggestions li {
  font-size: 13px;
  color: #aaa;
  margin-bottom: 8px;
  padding-left: 20px;
  position: relative;
}

.suggestions li::before {
  content: "→";
  position: absolute;
  left: 0;
  color: #00ffff;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 12px;
}

.action-btn {
  flex: 1;
  padding: 12px;
  background: rgba(0, 255, 255, 0.1);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 12px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
}

.action-btn:hover {
  background: rgba(0, 255, 255, 0.2);
}

.optimize-btn {
  background: linear-gradient(
    135deg,
    rgba(255, 0, 255, 0.2),
    rgba(0, 255, 255, 0.2)
  );
  border-color: #ff00ff;
}

@media (max-width: 968px) {
  .main-grid {
    grid-template-columns: 1fr;
  }
  .hero-title {
    font-size: 32px;
  }
  .preview-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
</style>
