package request

// CreateJournalRequest 创建期刊请求
type CreateJournalRequest struct {
	FullName     string  `json:"full_name" binding:"required"` // 期刊全称
	ShortName    string  `json:"short_name"`                   // 期刊简称
	ISSN         string  `json:"issn" binding:"required"`      // ISSN号
	ImpactFactor float64 `json:"impact_factor"`                // 影响因子
	Publisher    string  `json:"publisher"`                    // 出版商
}

// UpdateJournalRequest 更新期刊请求
type UpdateJournalRequest struct {
	FullName     string  `json:"full_name"`     // 期刊全称
	ShortName    string  `json:"short_name"`    // 期刊简称
	ImpactFactor float64 `json:"impact_factor"` // 影响因子
	Publisher    string  `json:"publisher"`     // 出版商
}
