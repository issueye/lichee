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
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "一页大小， 如果不分页 传 0",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
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
                                    "$ref": "#/definitions/common.Full"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.Job"
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
                            "$ref": "#/definitions/common.Base"
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
                            "$ref": "#/definitions/common.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/common.Base"
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
                            "$ref": "#/definitions/common.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/common.Base"
                        }
                    }
                }
            },
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
                            "$ref": "#/definitions/common.Base"
                        }
                    },
                    "500": {
                        "description": "错误返回内容",
                        "schema": {
                            "$ref": "#/definitions/common.Base"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "common.Base": {
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
        "common.Full": {
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
        },
        "model.Job": {
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
                }
            }
        },
        "model.ReqCreateJob": {
            "type": "object",
            "required": [
                "expr",
                "name"
            ],
            "properties": {
                "expr": {
                    "description": "时间表达式",
                    "type": "string"
                },
                "mark": {
                    "description": "备注",
                    "type": "string"
                },
                "name": {
                    "description": "名称",
                    "type": "string"
                }
            }
        },
        "model.ReqModifyJob": {
            "type": "object",
            "required": [
                "expr",
                "id",
                "name"
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
