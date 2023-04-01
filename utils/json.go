package utils

import (
	"encoding/json"
	"fmt"
)

type Ljson struct{}

// 结构体转为json
func (l Ljson) Struct2Json(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("[Struct2Json]转换异常: %v", err))
	}
	return string(str)
}

// 结构体转为json
func (l Ljson) Struct2JsonFmt(obj interface{}) string {
	str, err := json.MarshalIndent(obj, " ", " ")
	if err != nil {
		panic(fmt.Sprintf("[Struct2Json]转换异常: %v", err))
	}
	return string(str)
}

// json转为结构体
func (l Ljson) Json2Struct(str string, obj interface{}) {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		panic(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}

// json interface转为结构体
func (l Ljson) JsonI2Struct(str interface{}, obj interface{}) {
	JsonStr := str.(string)
	l.Json2Struct(JsonStr, obj)
}

// 结构体转结构体, json为中间桥梁, struct2必须以指针方式传递, 否则可能获取到空数据
func (l Ljson) Struct2StructByJson(struct1 interface{}, struct2 interface{}) {
	// 转换为响应结构体, 隐藏部分字段
	jsonStr := l.Struct2Json(struct1)
	l.Json2Struct(jsonStr, struct2)
}
