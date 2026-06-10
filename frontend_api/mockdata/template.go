package mockdata

import "fmt"

// ═══════════════════════════════════════════════
//  模板数据入口
//
//  当前阶段：所有模拟数据写死在代码中
//  后期对接数据库：只需替换 GetTemplate 内部实现，
//  外部调用方无需任何改动。
// ═══════════════════════════════════════════════

// GetTemplate 根据模板ID返回对应的模板数据
//
// 参数:
//   id - 模板编号（当前仅支持 1）
//
// 返回:
//   *Template1Data - 模板数据
//   error          - 模板不存在时返回错误
//
// 🔌 数据库接入路径：
//   1. 将 switch 内硬编码替换为 db.Query("SELECT ... WHERE id=?", id)
//   2. 将返回类型改为 interface{} 或泛型以支持多套模板结构
//   3. 删除 getTemplate1() 及相关模拟数据文件
func GetTemplate(id int) (*Template1Data, error) {
	// 📦 模拟数据分支（数据库接入后替换此 switch）
	switch id {
	case 1:
		return getTemplate1(), nil
	// 预留：后续新增模板套件时在此追加 case
	// case 2: return getTemplate2()
	// case 3: return getTemplate3()
	default:
		return nil, fmt.Errorf("模板 ID=%d 不存在", id)
	}
}
