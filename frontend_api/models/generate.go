package models

// DesignItem 设计稿项
type DesignItem struct {
	Image       string `json:"image" binding:"required"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// GenerateRequest AI 代码生成请求（阻塞 & SSE 共用）
type GenerateRequest struct {
	Designs   []DesignItem `json:"designs" binding:"required,min=1"`
	Framework string       `json:"framework" binding:"required"`
	Quality   int          `json:"quality" binding:"required,min=1,max=100"`
}

// ScoreDimension 评分维度
type ScoreDimension struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Icon  string `json:"icon"`
}

// GenerateResponse AI 代码生成响应
type GenerateResponse struct {
	Code       string           `json:"code"`
	Score      int              `json:"score"`
	Dimensions []ScoreDimension `json:"dimensions"`
}
