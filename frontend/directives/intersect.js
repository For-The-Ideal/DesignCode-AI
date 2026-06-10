
const isServer = import.meta.server;
const defaultConfig = {
  threshold: 0,
  root: null,
  rootMargin: '0px'
};
const defaultExtraData = {
  isRepeatIntersect: false,
  isIntersect: true,
  intersectKey: 'default',
  itemData: null,
};
const defaultCallback = () => { };
function parseIntersectValue(value) {
  const data = {
    config: {
      ...defaultConfig,
    },
    extraData: {
      ...defaultExtraData,
    },
    callback: defaultCallback
  };
  if (typeof value === 'function') {
    data.callback = value;
  } else {
    data.config = {
      ...data.config,
      ...(value.config || {})
    };
    data.extraData = {
      ...data.extraData,
      ...(value.extraData || {})
    };
    data.callback = value.callback || defaultCallback;
  }
  return data;
}
const intersectionObserverMap = {};
let map = null;
const customIntersect = {
  mounted: function (el, binding) {
    if (isServer) return;
    const { config, extraData, callback } = parseIntersectValue(binding.value || {});
    if (!extraData.isIntersect) return;
    try {
      if (!map) {
        map = new Map();
      }
      const intersectionObserverKey = extraData.intersectKey || 'default';
      if (!intersectionObserverMap[intersectionObserverKey]) {
        intersectionObserverMap[intersectionObserverKey] = new IntersectionObserver((entries) => {
          // callback(entry, extraData);
          entries.forEach(entry => {
            const intersecting = entry.isIntersecting || entry.isIntersecting === void 0;
            if (intersecting) {
              const curMap = map.get(entry.target);
              if (!curMap) return;
              const extraData = curMap.extraData || {};
              curMap.callback(curMap.extraData, entry);
              if (!extraData.isRepeatIntersect) {
                intersectionObserverMap[intersectionObserverKey].unobserve(entry.target);
                map.delete(entry.target);
              }
            }
          });
        }, config);
      }
      intersectionObserverMap[intersectionObserverKey].observe(el);
      map.set(el, { extraData, callback });
    } catch (error) {
      callback(extraData, null);
      console.log('intersect-error', error);
    }
  },
  unmounted: function (el, binding) {
    if (!map) return;
    const curMap = map.get(el);
    if (!curMap) return;
    const extraData = curMap.extraData || {};
    const intersectionObserverKey = extraData.intersectKey || 'default';
    if (intersectionObserverMap[intersectionObserverKey]) {
      intersectionObserverMap[intersectionObserverKey].unobserve(el);
    }
    map.delete(el);
  }
}

export default customIntersect
