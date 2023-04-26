package model

type DbSource struct {
	Id       int64  `json:"id"`       // ID
	Name     string `json:"name"`     // 名称
	Host     string `json:"host"`     // 地址
	Port     int64  `json:"port"`     // 端口号
	Database string `json:"database"` // 数据库
	User     string `json:"user"`     // 账号
	Password string `json:"password"` // 密码
	Type     int64  `json:"type"`     // 类型 0 sqlserver 1 mysql 2 oracle
	Mark     string `json:"mark"`     // 备注
}

// ReqCreateDbSource
// 创建
type ReqCreateDbSource struct {
	Name     string `json:"name" binding:"required" label:"名称" example:"测试SQLSERVER_001"` // 名称
	Host     string `json:"host" binding:"required" label:"服务地址" example:"."`             // 地址
	Port     int64  `json:"port" binding:"required" label:"端口号" example:"1433"`           // 端口号
	Database string `json:"database" binding:"required" label:"数据库名称" example:"TEST"`     // 数据库
	User     string `json:"user" binding:"required" label:"账户" example:"sa"`              // 账号
	Password string `json:"password" binding:"required" label:"密码" example:"123456"`      // 密码
	Type     int64  `json:"type" label:"数据库类型" example:"0"`                               // 类型
	Mark     string `json:"mark"`                                                         // 备注
}

// 修改
type ReqModifyDbSource struct {
	Id       int64  `json:"id"  binding:"required" label:"id" example:"10000"`            // ID
	Name     string `json:"name" binding:"required" label:"名称" example:"测试SQLSERVER_001"` // 名称
	Host     string `json:"host" binding:"required" label:"服务地址" example:"."`             // 地址
	Port     int64  `json:"port" binding:"required" label:"端口号" example:"1433"`           // 端口号
	Database string `json:"database" binding:"required" label:"数据库名称" example:"TEST"`     // 数据库
	User     string `json:"user" binding:"required" label:"账户" example:"sa"`              // 账号
	Password string `json:"password" binding:"required" label:"密码" example:"123456"`      // 密码
	Type     int64  `json:"type" label:"数据库类型" example:"0"`                               // 类型
	Mark     string `json:"mark"`                                                         // 备注
}

// ReqQueryDbSource
// 请求 查询
type ReqQueryDbSource struct {
	Name string `json:"name" form:"name"` // 名称
	PageInfo
}
