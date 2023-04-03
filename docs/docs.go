// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/job": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取定时任务列表",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "定时任务管理"
                ],
                "summary": "获取定时任务列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "是否需要分页， 默认需要， 如果不分页 传 true",
                        "name": "isNotPaging",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "页码， 如果不分页 传 0",
                        "name": "pageNum",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "一页大小， 如果不分页 传 0",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "任务名称",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "任务描述",
                        "name": "mark",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/res.Full"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.ResQueryJob"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "修改定时任务数据",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "定时任务管理"
                ],
                "summary": "修改定时任务数据",
                "parameters": [
                    {
                        "description": "修改定时任务数据",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ReqModifyJob"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "添加定时任务数据",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "定时任务管理"
                ],
                "summary": "添加定时任务数据",
                "parameters": [
                    {
                        "description": "添加定时任务数据",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ReqCreateJob"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            }
        },
        "/api/job/status/{id}": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "修改定时任务数据",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "定时任务管理"
                ],
                "summary": "修改定时任务数据",
                "parameters": [
                    {
                        "type": "string",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            }
        },
        "/api/job/{id}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "删除定时任务",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "定时任务管理"
                ],
                "summary": "删除定时任务",
                "parameters": [
                    {
                        "type": "string",
                        "description": "任务ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            }
        },
        "/api/login": {
            "post": {
                "description": "` + "`" + `` + "`" + `` + "`" + `\n用户登录\n` + "`" + `` + "`" + `` + "`" + `",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "基本接口"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "登录信息",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/middleware.LoginUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/res.Full"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            }
        },
        "/api/logout": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "用户登出时，调用此接口",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "基本接口"
                ],
                "summary": "用户登出",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            }
        },
        "/api/refreshToken": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "当token即将获取或者过期时刷新token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "基本接口"
                ],
                "summary": "刷新token",
                "responses": {
                    "200": {
                        "description": "code:200 成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/res.Full"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/middleware.JwtToken"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            }
        },
        "/api/user": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "获取用户列表",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理"
                ],
                "summary": "获取用户列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "是否需要分页， 默认需要， 如果不分页 传 true",
                        "name": "isNotPaging",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "页码， 如果不分页 传 0",
                        "name": "pageNum",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "一页大小， 如果不分页 传 0",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "任务名称",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "任务名称",
                        "name": "account",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "任务描述",
                        "name": "mark",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/res.Full"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.ReqQueryUser"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "修改用户信息",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理"
                ],
                "summary": "修改用户信息",
                "parameters": [
                    {
                        "description": "修改用户信息",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ReqModifyUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "添加用户",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理"
                ],
                "summary": "添加用户",
                "parameters": [
                    {
                        "description": "添加用户",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ReqCreateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            }
        },
        "/api/user/status/{id}": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "停用/启用用户",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理"
                ],
                "summary": "停用/启用用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            }
        },
        "/api/user/{id}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "删除用户",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理"
                ],
                "summary": "删除用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code: 200 成功",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/res.Base"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "middleware.JwtToken": {
            "type": "object",
            "properties": {
                "expires": {
                    "description": "时间",
                    "type": "string"
                },
                "id": {
                    "description": "id",
                    "type": "integer"
                },
                "name": {
                    "description": "用户名",
                    "type": "string"
                },
                "token": {
                    "description": "token",
                    "type": "string"
                },
                "uid": {
                    "description": "用户ID",
                    "type": "string"
                }
            }
        },
        "middleware.LoginUser": {
            "type": "object",
            "properties": {
                "account": {
                    "description": "登录名",
                    "type": "string"
                },
                "password": {
                    "description": "用户密码",
                    "type": "string"
                }
            }
        },
        "model.ReqCreateJob": {
            "type": "object",
            "required": [
                "expr",
                "name",
                "path"
            ],
            "properties": {
                "expr": {
                    "description": "时间表达式",
                    "type": "string",
                    "example": "0/5 * * * * ?"
                },
                "mark": {
                    "description": "备注",
                    "type": "string",
                    "example": "每五秒执行一次脚本"
                },
                "name": {
                    "description": "名称",
                    "type": "string",
                    "example": "测试定时任务"
                },
                "path": {
                    "description": "脚本路径",
                    "type": "string",
                    "example": "test.js"
                }
            }
        },
        "model.ReqCreateUser": {
            "type": "object",
            "properties": {
                "account": {
                    "description": "登录名",
                    "type": "string"
                },
                "mark": {
                    "description": "备注",
                    "type": "string"
                },
                "name": {
                    "description": "用户名",
                    "type": "string"
                },
                "password": {
                    "description": "用户密码",
                    "type": "string"
                }
            }
        },
        "model.ReqModifyJob": {
            "type": "object",
            "required": [
                "expr",
                "id",
                "name",
                "path"
            ],
            "properties": {
                "expr": {
                    "description": "时间表达式",
                    "type": "string"
                },
                "id": {
                    "description": "任务ID",
                    "type": "integer"
                },
                "mark": {
                    "description": "备注",
                    "type": "string"
                },
                "name": {
                    "description": "名称",
                    "type": "string"
                },
                "path": {
                    "description": "脚本路径",
                    "type": "string"
                }
            }
        },
        "model.ReqModifyUser": {
            "type": "object",
            "properties": {
                "account": {
                    "description": "登录名",
                    "type": "string"
                },
                "id": {
                    "description": "用户ID",
                    "type": "integer"
                },
                "mark": {
                    "description": "备注",
                    "type": "string"
                },
                "name": {
                    "description": "用户名",
                    "type": "string"
                },
                "password": {
                    "description": "用户密码",
                    "type": "string"
                }
            }
        },
        "model.ReqQueryUser": {
            "type": "object",
            "properties": {
                "account": {
                    "description": "登录名",
                    "type": "string"
                },
                "isNotPaging": {
                    "description": "是否分页",
                    "type": "boolean"
                },
                "mark": {
                    "description": "备注",
                    "type": "string"
                },
                "name": {
                    "description": "用户名",
                    "type": "string"
                },
                "pageNum": {
                    "description": "页数",
                    "type": "integer"
                },
                "pageSize": {
                    "description": "页码",
                    "type": "integer"
                },
                "total": {
                    "description": "总数  由服务器返回回去",
                    "type": "integer"
                }
            }
        },
        "model.ResQueryJob": {
            "type": "object",
            "properties": {
                "create_time": {
                    "description": "创建时间",
                    "type": "string"
                },
                "enable": {
                    "description": "状态",
                    "type": "boolean"
                },
                "expr": {
                    "description": "时间表达式",
                    "type": "string"
                },
                "id": {
                    "description": "任务ID",
                    "type": "integer"
                },
                "mark": {
                    "description": "备注",
                    "type": "string"
                },
                "name": {
                    "description": "名称",
                    "type": "string"
                },
                "path": {
                    "description": "脚本路径",
                    "type": "string"
                }
            }
        },
        "res.Base": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "状态码",
                    "type": "integer"
                },
                "message": {
                    "description": "状态消息",
                    "type": "string"
                },
                "requestDatetime": {
                    "description": "请求时间",
                    "type": "string"
                },
                "requestId": {
                    "description": "请求ID",
                    "type": "string"
                },
                "responseDatetime": {
                    "description": "返回时间",
                    "type": "string"
                }
            }
        },
        "res.Full": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "状态码",
                    "type": "integer"
                },
                "data": {
                    "description": "返回数据"
                },
                "message": {
                    "description": "状态消息",
                    "type": "string"
                },
                "requestDatetime": {
                    "description": "请求时间",
                    "type": "string"
                },
                "requestId": {
                    "description": "请求ID",
                    "type": "string"
                },
                "responseDatetime": {
                    "description": "返回时间",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
