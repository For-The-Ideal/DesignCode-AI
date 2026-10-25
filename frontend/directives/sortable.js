/**
 * 拖拽排序指令（事件委托版）
 *
 * 挂在 v-for 的「容器」元素上，通过事件委托协调兄弟项的拖拽与落点，
 * 命中落点后对绑定的数组做原地重排（方案 A：指令直接改 list）。
 *
 * 用法:
 *   <ul v-sortable="{ list: items, disabled: items.length <= 1 }">
 *     <li v-for="(it, i) in items" :key="it.id" data-sortable-item :data-idx="i">...</li>
 *   </ul>
 *
 * @param {Array}    list      需要原地重排的响应式数组（必填）
 * @param {Boolean}  disabled  为 true 时禁用拖拽（单项 / 只读场景），代码保留
 * @param {String}   handle    可选，拖拽把手选择器；设置后仅按住把手才可拖动
 *
 * 说明：每个可拖动子项需带 data-sortable-item 与 data-idx（其在 list 中的下标）。
 * 拖拽 / 悬停的视觉通过 .sortable-dragging、.sortable-over 两个 class 呈现，
 * 由使用方在自身样式中定义外观。
 */

const DRAGGING_CLS = 'sortable-dragging'
const OVER_CLS = 'sortable-over'

const sortable = {
  mounted(el, binding) {
    el._sortableValue = binding.value
    const state = { from: -1, over: -1 }

    const cfg = () => el._sortableValue || {}
    const getList = () => (Array.isArray(cfg().list) ? cfg().list : null)
    const isDisabled = () => !!cfg().disabled
    const getHandle = () => cfg().handle || null

    // 找到触发事件所属的可拖子项
    const findItem = (target) =>
      target && target.closest ? target.closest('[data-sortable-item]') : null
    const idxOf = (item) => parseInt(item.getAttribute('data-idx'), 10)

    const clearOver = () => {
      el.querySelectorAll('.' + OVER_CLS).forEach((n) => n.classList.remove(OVER_CLS))
      state.over = -1
    }
    const cleanup = () => {
      const dragging = el.querySelector('.' + DRAGGING_CLS)
      if (dragging) dragging.classList.remove(DRAGGING_CLS)
      clearOver()
      state.from = -1
    }

    const onDragStart = (e) => {
      if (isDisabled()) return
      const item = findItem(e.target)
      if (!item) return
      // 指定了 handle 时，仅允许从把手处发起拖动
      const handle = getHandle()
      if (handle && !(e.target.closest && e.target.closest(handle))) return

      state.from = idxOf(item)
      item.classList.add(DRAGGING_CLS)
      e.dataTransfer.effectAllowed = 'move'
      try {
        e.dataTransfer.setData('text/plain', String(state.from))
      } catch (_) {
        /* 部分浏览器 setData 可选 */
      }
    }

    const onDragOver = (e) => {
      if (isDisabled() || state.from < 0) return
      const item = findItem(e.target)
      if (!item) return
      e.preventDefault()
      const idx = idxOf(item)
      if (state.over !== idx) {
        clearOver()
        item.classList.add(OVER_CLS)
        state.over = idx
      }
    }

    const onDrop = (e) => {
      if (isDisabled() || state.from < 0) return
      e.preventDefault()
      const item = findItem(e.target)
      if (!item) {
        cleanup()
        return
      }
      const to = idxOf(item)
      const list = getList()
      if (list && to >= 0 && state.from !== to && to < list.length) {
        const [moved] = list.splice(state.from, 1)
        list.splice(to, 0, moved)
      }
      cleanup()
    }

    const onDragEnd = () => cleanup()

    el._sortableHandlers = { onDragStart, onDragOver, onDrop, onDragEnd }
    el.addEventListener('dragstart', onDragStart)
    el.addEventListener('dragover', onDragOver)
    el.addEventListener('drop', onDrop)
    el.addEventListener('dragend', onDragEnd)

    applyDraggable(el)
  },

  updated(el, binding) {
    el._sortableValue = binding.value
    applyDraggable(el)
  },

  unmounted(el) {
    const h = el._sortableHandlers
    if (h) {
      el.removeEventListener('dragstart', h.onDragStart)
      el.removeEventListener('dragover', h.onDragOver)
      el.removeEventListener('drop', h.onDrop)
      el.removeEventListener('dragend', h.onDragEnd)
      el._sortableHandlers = null
    }
  }
}

// 依据 disabled 同步每个子项的 draggable 属性；有 handle 时默认关闭整项拖动
function applyDraggable(el) {
  const value = el._sortableValue || {}
  const disabled = !!value.disabled
  const hasHandle = !!value.handle
  el.querySelectorAll('[data-sortable-item]').forEach((item) => {
    const canDrag = !disabled && !hasHandle
    item.setAttribute('draggable', canDrag ? 'true' : 'false')
  })
}

export default sortable
