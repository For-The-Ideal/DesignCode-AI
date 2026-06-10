// 添加css单一属性-快捷指令
const cssDirectiveHandle = (item, el, bind) => {
  if (!bind.value) {
    return;
  }
  let val = bind.value;
  if (item === 'background' || item === 'backgroundColor') {
    val = val.indexOf('#') === -1 ? `#${val}` : val;
  } else if (item === 'backgroundImage') {
    val = val.indexOf('url') === -1 ? `url(${val})` : val;
  }
  if (item === 'color') {
    val = val.indexOf('#') > -1 ? val : val.indexOf('rgb') === -1 ? `#${val}` : val;
  }
  if (item === 'border') { // border 默认1px
    if (val.indexOf('#') > -1 || val.indexOf('rgb') > -1) {
      val = `1px solid ${val}`;
    } else {
      val = `1px solid #${val}`;
    }
  }
  el.style[item] = val;
};

const css = {}
const arr = [
  'background', 'backgroundColor', 'backgroundImage', 'width', 'height', 'color', 'border',
  'left', 'right', 'top', 'bottom', 'margin', 'marginTop', 'marginBottom', 'marginLeft', 'marginRight',
  'padding', 'paddingTop', 'paddingBottom', 'paddingLeft', 'paddingRight'
]
arr.forEach(item => {
  css[item] = {
    mounted(el, bind) {
      cssDirectiveHandle(item, el, bind)
    },
    updated(el, bind) {
      cssDirectiveHandle(item, el, bind)
    }
  }
})

export default css
