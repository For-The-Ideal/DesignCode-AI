<template>
  <AppSidebar
    brand="代码生成"
    :navItems="menuStore.navItemsWithActive"
    expandedWidth="225px"
    @navClick="handleNavClick">
    <!-- 底部区域：用量 + 升级（仅登录后显示） -->
    <div class="sidebar-bottom" v-if="isLogin">
      <!-- 使用情况（SVG 圆环） -->
      <div class="usage-card">
        <div class="usage-label">
          <span>本月使用情况</span>
        </div>
        <div class="usage-row">
          <div class="usage-ring-wrap">
            <svg class="usage-ring" viewBox="0 0 80 80">
              <defs>
                <linearGradient id="ringGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#00ffff" />
                  <stop offset="100%" stop-color="#ff00ff" />
                </linearGradient>
              </defs>
              <circle cx="40" cy="40" r="34" class="ring-bg" />
              <circle
                cx="40" cy="40" r="34"
                class="ring-fill"
                :style="ringStyle"
              />
            </svg>
            <div class="ring-center">
              <span class="ring-value">{{ usagePercent }}%</span>
            </div>
          </div>
          <div class="usage-info">
            <span class="usage-remaining">
              剩余 <span class="usage-remaining-num">{{ credits }}</span> 次
            </span>
            <p class="usage-total">共 {{ usageTotal }} 次</p>
          </div>
        </div>
      </div>

      <!-- 升级会员 -->
      <div class="upgrade-card">
        <div class="upgrade-glow"></div>
        <div class="upgrade-title">💎 升级会员</div>
        <div class="upgrade-desc">解锁更多高级功能</div>
        <button class="upgrade-btn">立即升级</button>
      </div>
    </div>
  </AppSidebar>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useUserStore } from '~/stores/user'
import { useCommonStore } from '~/stores/common'
import { useMenuListStore } from '~/stores/menuList'
import AppSidebar from '~/components/layout/AppSidebar.vue'

const userStore = useUserStore()
const commonStore = useCommonStore()
const { isLogin, userInfo } = storeToRefs(userStore)
const credits = computed(() => userInfo.value.credits ?? 0)
const openLoginModal = () => commonStore.setLoginModalVisible(true)
const router = useRouter()
const menuStore = useMenuListStore()

const handleNavClick = (item) => {
  // 1. 已激活 → 不重复跳转
  if (item.active) return
  // 2. 未登录 → 弹出登录弹窗
  if (!isLogin.value) {
    openLoginModal()
    return
  }
  router.push(item.to)
}


const usageTotal = 100

const usagePercent = computed(() => {
  if (usageTotal <= 0) return 0
  const used = Math.max(0, usageTotal - credits.value)
  return Math.round((used / usageTotal) * 100)
})

const ringStyle = computed(() => {
  const r = 34
  const circumference = 2 * Math.PI * r
  const offset = circumference * (1 - usagePercent.value / 100)
  return {
    strokeDasharray: `${circumference}`,
    strokeDashoffset: `${offset}`,
  }
})

</script>

<style scoped>
/* 底部 */
.sidebar-bottom {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* ═══ 使用量卡片（SVG 圆环） ═══ */
.usage-card {
  background: rgba(0, 255, 255, 0.03);
  border-radius: 14px;
  padding: 16px;
  border: 1px solid rgba(0, 255, 255, 0.08);
}
.usage-label {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  margin-bottom: 12px;
}

.usage-ring-wrap {
  position: relative;
  flex-shrink: 0;
}
.usage-row {
  display: flex;
  align-items: center;
  gap: 2px;
}
.usage-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  margin-left: 12px;
}
.usage-ring {
  width: 80px;
  height: 80px;
  transform: rotate(-90deg);
}
.ring-bg {
  fill: none;
  stroke: rgba(255, 255, 255, 0.05);
  stroke-width: 4;
}
.ring-fill {
  fill: none;
  stroke: url(#ringGradient);
  stroke-width: 4;
  stroke-linecap: round;
  transition: stroke-dashoffset 0.6s ease;
}
.ring-center {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.ring-value {
  font-size: 18px;
  font-weight: 700;
  color: #00ffff;
}

.usage-remaining {
  font-size: 12px;
  color: #fff;
  text-align: left;
}
.usage-remaining-num {
  font-weight: 500;
  margin: 0px 3px;
}


.usage-total {
  font-size: 12px;
  color: gray;
}

/* ═══ 升级卡片（主题色版） ═══ */
.upgrade-card {
  position: relative;
  border-radius: 14px;
  padding: 18px 16px;
  text-align: center;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.04), rgba(255, 0, 255, 0.04));
  border: 1px solid rgba(0, 255, 255, 0.1);
}
.upgrade-glow {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(
    circle at 50% 50%,
    rgba(0, 255, 255, 0.04) 0%,
    transparent 60%
  );
  pointer-events: none;
}
.upgrade-title {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  position: relative;
}
.upgrade-desc {
  font-size: 12px;
  color: #fff;
  margin-top: 2px;
  position: relative;
}
.upgrade-btn {
  margin-top: 10px;
  width: 100%;
  padding: 8px 0;
  border-radius: 8px;
  border: none;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 2px 16px rgba(0, 255, 255, 0.15);
  transition: all 0.3s;
  position: relative;
  background: linear-gradient(90deg, #60a5fa, #818cf8);
}
.upgrade-btn:hover {
  box-shadow: 0 4px 28px rgba(0, 255, 255, 0.35);
  transform: translateY(-1px);
}
</style>