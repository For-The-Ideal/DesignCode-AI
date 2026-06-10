const isServer = import.meta.server;
let intersectionObserver = null;
let map = null;

const customAnimate = {
  mounted: function (el, binding) {
    const value = binding.value;
    const modifiers = binding.modifiers || {};
    if (!value) return;
    if (isServer) return;
    el.className = 'animated ' + el.className;
    try {
      if (!map) {
        map = new Map();
      }
      if (!intersectionObserver) {
        intersectionObserver = new IntersectionObserver((entries) => {
          for (const entry of entries) {
            const intersecting = entry.isIntersecting || entry.isIntersecting === void 0;
            if (intersecting) {
              const curMap = map.get(entry.target);
              if (!curMap) {
                intersectionObserver.unobserve(entry.target);
                continue;
              }
              const className = entry.target.className || '';
              if (className.indexOf(curMap) === -1) {
                entry.target.className = `${entry.target.className} ${curMap}`;
                intersectionObserver.unobserve(entry.target);
                map.delete(entry.target);
              }
            }
          }
        }, { rootMargin: '0px' });
      }
      intersectionObserver.observe(el);
      map.set(el, value);
    } catch (error) {
      el.className = `${el.className} ${value}`;
      console.log('v-animate', error);
    }
    try {
      // 带有.move 动画完成后移除相关类名，不影响其它动效
      const move = modifiers.hasOwnProperty('.move') && value.move;
      if (move) {
        el.addEventListener('animationend', function handler() {
          el.className = el.className
            .split(' ')
            .filter(e => ['animated', 'opacity-0', ...value.split(' ')].every(f => e.indexOf(f) === -1))
            .join(' ');
          el.removeEventListener('animationend', handler);
        });
      }
    } catch (err) {
      console.log('v-animate.move', err);
    }
  },
  unmounted: function (el, binding) {
    if (!map) return;
    const curMap = map.get(el);
    if (!curMap) return;
    intersectionObserver.unobserve(el);
    map.delete(el);
  }
}

export default customAnimate
