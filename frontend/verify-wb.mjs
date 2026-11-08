import fs from 'node:fs'
import { parse, compileScript, compileTemplate } from '@vue/compiler-sfc'

const file = 'components/workbench/WorkbenchUpload.vue'
const source = fs.readFileSync(file, 'utf8')
const { descriptor, errors } = parse(source, { filename: file })
if (errors.length) {
  console.error('PARSE_ERR', errors)
  process.exit(1)
}
const script = compileScript(descriptor, { id: 'x' })
const tpl = compileTemplate({
  source: descriptor.template.content,
  filename: file,
  id: 'x',
  compilerOptions: { bindingMetadata: script.bindings },
})
if (tpl.errors.length) {
  console.error('TPL_ERR', tpl.errors)
  process.exit(1)
}
console.log('ALL_OK')
