package model

import "time"

type Job struct {
	Id         int64     `json:"id"`          // 任务ID
	Name       string    `json:"name"`        // 名称
	Expr       string    `json:"expr"`        // 时间表达式
	Mark       string    `json:"mark"`        // 备注
	Enable     bool      `json:"enable"`      // 状态
	CreateTime time.Time `json:"create_time"` // 创建时间
}

type ReqCreateJob struct {
	Name   string `json:"name"`   // 名称
	Expr   string `json:"expr"`   // 时间表达式
	Mark   string `json:"mark"`   // 备注
	Enable bool   `json:"enable"` // 状态
}
