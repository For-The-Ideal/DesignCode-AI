const isServer = import.meta.server
let isEventTouchend = false; // 是否注册了 Touchend 事件
// v-title 指令（自定义 title 效果）
// data-showtitleel (自定义 指定子元素 title 效果，这个一般用于无法直接在子元素指定v-title）
const customTitle = {
  createSpan() {
    return document.createElement('span');
  },
  setCss(el, css = {}) {
    const style = el.style;
    for (const name in css) {
      style[name] = css[name];
    }
  },
  init() {
    if (isServer || this.$el) return;
    this.$el = this.createSpan();
    this.isTouch = 'ontouchstart' in document;
    this.setCss(this.$el, {
      display: 'none',
      maxWidth: '320px',
      position: 'absolute',
      top: '0px',
      left: '0px',
      right: 'auto',
      border: '1px solid #c3c3c3',
      padding: '4px',
      color: '#555555',
      lineHeight: '16px',
      borderRadius: '3px',
      background: '#ffffff',
      wordBreak: 'break-all',
      fontSize: '12px',
      zIndex: 999999
    });
  },
  getCurEl(el) {
    let bindEl = el;
    const customSubEl = el.getAttribute('data-showtitleel');
    if (el.getAttribute('data-showtitleel')) {
      bindEl = el.querySelector(customSubEl) || el
    }
    return bindEl;
  },
  getPosition(pageX) {
    const width = document.documentElement.clientWidth ||
      document.body.clientWidth ||
      (window.innerWidth - 17);
    let left, right;
    if (pageX >= (width - 80)) {
      left = 'auto';
      right = `${width - pageX - 4}px`;
    } else {
      left = `${pageX + 4}px`;
      right = 'auto';
    }
    return { left, right };
  },
  events: {
    handleEvent(event) {
      switch (event.type) {
        case 'mouseenter':
          if (!customTitle.isTouch) {
            this.mouseenter(event);
          }
          break;
        case 'mousemove':
          if (!customTitle.isTouch) {
            this.mousemove(event);
          }
          break;
        case 'mouseleave':
          if (!customTitle.isTouch) {
            this.mouseleave(event);
          }
          break;
        case 'touchstart':
          this.mouseenter(event);
          break;
        case 'touchmove':
          this.mousemove(event);
          break;
        case 'touchend':
          this.mouseleave(event);
          break;
        default:
      }
    },
    mouseenter(event) {
      const el = customTitle.$el;
      const text = event.currentTarget.__title;
      if (!text || !text.length) {
        return;
      }
      el.innerHTML = text;
      !el.parentNode && document.body.appendChild(el);
      const { pageY, pageX } = customTitle.isTouch ? event.touches[0] : event;
      const { left, right } = customTitle.getPosition(pageX);

      customTitle.setCss(el, {
        display: 'block',
        top: `${pageY + 23}px`,
        left, right
      });

      // 解决移动端不能消失的问题
      if (!isEventTouchend && customTitle.isTouch) { // 防止多次注册
        window.addEventListener('touchend', customTitle.events, false);
        isEventTouchend = true;
      }

    },
    mousemove(e) {
      const el = customTitle.$el;
      if (!el.innerText || !el.innerText.length) {
        customTitle.setCss(el, {
          display: 'none'
        });
        return;
      }
      const { pageY, pageX } = customTitle.isTouch ? e.touches[0] : e;
      const { left, right } = customTitle.getPosition(pageX);
      if (e.target.getAttribute('data-type') === 'NoTitle') {
        customTitle.setCss(el, {
          display: 'none'
        });
      } else {
        customTitle.setCss(el, {
          display: 'block',
          top: `${pageY + 24}px`,
          left, right
        });
      }

    },
    mouseleave(event) {
      customTitle.setCss(customTitle.$el, {
        display: 'none'
      });
      if (isEventTouchend && customTitle.isTouch) {
        isEventTouchend = false;
        window.removeEventListener('touchend', customTitle.events, false);
      }
    }
  },
  mounted(el, binding, vnode) {
    if (isServer) return;
    const title = binding.value || el.getAttribute('title') || '';
    if (!title || !title.length) {
      return;
    }
    if (typeof title !== 'string') {
      throw new Error('v-title params not string');
    }
    if (!title || !title.trim()) return;
    customTitle.init();
    // 删除原有的title提示
    el.removeAttribute('title');

    const bindEl = customTitle.getCurEl(el);
    bindEl.__title = title;
    bindEl.addEventListener('mouseenter', customTitle.events, false);
    bindEl.addEventListener('mousemove', customTitle.events, false);
    bindEl.addEventListener('mouseleave', customTitle.events, false);
    bindEl.addEventListener('touchstart', customTitle.events, false);
    bindEl.addEventListener('touchmove', customTitle.events, false);
    bindEl.addEventListener('touchend', customTitle.events, false);
  },
  updated(el, binding) {
    if (binding.oldValue != binding.value) {
      const bindEl = customTitle.getCurEl(el);
      bindEl.__title = binding.value;
      customTitle.$el.innerHTML = bindEl.__title;
    }
  },
  unmounted(el) {
    const bindEl = customTitle.getCurEl(el);
    bindEl.removeEventListener('mouseenter', customTitle.events, false);
    bindEl.removeEventListener('mousemove', customTitle.events, false);
    bindEl.removeEventListener('mouseleave', customTitle.events, false);
    bindEl.removeEventListener('touchstart', customTitle.events, false);
    bindEl.removeEventListener('touchmove', customTitle.events, false);
    bindEl.removeEventListener('touchend', customTitle.events, false);
    const node = customTitle.$el;
    node.parentNode && node.parentNode.removeChild(node);
  }
}

export default customTitle
