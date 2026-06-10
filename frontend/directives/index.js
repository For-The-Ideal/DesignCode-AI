import title from './title.js'
import css from './css.js'
import wave from './wave.js'
import animate from './animate.js'
import intersect from './intersect.js'
import longpress from './longpress.js'
import borderGradient from './borderGradient.js'

/**
 * Vue指令注册入口
 *
 * 导出对应 { key: object } --> Vue.directive(key, object);
 * */

export default {
  'title': title,
  ...css,
  'wave': wave,
  'animate': animate,
  'intersect': intersect,
  'longpress': longpress,
  'border-gradient': borderGradient
}
