package model

// PageInfo 数据分页
type PageInfo struct {
	IsNotPaging bool  `json:"isNotPaging" form:"isNotPaging"` // 是否分页
	PageNum     int64 `json:"pageNum" form:"pageNum"`         // 页数
	PageSize    int64 `json:"pageSize" form:"pageSize"`       // 页码
	Total       int64 `json:"total"`                          // 总数  由服务器返回回去
}
