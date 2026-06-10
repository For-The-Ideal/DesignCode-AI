import { ref, nextTick } from 'vue'

/**
 * 卡片扩散动画 composable
 * 用法：卡片从中心点扩散到周围预设位置
 *
 * @param {Object} options
 * @param {Ref<HTMLElement>} options.containerRef  - 容器 DOM 引用
 * @param {number}        options.cardWidth       - 卡片宽度 px
 * @param {number}        options.cardHeight      - 卡片高度 px
 * @param {number}        options.centerX         - 中心点 x（相对于容器）
 * @param {number}        options.centerY         - 中心点 y（相对于容器）
 * @param {Array}         options.anglePoints     - [{ angle, rFactor }] 角度与半径倍率
 * @param {number}        options.radius          - 基础半径
 */
export function useCardSpreadAnimation(options = {}) {
  const {
    cardWidth = 140,
    cardHeight = 70,
    anglePoints = [],
    radius = 280,
  } = options

  const animationDone = ref(false)
  const cards = ref([])

  /** 初始化卡片：全部归位到中心点 */
  const initCards = (configs, cx, cy) => {
    const shuffled = [...configs].sort(() => Math.random() - 0.5)
    const list = []
    for (let i = 0; i < shuffled.length; i++) {
      const cfg = shuffled[i]
      const pt = anglePoints[i] || { angle: 0, rFactor: 1 }
      const rad = (pt.angle * Math.PI) / 180
      const r = radius * pt.rFactor

      list.push({
        ...cfg,
        currentTop: `${cy - cardHeight / 2}px`,
        currentLeft: `${cx - cardWidth / 2}px`,
        targetTop: `${cy + Math.sin(rad) * r - cardHeight / 2}px`,
        targetLeft: `${cx + Math.cos(rad) * r - cardWidth / 2}px`,
        duration: `${3 + Math.random() * 2}s`,
        delay: `${i * 0.15}s`,
        isAnimating: false,
        isVisible: false,
      })
    }
    cards.value = list
  }

  /** 从中心向外扩散 */
  const spreadBubbles = () => {
    cards.value.forEach(c => { c.isVisible = true })

    setTimeout(() => {
      cards.value.forEach(c => {
        c.isAnimating = true
        c.currentTop = c.targetTop
        c.currentLeft = c.targetLeft
      })
      setTimeout(() => {
        cards.value.forEach(c => { c.isAnimating = false })
        animationDone.value = true
      }, 550)
    }, 50)
  }

  /** 一步到位初始化 + 扩散 */
  const start = async (configs, cx, cy) => {
    initCards(configs, cx, cy)
    await nextTick()
    requestAnimationFrame(() => spreadBubbles())
  }

  return { cards, animationDone, start }
}
