package model

import (
	"time"
)

type Job struct {
	Id         int64     `json:"id"`          // 任务ID
	Name       string    `json:"name"`        // 名称
	Expr       string    `json:"expr"`        // 时间表达式
	Mark       string    `json:"mark"`        // 备注
	Enable     bool      `json:"enable"`      // 状态
	Path       string    `json:"path"`        // 脚本路径
	AreaId     int64     `json:"area_id"`     // 参数域ID
	CreateTime time.Time `json:"create_time"` // 创建时间
}

// ReqCreateJob
// 请求 创建定时任务结构体
type ReqCreateJob struct {
	Name   string `json:"name" binding:"required" label:"任务名称" example:"测试定时任务"`         // 名称
	Expr   string `json:"expr" binding:"required" label:"时间表达式" example:"0/5 * * * * ?"` // 时间表达式
	Mark   string `json:"mark" label:"任务说明" example:"每五秒执行一次脚本"`                         // 备注
	Path   string `json:"path" binding:"required" label:"脚本路径" example:"test.js"`        // 脚本路径
	AreaId int64  `json:"area_id" binding:"required" label:"参数域" example:"10000"`        // 参数域ID
}

// ReqModifyJob
// 请求 修改定时任务结构体
type ReqModifyJob struct {
	Id     int64  `json:"id" binding:"required" label:"任务ID"`                     // 任务ID
	Name   string `json:"name" binding:"required" label:"任务名称"`                   // 名称
	Expr   string `json:"expr" binding:"required" label:"时间表达式"`                  // 时间表达式
	Mark   string `json:"mark" label:"任务说明"`                                      // 备注
	Path   string `json:"path" binding:"required" label:"脚本路径"`                   // 脚本路径
	AreaId int64  `json:"area_id" binding:"required" label:"参数域" example:"10000"` // 参数域ID
}

// ReqQueryJob
// 请求 查询定时任务结构体
type ReqQueryJob struct {
	Name string `json:"name" form:"name"` // 名称
	PageInfo
}

type ResQueryJob struct {
	Id         int64  `json:"id"`          // 任务ID
	Name       string `json:"name"`        // 名称
	Expr       string `json:"expr"`        // 时间表达式
	Mark       string `json:"mark"`        // 备注
	Enable     bool   `json:"enable"`      // 状态
	Path       string `json:"path"`        // 脚本路径
	Area       string `json:"area"`        // 参数域
	AreaId     int64  `json:"area_id"`     //
	CreateTime string `json:"create_time"` // 创建时间
}
