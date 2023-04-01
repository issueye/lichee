package model

import (
	"bytes"
	"encoding/gob"

	"github.com/issueye/lichee/utils"
)

type Job struct {
	Id         int64              `json:"id"`          // 任务ID
	Name       string             `json:"name"`        // 名称
	Expr       string             `json:"expr"`        // 时间表达式
	Mark       string             `json:"mark"`        // 备注
	Enable     bool               `json:"enable"`      // 状态
	CreateTime utils.LongDateTime `json:"create_time"` // 创建时间
}

// ReqCreateJob
// 请求 创建定时任务结构体
type ReqCreateJob struct {
	Name string `json:"name" binding:"required" label:"任务名称"`  // 名称
	Expr string `json:"expr" binding:"required" label:"时间表达式"` // 时间表达式
	Mark string `json:"mark" label:"任务说明"`                     // 备注
}

// ReqModifyJob
// 请求 修改定时任务结构体
type ReqModifyJob struct {
	Id   int64  `json:"id" binding:"required" label:"任务ID"`    // 任务ID
	Name string `json:"name" binding:"required" label:"任务名称"`  // 名称
	Expr string `json:"expr" binding:"required" label:"时间表达式"` // 时间表达式
	Mark string `json:"mark" label:"任务说明"`                     // 备注
}

// ReqQueryJob
// 请求 查询定时任务结构体
type ReqQueryJob struct {
	Name string `json:"name" form:"name"` // 名称
	Mark string `json:"mark" form:"mark"` // 备注
	PageInfo
}

func (job Job) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	err := enc.Encode(job.CreateTime)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (job *Job) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&job.Id); err != nil {
		return err
	}

	if err := dec.Decode(&job.Name); err != nil {
		return err
	}

	if err := dec.Decode(&job.Mark); err != nil {
		return err
	}

	if err := dec.Decode(&job.Enable); err != nil {
		return err
	}

	if err := dec.Decode(&job.Expr); err != nil {
		return err
	}

	var createTime utils.LongDateTime
	if err := dec.Decode(&createTime); err != nil {
		return err
	}

	job.CreateTime = createTime

	return nil
}
