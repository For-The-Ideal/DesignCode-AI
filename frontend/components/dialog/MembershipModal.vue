<template>
    <DialogModel ref="dialogRef" width="720px" top="5vh" :showClose="true" customClass="membership-dialog"
        @close="handleClose">

        <!-- 标题 -->
        <template #header>
            <div class="membership-dialog-title">
                <svg class="title-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                    <path d="M12 2L20 7v10l-8 5-8-5V7z" />
                    <path d="M12 7v5l3 2" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
                <span>升级会员</span>
            </div>
        </template>

        <!-- 内容 -->
        <template #default>
            <div class="membership-body">
                <!-- 当前状态 -->
                <div class="current-status" v-if="userInfo?.email">
                    <span class="status-label">当前方案</span>
                    <span class="status-badge" :class="`badge-lv-${userInfo.level}`">
                        {{ planLabel(userInfo.level) }}
                    </span>
                    <span class="status-credits">剩余 <b>{{ userInfo.credits ?? 0 }}</b> 积分</span>
                </div>

                <!-- 套餐卡片（3D 景深轮播） -->
                <div class="plan-cards-stage" ref="stageRef" @wheel="onWheel">
                    <div class="plan-cards">
                        <div v-for="(plan, idx) in sortedPlans" :key="plan.level" class="plan-card" :class="{
                            'plan-active': userInfo.level === plan.level,
                            'plan-recommend': plan.level === 1,
                            'plan-center': idx === 1,
                        }" :style="cardStyle(idx)" @click.stop="switchTo(idx)">
                            <div class="plan-ribbon" v-if="plan.level === 1">推荐</div>

                            <div class="plan-header">
                                <h3 class="plan-name">{{ plan.name }}</h3>
                                <div class="plan-price">
                                    <span class="price-num">{{ plan.price === 0 ? '免费' : '¥' + plan.price }}</span>
                                    <span class="price-unit" v-if="plan.price > 0">/月</span>
                                </div>
                                <div class="plan-original" v-if="plan.original_price > plan.price">
                                    ¥{{ plan.original_price }}
                                </div>
                            </div>

                            <!-- 核心指标 -->
                            <div class="plan-credits">
                                <span class="credits-num">{{ plan.credits_per_month }}</span>
                                <span class="credits-text">次/月</span>
                            </div>

                            <!-- 功能列表 -->
                            <ul class="plan-features">
                                <li v-for="(f, i) in plan.features" :key="i">
                                    <svg class="check-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                                        stroke-width="2">
                                        <path d="M5 13l4 4L19 7" stroke-linecap="round" stroke-linejoin="round" />
                                    </svg>
                                    {{ f }}
                                </li>
                            </ul>

                            <button class="plan-btn"
                                :class="{ 'btn-current': userInfo.level === plan.level }"
                                :disabled="!canUpgrade(plan) || loadingPlan === plan.level"
                                @click.stop="handleUpgrade(plan)">
                                <span v-if="loadingPlan === plan.level">处理中...</span>
                                <span v-else-if="userInfo.level === plan.level">当前方案</span>
                                <span v-else-if="!canUpgrade(plan)">无法降级</span>
                                <span v-else>立即升级</span>
                            </button>
                        </div>
                    </div>
                </div>

                <!-- 积分充值 -->
                <div class="credits-section">
                    <div class="credits-label">
                        <span>购买积分包</span>
                        <span class="credits-desc">一次购买，永久有效</span>
                    </div>
                    <div class="credits-packages">
                        <button v-for="pkg in creditPackages" :key="pkg.id" class="credit-pkg" :disabled="!!loadingPkg"
                            @click="handleBuyCredits(pkg)">
                            <div class="pkg-top">
                                <span class="pkg-credits">{{ pkg.credits }}<small> 次</small></span>
                                <span class="pkg-badge" v-if="loadingPkg === pkg.id">···</span>
                                <span class="pkg-badge" v-else>购买</span>
                            </div>
                            <div class="pkg-bottom">
                                <span class="pkg-price">¥{{ pkg.price }}</span>
                                <s class="pkg-original" v-if="pkg.original_price > pkg.price">¥{{ pkg.original_price
                                    }}</s>
                            </div>
                        </button>
                    </div>
                </div>
            </div>
        </template>

        <template #footer>
            <div class="membership-footer">
                <button class="footer-close-btn" @click="handleClose">关闭</button>
            </div>
        </template>
    </DialogModel>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { storeToRefs } from 'pinia'
import DialogModel from '~/components/dialog/DialogModel.vue'
import { useUserStore } from '~/stores/user'
import { membershipApi } from '~/api/membership'

const userStore = useUserStore()
const { userInfo } = storeToRefs(userStore)

const emit = defineEmits(['close', 'success'])
const dialogRef = ref(null)

const plans = ref([])
const creditPackages = ref([])
const loadingPlan = ref(null)
const loadingPkg = ref(null)

const planLabel = (lv) => ({ 0: '免费版', 1: '专业版', 2: '旗舰版' }[lv] || '免费版')
const canUpgrade = (plan) => plan.level > userInfo.value.level

const fetchData = async () => {
    try {
        const res = await membershipApi.getPlans()
        if (res.code === 200 && res.data) {
            plans.value = res.data.plans || []
            creditPackages.value = res.data.packages || []
            // 当前套餐默认居中
            const cur = userInfo.value.level
            const curIdx = plans.value.findIndex(p => p.level === cur)
            if (curIdx >= 0) {
                const rest = plans.value.map((_, i) => i).filter(i => i !== curIdx)
                displayOrder.value = [rest[0], curIdx, rest[1] || rest[0]]
            }
        }
    } catch {
        // 静默
    }
}

// 显示顺序：左/中/右 对应原始 plans 索引，默认 [free, pro, premium]
const displayOrder = ref([0, 1, 2])

const sortedPlans = computed(() =>
    displayOrder.value.map(i => plans.value[i]).filter(Boolean)
)

// ── 动画工具 ──
const buildTransform = (rotate, scale, z = 0) =>
    `translateY(-50%) translateZ(${z}px) rotateY(${rotate}deg) scale(${scale})`

const cardStyle = (idx) => {
    const n = plans.value.length
    const center = Math.floor(n / 2)
    const offset = idx - center
    const abs = Math.abs(offset)

    const leftPct = n <= 1 ? 50 : 50 + offset * 26

    if (offset === 0) {
        return {
            left: `${leftPct}%`,
            marginLeft: '-110px',
            transform: buildTransform(0, 1.04, 15),
            opacity: 1,
            zIndex: 10,
            filter: 'brightness(1)',
        }
    }
    return {
        left: `${leftPct}%`,
        marginLeft: '-110px',
        transform: buildTransform(offset < 0 ? 40 : -40, 1 - abs * 0.10, -30),
        opacity: 1 - abs * 0.3,
        zIndex: 10 - abs,
        filter: `brightness(${1 - abs * 0.2})`,
    }
}

// ── 交换动画（CSS transition 自驱动淡入淡出）──
const switching = ref(false)

const animateSwap = () => {
    switching.value = true
    setTimeout(() => { switching.value = false }, 600)
}

// ── 点击交换 ──
const switchTo = (displayIdx) => {
    if (switching.value || displayIdx === 1) return
    const next = [...displayOrder.value]
        ;[next[1], next[displayIdx]] = [next[displayIdx], next[1]]
    displayOrder.value = next
    animateSwap()
}

// ── 滚轮循环翻页 ──
const wheelLock = ref(false)

const rotateTo = (dir) => {
    if (switching.value || wheelLock.value) return
    const next = [...displayOrder.value]
    if (dir > 0) {
        // 向下滚 → 所有卡片左移（右侧卡进中间）
        next.push(next.shift())
    } else {
        // 向上滚 → 所有卡片右移（左侧卡进中间）
        next.unshift(next.pop())
    }
    displayOrder.value = next
    animateSwap()
    wheelLock.value = true
    setTimeout(() => { wheelLock.value = false }, 650)
}

const onWheel = (e) => {
    e.preventDefault()
    const dir = e.deltaY > 0 ? 1 : -1
    rotateTo(dir)
}

const handleUpgrade = async (plan) => {
    if (loadingPlan.value) return
    loadingPlan.value = plan.level
    try {
        const res = await membershipApi.upgrade({ level: plan.level, payment_method: 'credits' })
        if (res.code === 200) {
            ElMessage.success(res.message || '升级成功！')
            // 刷新用户信息
            await userStore.setUserInfo(
                {
                    ...userInfo.value,
                    credits: res.data.credits,
                    current_plan: res.data.level,
                    level: res.data.level,
                }
            )
            emit('success')
        } else {
            ElMessage.error(res.message || '升级失败')
        }
    } catch {
        ElMessage.error('网络错误，请稍后再试')
    } finally {
        loadingPlan.value = null
    }
}

const handleBuyCredits = async (pkg) => {
    if (loadingPkg.value) return
    loadingPkg.value = pkg.id
    try {
        const res = await membershipApi.buyCredits({ package_id: pkg.id })
        if (res.code === 200) {
            const d = res.data
            if (d.status === 'pending') {
                ElMessage.info(d.message || '支付接口即将上线')
            } else {
                ElMessage.success('购买成功！')
                await userStore.setUserInfo(
                    {
                        ...userInfo.value,
                        credits: userInfo.value.credits + pkg.credits,
                    }
                )
            }
        } else {
            ElMessage.error(res.message || '购买失败')
        }
    } catch {
        ElMessage.error('网络错误，请稍后再试')
    } finally {
        loadingPkg.value = null
    }
}

const open = () => {
    dialogRef.value?.open()
    fetchData()
}
const handleClose = () => {
    emit('close')
}

defineExpose({ open })
</script>

<style lang="scss" scoped>
/* ═══ Dialog 面板 ═══ */
.membership-body {
    max-height: calc(85vh - 140px);
    overflow-y: auto;
    padding-right: 4px;
}

.membership-dialog-title {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 18px;
    font-weight: 700;
    color: #e2e8f0;
}

.title-icon {
    width: 28px;
    height: 28px;
    color: #00ffff;
    filter: drop-shadow(0 0 6px rgba(0, 255, 255, 0.4));
}

/* ═══ 当前状态 ═══ */
.current-status {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: rgba(0, 255, 255, 0.04);
    border: 1px solid rgba(0, 255, 255, 0.08);
    border-radius: 12px;
    margin-bottom: 20px;
}

.status-label {
    font-size: 13px;
    color: rgba(255, 255, 255, 0.5);
}

.status-badge {
    font-size: 12px;
    font-weight: 600;
    padding: 3px 10px;
    border-radius: 6px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.badge-lv-0 {
    background: rgba(255, 255, 255, 0.06);
    color: #94a3b8;
}

.badge-lv-1 {
    background: rgba(0, 255, 255, 0.12);
    color: #22d3ee;
    border: 1px solid rgba(0, 255, 255, 0.2);
}

.badge-lv-2 {
    background: rgba(255, 0, 255, 0.1);
    color: #f0abfc;
    border: 1px solid rgba(255, 0, 255, 0.2);
}

.status-credits {
    margin-left: auto;
    font-size: 13px;
    color: rgba(255, 255, 255, 0.6);
}

.status-credits b {
    color: #22d3ee;
}

/* ═══ 套餐卡片（3D 景深轮播）═══ */
.plan-cards-stage {
    height: 420px;
    margin-bottom: 24px;
    overflow: visible;
    perspective: 900px;
    perspective-origin: center center;
}

.plan-cards {
    position: relative;
    width: 100%;
    height: 100%;
}

.plan-card {
    position: absolute;
    width: 220px;
    background: rgba(8, 10, 20, 0.92);
    border: 1px solid rgba(255, 255, 255, 0.06);
    border-radius: 16px;
    padding: 20px 16px;
    text-align: center;
    transform-style: preserve-3d;
    backface-visibility: hidden;
    transition: transform 0.55s cubic-bezier(0.23, 1, 0.32, 1),
        opacity 0.55s cubic-bezier(0.23, 1, 0.32, 1),
        filter 0.55s cubic-bezier(0.23, 1, 0.32, 1),
        left 0.55s cubic-bezier(0.23, 1, 0.32, 1),
        box-shadow 0.55s cubic-bezier(0.23, 1, 0.32, 1),
        border-color 0.55s cubic-bezier(0.23, 1, 0.32, 1);
    top: 50%;
    transform-origin: center;
    cursor: pointer;
}

.plan-center {
    border-color: rgba(0, 255, 255, 0.25) !important;
    box-shadow: 0 0 40px rgba(0, 255, 255, 0.1), 0 20px 40px rgba(0, 0, 0, 0.4);
}

.plan-center.plan-recommend {
    border-color: rgba(0, 255, 255, 0.35) !important;
    box-shadow: 0 0 50px rgba(0, 255, 255, 0.15), 0 24px 48px rgba(0, 0, 0, 0.5);
}

/* 推荐标签 */
.plan-ribbon {
    position: absolute;
    top: -1px;
    right: -1px;
    padding: 4px 14px;
    font-size: 11px;
    font-weight: 700;
    color: #000;
    background: linear-gradient(135deg, #00ffff, #22d3ee);
    border-radius: 0 14px 0 12px;
    letter-spacing: 1px;
}

.plan-header {
    margin-bottom: 14px;
}

.plan-name {
    font-size: 15px;
    font-weight: 700;
    color: #e2e8f0;
    margin-bottom: 6px;
}

.plan-price {
    margin-bottom: 2px;
}

.price-num {
    font-size: 28px;
    font-weight: 800;
    color: #fff;
}

.price-unit {
    font-size: 13px;
    color: rgba(255, 255, 255, 0.4);
}

.plan-original {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.25);
    text-decoration: line-through;
}

/* 积分数 */
.plan-credits {
    margin-bottom: 14px;
    padding: 8px 0;
    border-top: 1px solid rgba(255, 255, 255, 0.05);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.credits-num {
    font-size: 22px;
    font-weight: 800;
    color: #22d3ee;
}

.credits-text {
    font-size: 12px;
    color: rgba(0, 255, 255, 0.5);
    margin-left: 4px;
}

/* 功能列表 */
.plan-features {
    list-style: none;
    padding: 0;
    margin: 0 0 16px;
    text-align: left;
}

.plan-features li {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.65);
    padding: 4px 0;
    display: flex;
    align-items: center;
    gap: 6px;
}

.check-icon {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
    color: #22c55e;
}

/* 按钮 */
.plan-btn {
    width: 100%;
    padding: 10px 0;
    border-radius: 10px;
    border: 1px solid rgba(0, 255, 255, 0.2);
    background: linear-gradient(135deg, rgba(0, 255, 255, 0.08), rgba(255, 0, 255, 0.06));
    color: #22d3ee;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
}

.plan-btn:hover:not(:disabled) {
    background: linear-gradient(135deg, rgba(0, 255, 255, 0.15), rgba(255, 0, 255, 0.1));
    box-shadow: 0 0 20px rgba(0, 255, 255, 0.15);
}

.btn-current {
    background: transparent;
    color: rgba(255, 255, 255, 0.3);
    border-color: rgba(255, 255, 255, 0.06);
    cursor: default;
}

.plan-btn:disabled {
    cursor: not-allowed;
    opacity: 0.6;
}

/* ═══ 积分充值 ═══ */
.credits-section {
    background: rgba(8, 10, 20, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 16px;
    padding: 20px;
}

.credits-label {
    display: flex;
    align-items: baseline;
    gap: 12px;
    margin-bottom: 16px;
    font-size: 15px;
    font-weight: 700;
    color: #e2e8f0;
}

.credits-desc {
    font-size: 12px;
    font-weight: 400;
    color: rgba(255, 255, 255, 0.35);
}

.credits-packages {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
}

.credit-pkg {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 16px;
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.07);
    background: rgba(255, 255, 255, 0.02);
    cursor: pointer;
    transition: all 0.3s;
    font-family: inherit;
    text-align: left;
}

.credit-pkg:hover:not(:disabled) {
    border-color: rgba(0, 255, 255, 0.2);
    background: rgba(0, 255, 255, 0.04);
    box-shadow: 0 0 18px rgba(0, 255, 255, 0.06);
}

.credit-pkg:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.pkg-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.pkg-credits {
    font-size: 20px;
    font-weight: 800;
    color: #22d3ee;
}

.pkg-credits small {
    font-size: 12px;
    font-weight: 500;
    color: rgba(0, 255, 255, 0.5);
}

.pkg-badge {
    font-size: 11px;
    font-weight: 600;
    color: rgba(0, 255, 255, 0.7);
    padding: 3px 10px;
    border-radius: 6px;
    border: 1px solid rgba(0, 255, 255, 0.15);
    white-space: nowrap;
}

.pkg-bottom {
    display: flex;
    align-items: baseline;
    gap: 8px;
}

.pkg-price {
    font-size: 15px;
    font-weight: 700;
    color: #e2e8f0;
}

.pkg-original {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.25);
}

.credits-hint {
    margin-top: 14px;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.3);
    text-align: center;
}

/* ═══ 底部 ═══ */
.membership-footer {
    display: flex;
    justify-content: center;
}

.footer-close-btn {
    padding: 8px 32px;
    border-radius: 10px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: transparent;
    color: rgba(255, 255, 255, 0.5);
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s;
}

.footer-close-btn:hover {
    color: #fff;
    border-color: rgba(255, 255, 255, 0.25);
}
</style>
