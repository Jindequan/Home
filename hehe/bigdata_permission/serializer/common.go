package serializer

//前端传入分页参数
type Page struct {
	PageSize int `json:"page_size" form:"page_size"`
	PageIndex int `json:"page_index" form:"page_index"`
}
//后端返回分页信息
type PageInfo struct {
	PageSize int `json:"page_size" form:"page_size"`
	PageIndex int `json:"page_index" form:"page_index"`
	Total int `json:"total" form:"total"`
	TotalPage int `json:"total_page" form:"total_page"`
}
//后端返回可视的时间字符串
type DaoTimeString struct {
	CreatedAtStr string `json:"created_at_str"`
	UpdatedAtStr string `json:"updated_at_str"`
	DeletedAtStr string `json:"deleted_at_str"`
}
