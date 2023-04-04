package model

import "time"

type User struct {
	Id         int64     `json:"id"`          // 用户ID
	Account    string    `json:"account"`     // 登录名
	Name       string    `json:"name"`        // 用户名
	Password   string    `json:"password"`    // 用户密码
	Mark       string    `json:"mark"`        // 备注
	Enable     int64     `json:"enable"`      // 启用
	LoginTime  time.Time `json:"login_time"`  // 登录时间
	CreateTime time.Time `json:"create_time"` // 创建时间
}

type ReqCreateUser struct {
	Account  string `json:"account"`  // 登录名
	Name     string `json:"name"`     // 用户名
	Password string `json:"password"` // 用户密码
	Mark     string `json:"mark"`     // 备注
}

type ReqModifyUser struct {
	Id       int64  `json:"id"`       // 用户ID
	Account  string `json:"account"`  // 登录名
	Name     string `json:"name"`     // 用户名
	Password string `json:"password"` // 用户密码
	Mark     string `json:"mark"`     // 备注
}

type ReqQueryUser struct {
	Account string `json:"account" form:"account"` // 登录名
	Name    string `json:"name" form:"name"`       // 用户名
	Mark    string `json:"mark" form:"mark"`       // 备注
	PageInfo
}

type ResQueryUser struct {
	Id         int64  `json:"id"`          // 用户ID
	Account    string `json:"account"`     // 登录名
	Name       string `json:"name"`        // 用户名
	Mark       string `json:"mark"`        // 备注
	Enable     int64  `json:"enable"`      // 启用
	LoginTime  string `json:"login_time"`  // 登录时间
	CreateTime string `json:"create_time"` // 创建时间
}

type LoginUser struct {
	Account  string `json:"account"`  // 登录名
	Password string `json:"password"` // 用户密码
}
