const isServer = import.meta.server
let injectCount = 0
let svgEl = null

/**
 * v-borderGradient
 * 注入全局 SVG 渐变定义 #borderGradient，供边框光带动画使用
 * 多个组件同时使用时自动引用计数，确保最后一个 unmount 时才移除
 *
 * 用法: <div v-borderGradient>...</div>
 */
const borderGradient = {
  getSSRProps() {
    return {}
  },
  mounted() {
    if (isServer) return
    if (injectCount === 0) {
      svgEl = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
      svgEl.setAttribute('style', 'position:absolute; width:0; height:0')
      svgEl.id = 'global-gradient-def'

      const defs = document.createElementNS('http://www.w3.org/2000/svg', 'defs')
      const grad = document.createElementNS('http://www.w3.org/2000/svg', 'linearGradient')
      grad.setAttribute('id', 'borderGradient')
      grad.setAttribute('x1', '0%')
      grad.setAttribute('y1', '0%')
      grad.setAttribute('x2', '100%')
      grad.setAttribute('y2', '0%')

      const stops = [
        { offset: '0%', color: '#00ffff' },
        { offset: '50%', color: '#ff00ff' },
        { offset: '100%', color: '#00ffff' }
      ]
      stops.forEach(s => {
        const stop = document.createElementNS('http://www.w3.org/2000/svg', 'stop')
        stop.setAttribute('offset', s.offset)
        stop.setAttribute('stop-color', s.color)
        grad.appendChild(stop)
      })

      defs.appendChild(grad)
      svgEl.appendChild(defs)
      document.body.appendChild(svgEl)
    }
    injectCount++
  },
  unmounted() {
    if (isServer) return
    injectCount--
    if (injectCount === 0 && svgEl && svgEl.parentNode) {
      svgEl.parentNode.removeChild(svgEl)
      svgEl = null
    }
  }
}

export default borderGradient
