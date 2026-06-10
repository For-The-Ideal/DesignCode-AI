/**
 * 点击产生波纹效果
 * v-wave="{el:xxx, color: #1f1f1f}"
 * el: 波纹效果的元素，没值的话默认为当前元素
 * color: 波纹颜色，没值的话默认为 #1f1f1f
 */
const isServer = import.meta.server
const customWave = {
  createStyle({ width, height }) {
    const style = document.createElement('style');
    // opacity: 0;
    style.innerHTML = `
      @keyframes v-wave-enter {
        to { 
          opacity: 0;
          width: ${width}px;
          height: ${height}px;
          border-width: 1px;
        }
      }
    `;
    return style;
  },
  createWave(confValue = {}) {
    const { color, borderRadius, time } = confValue;
    const wave = document.createElement('div');
    this.setCss(wave, {
      position: 'absolute',
      left: '50%',
      top: '50%',
      transform: 'translate(-50%, -50%)',
      width: '0px',
      height: '0px',
      border: `0px solid ${color}`,
      borderRadius: borderRadius,
      boxSizing: 'border-box',
      zIndex: 1,
      opacity: 1,
      animation: `v-wave-enter ${time}ms cubic-bezier(0.08, 0.82, 0.17, 1) forwards`
    })
    return wave;
  },
  setCss(el, css = {}) {
    const style = el.style;
    for (const name in css) {
      style[name] = css[name];
    }
  },
  updateConfValue(el, value = {}) {
    el._waveConfValue = { eventName: 'mousedown', color: '#1f1f1f', time: 500, offset: 8, borderRadius: '50%', ...(value || {}) };
  },
  waveClickEvent(el) {
    const confValue = el._waveConfValue || {};
    const { offset, time } = confValue;
    return (e) => {
      // console.log('waveClickEvent', el, e, confValue);
      const target = confValue.el ? el.querySelector(confValue.el) : el;
      if (!target) return console.warn('wave 指令配置的元素不存在');
      // console.log('waveClickEvent', target);
      const rect = target.getBoundingClientRect();
      const aminWidth = rect.width + offset;
      const aminHeight = rect.height + offset;
      const style = this.createStyle({ width: aminWidth, height: aminHeight });
      target.appendChild(style)
      const wave = this.createWave(confValue);
      target.appendChild(wave);

      setTimeout(() => {
        target.removeChild(style);
        target.removeChild(wave);
      }, time);
    }
  },
  initWaveEvent(el) {
    const confValue = el._waveConfValue || {};
    if (el._waveClickEvent) {
      el.removeEventListener(confValue.eventName, el._waveClickEvent);
    }
    el._waveClickEvent = this.waveClickEvent(el);

    el.addEventListener(confValue.eventName, el._waveClickEvent);
  },
  init(el, confValue = {}) {
    this.$el = el;
    this.updateConfValue(el, confValue)
    this.initWaveEvent(el)
  },
  mounted(el, binding, vnode) {
    if (isServer) return;
    if (typeof binding.value !== 'object') return console.warn('wave 指令参数值必须为对象');
    customWave.init(el, binding.value)
  },
  updated(el, binding) {
    customWave.init(el, binding.value)
  },
  unmounted(el) {
    const confValue = el._waveConfValue || {};
    el._waveClickEvent && el.removeEventListener(confValue.eventName, el._waveClickEvent);
  }
}

export default customWave
