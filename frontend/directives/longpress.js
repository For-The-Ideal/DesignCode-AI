/**
 * 长按指令
 * 使用方法: v-longpress="handler" 或 v-longpress="{ handler, duration: 500 }"
 * 
 * @example
 * // 基础用法
 * <div v-longpress="onLongPress">内容</div>
 * 
 * // 自定义时长
 * <div v-longpress="{ handler: onLongPress, duration: 800 }">内容</div>
 * 
 * // 阻止点击事件（长按后不触发click）
 * <el-image v-longpress="{ handler: onLongPress, preventClick: true }" />
 */

const longpress = {
  mounted(el, binding, vnode) {
    if (typeof binding.value === 'function') {
      // 简单用法: v-longpress="handler"
      const handler = binding.value;
      init(el, { handler, duration: 500, preventClick: false }, vnode);
    } else if (typeof binding.value === 'object') {
      // 对象用法: v-longpress="{ handler, duration, preventClick }"
      const { handler, duration = 500, preventClick = false } = binding.value;
      init(el, { handler, duration, preventClick }, vnode);
    }
  },

  unmounted(el) {
    // 清理定时器
    if (el._longpressTimer) {
      clearTimeout(el._longpressTimer);
      el._longpressTimer = null;
    }
  }
};

function init(el, options, vnode) {
  const { handler, duration, preventClick } = options;

  // 在元素上存储状态
  el._longpressTimer = null;
  el._isLongPressing = false;
  el._shouldPreventClick = false;
  el._longpressHandler = handler;
  el._longpressPreventClick = preventClick;

  // 鼠标事件
  el._handleMouseDown = function (e) {
    el._isLongPressing = true;
    el._shouldPreventClick = false;

    el._longpressTimer = setTimeout(() => {
      if (el._isLongPressing) {
        el._shouldPreventClick = true;
        // 触发长按回调
        const event = new Event('longpress', { bubbles: true });
        el.dispatchEvent(event);

        // 调用传入的处理函数
        if (typeof el._longpressHandler === 'function') {
          el._longpressHandler.call(this, e);
        }
      }
    }, duration);
  };

  el._handleMouseUp = function (e) {
    el._isLongPressing = false;
    if (el._longpressTimer) {
      clearTimeout(el._longpressTimer);
      el._longpressTimer = null;
    }

    // 如果需要阻止点击，延迟恢复
    if (el._longpressPreventClick && el._shouldPreventClick) {
      setTimeout(() => {
        el._shouldPreventClick = false;
      }, 100);
    }
  };

  el._handleMouseLeave = function () {
    el._isLongPressing = false;
    if (el._longpressTimer) {
      clearTimeout(el._longpressTimer);
      el._longpressTimer = null;
    }
    el._shouldPreventClick = false;
  };

  // 触摸事件
  el._handleTouchStart = function (e) {
    el._isLongPressing = true;
    el._shouldPreventClick = false;

    el._longpressTimer = setTimeout(() => {
      if (el._isLongPressing) {
        el._shouldPreventClick = true;
        // 触发长按回调
        const event = new Event('longpress', { bubbles: true });
        el.dispatchEvent(event);

        // 调用传入的处理函数
        if (typeof el._longpressHandler === 'function') {
          el._longpressHandler.call(this, e);
        }
      }
    }, duration);
  };

  el._handleTouchEnd = function (e) {
    el._isLongPressing = false;
    if (el._longpressTimer) {
      clearTimeout(el._longpressTimer);
      el._longpressTimer = null;
    }

    // 如果需要阻止点击，延迟恢复
    if (el._longpressPreventClick && el._shouldPreventClick) {
      setTimeout(() => {
        el._shouldPreventClick = false;
      }, 100);
    }
  };

  el._handleTouchMove = function () {
    el._isLongPressing = false;
    if (el._longpressTimer) {
      clearTimeout(el._longpressTimer);
      el._longpressTimer = null;
    }
    el._shouldPreventClick = false;
  };

  // 阻止点击事件（如果启用）
  el._handleClick = function (e) {
    if (el._longpressPreventClick && el._shouldPreventClick) {
      e.stopPropagation();
      e.preventDefault();
      return false;
    }
  };

  // 绑定事件
  el.addEventListener('mousedown', el._handleMouseDown);
  el.addEventListener('mouseup', el._handleMouseUp);
  el.addEventListener('mouseleave', el._handleMouseLeave);
  el.addEventListener('touchstart', el._handleTouchStart, { passive: true });
  el.addEventListener('touchend', el._handleTouchEnd);
  el.addEventListener('touchmove', el._handleTouchMove, { passive: true });

  if (preventClick) {
    el.addEventListener('click', el._handleClick, true); // 使用捕获阶段
  }
}

export default longpress
