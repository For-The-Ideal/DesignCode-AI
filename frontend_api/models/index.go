package models

// AllModels 包含所有需要进行数据库迁移的模型
// 当你有 100 个模型时，只需要在这里添加一行即可
var AllModels = []interface{}{
	&User{},
}
