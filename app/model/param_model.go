package model

type Param struct {
	Id     int64  `json:"id"`      // id
	Name   string `json:"name"`    // 参数名称
	AreaId int64  `json:"area_id"` // 参数域id
	Area   string `json:"area"`    // 参数域
	Value  string `json:"value"`   // 参数值
	Mark   string `json:"mark"`    // 备注
}

type ParamArea struct {
	Id   int64  `json:"id"`   // id
	Name string `json:"name"` // 域名称
	Mark string `json:"mark"` // 备注
}

type ReqCreateParam struct {
	Name   string `json:"name"`    // 参数名称
	AreaId int64  `json:"area_id"` // 参数域id
	Value  string `json:"value"`   // 参数值
	Mark   string `json:"mark"`    // 备注
}

type ReqModifyParam struct {
	Id     int64  `json:"id"`      // id
	Name   string `json:"name"`    // 参数名称
	AreaId int64  `json:"area_id"` // 参数域id
	Area   string `json:"area"`    // 参数域
	Value  string `json:"value"`   // 参数值
	Mark   string `json:"mark"`    // 备注
}

type ReqQueryParam struct {
	Name   string `json:"name" form:"name"`       // 参数名称
	AreaId int64  `json:"area_id" form:"area_id"` // 参数域id
	PageInfo
}

type ReqCreateArea struct {
	Name string `json:"name"` // 域名称
	Mark string `json:"mark"` // 备注
}

type ReqQueryArea struct {
	Name string `json:"name" form:"name"` // 参数名称
	PageInfo
}

type ReqModifyArea struct {
	Id   int64  `json:"id"`   // id
	Name string `json:"name"` // 域名称
	Mark string `json:"mark"` // 备注
}
