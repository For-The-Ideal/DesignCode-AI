// 页面关系边表（多图模式）
// 真源：state.modes[mode].pageLinks；每页的 parentPageId 只是旧数据的派生视图，
// 新交互（关系段/关系板）统一写边表，提交时由 task.js 转成 page_links。
// kind 是「关系类型」（route 路由跳转 / modal 弹层 / tab 底部 Tab / back 返回 / auto 自动），
// trigger 是后端可消费的触发方式，由 kind 派生，保证提交字段始终合法。
export const KIND_TRIGGER = { route: 'node_click', modal: 'node_click', tab: 'bottom_tab', back: 'back', auto: 'auto' }

export const createRelationActions = (ctx) => {
  const { images, pageLinks } = ctx

  const uid = () => `r${Date.now()}_${Math.random().toString(36).slice(2, 7)}`
  const pageIdSet = () => new Set(images.value.map((i) => i.pageId))

  /** 旧数据迁移：边表为空时，把各页 parentPageId 转成显式边（迁移后清空旧字段，避免删边后又复活） */
  const ensurePageLinks = () => {
    if (pageLinks.value.length) return
    const ids = pageIdSet()
    let migrated = false
    images.value.forEach((img) => {
      const from = img.parentPageId
      if (!from || !ids.has(from) || from === img.pageId) return
      if (pageLinks.value.some((l) => l.fromPageId === from && l.toPageId === img.pageId)) return
      pageLinks.value.push({
        id: uid(),
        fromPageId: from,
        toPageId: img.pageId,
        trigger: img.linkTrigger || 'node_click',
        kind: 'route',
        nodeKey: '',
        label: '',
        desc: '',
      })
      migrated = true
    })
    if (migrated) images.value.forEach((img) => { img.parentPageId = '' })
  }

  /** 新增显式边：页面必须存在、禁止自环、同 from→to 只保留一条（重复时返回已有边）
   *  extra：{ kind 关系类型, nodeKey 源组件节点, label 源组件名, desc 说明 } */
  const addPageLink = (fromPageId, toPageId, trigger = 'node_click', extra = {}) => {
    const ids = pageIdSet()
    if (!ids.has(fromPageId) || !ids.has(toPageId) || fromPageId === toPageId) return null
    const exist = pageLinks.value.find((l) => l.fromPageId === fromPageId && l.toPageId === toPageId)
    if (exist) return exist
    const kind = extra.kind || 'route'
    const link = {
      id: uid(),
      fromPageId,
      toPageId,
      kind,
      trigger: KIND_TRIGGER[kind] || trigger || 'node_click',
      nodeKey: extra.nodeKey || '',
      label: extra.label || '',
      desc: extra.desc || '',
    }
    pageLinks.value.push(link)
    return link
  }

  /** 修改边（源/目标页/关系类型/源组件/说明）；非法引用或改重复时保持原值 */
  const updatePageLink = (id, patch = {}) => {
    const link = pageLinks.value.find((l) => l.id === id)
    if (!link) return
    if (patch.fromPageId && patch.fromPageId !== link.fromPageId) {
      const ids = pageIdSet()
      if (!ids.has(patch.fromPageId) || patch.fromPageId === link.toPageId) return
      if (pageLinks.value.some((l) => l.id !== id && l.fromPageId === patch.fromPageId && l.toPageId === link.toPageId)) return
      link.fromPageId = patch.fromPageId
    }
    if (patch.toPageId && patch.toPageId !== link.toPageId) {
      const ids = pageIdSet()
      if (!ids.has(patch.toPageId) || patch.toPageId === link.fromPageId) return
      if (pageLinks.value.some((l) => l.id !== id && l.fromPageId === link.fromPageId && l.toPageId === patch.toPageId)) return
      link.toPageId = patch.toPageId
    }
    if (patch.kind) {
      link.kind = patch.kind
      link.trigger = KIND_TRIGGER[patch.kind] || link.trigger
    } else if (patch.trigger) {
      link.trigger = patch.trigger
    }
    if (patch.nodeKey !== undefined) link.nodeKey = patch.nodeKey || ''
    if (patch.label !== undefined) link.label = String(patch.label || '').slice(0, 40)
    if (patch.desc !== undefined) link.desc = String(patch.desc).slice(0, 100)
  }

  const removePageLink = (id) => {
    const i = pageLinks.value.findIndex((l) => l.id === id)
    if (i > -1) pageLinks.value.splice(i, 1)
  }

  /** 清理悬空/自环边（页面删除、重排后调用） */
  const prunePageLinks = () => {
    const ids = pageIdSet()
    const keep = pageLinks.value.filter(
      (l) => ids.has(l.fromPageId) && ids.has(l.toPageId) && l.fromPageId !== l.toPageId,
    )
    pageLinks.value.splice(0, pageLinks.value.length, ...keep)
  }

  const linksFrom = (pageId) => pageLinks.value.filter((l) => l.fromPageId === pageId)
  const linksTo = (pageId) => pageLinks.value.filter((l) => l.toPageId === pageId)

  return { ensurePageLinks, addPageLink, updatePageLink, removePageLink, prunePageLinks, linksFrom, linksTo }
}
