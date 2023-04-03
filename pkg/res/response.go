package res

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/issueye/lichee/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	orange_validator "github.com/issueye/lichee/pkg/validator"
)

type RCode int

const (
	// OK (成功)
	OK RCode = 200
	// 创建成功
	OK_CREATE RCode = 201
	// 错误请求, 参数错误
	BAD_REQUEST RCode = 400
	// 未授权请求
	UNAUTHORIZED RCode = 401
	// 找不到资源
	NOT_FOUND RCode = 404
	// 请求方法不允许
	METHOD_NOT_ALLOWED RCode = 405
	// 服务器内部错误
	INTERNAL_SERVER_ERROR RCode = 500
	// 未实现
	NOT_IMPLEMENTED RCode = 501
	// 错误网关
	BAD_GATEWAY RCode = 502
	// 网关超时
	GATEWAY_TIMEOUT RCode = 504
	// HTTP版本不支持
	HTTP_VERSION_NOT_SUPPORTED RCode = 505

	// 16+ 创建
	CREATE_FAIL = 1601
	// 17+ 修改
	MODIFY_FAIL = 1701
	// 18+ 删除
	DELETE_FAIL = 1801
	// 19+ 查询
	QUERY_FAIL = 1091
)

var CodeMessage = map[RCode]string{
	OK:                         "成功",
	OK_CREATE:                  "创建成功",
	BAD_REQUEST:                "错误请求, 参数错误",
	UNAUTHORIZED:               "未授权请求",
	NOT_FOUND:                  "找不到资源",
	METHOD_NOT_ALLOWED:         "请求方法不允许",
	INTERNAL_SERVER_ERROR:      "服务器内部错误",
	NOT_IMPLEMENTED:            "未实现",
	BAD_GATEWAY:                "错误网关",
	GATEWAY_TIMEOUT:            "网关超时",
	HTTP_VERSION_NOT_SUPPORTED: "HTTP版本不支持",

	CREATE_FAIL: "创建数据失败",
	MODIFY_FAIL: "修改数据失败",
	DELETE_FAIL: "删除数据失败",
	QUERY_FAIL:  "查询数据失败",
}

type PageInfo struct {
	IsNotPaging bool  `json:"isNotPaging" form:"isNotPaging"` // 是否分页
	PageNum     int64 `json:"pageNum" form:"pageNum"`         // 页数
	PageSize    int64 `json:"pageSize" form:"pageSize"`       // 页码
	Total       int64 `json:"total"`                          // 总数  由服务器返回回去
}

type Base struct {
	Code             RCode  `json:"code"`             // 状态码
	Message          string `json:"message"`          // 状态消息
	RequestDatetime  string `json:"requestDatetime"`  // 请求时间
	ResponseDatetime string `json:"responseDatetime"` // 返回时间
	RequestId        string `json:"requestId"`        // 请求ID
}

type Full struct {
	Base
	Data interface{} `json:"data"` // 返回数据
}

type Null struct {
	Base
}

type Page struct {
	Base
	Page PageInfo    `json:"pageInfo"` // 分页信息
	Data interface{} `json:"data"`     // 返回数据
}

// SuccessData
// 返回成功，并且包含所有返回数据
func SuccessData(ctx *gin.Context, data interface{}) {
	res := new(Full)
	res.Code = OK
	res.RequestDatetime = ctx.GetString("RQ_DATETIME")
	res.RequestId = ctx.GetString("RQ_ID")
	res.ResponseDatetime = time.Now().Format(utils.FormatDateTimeMs)
	res.Message = CodeMessage[OK]
	res.Data = data

	ctx.Set("res", res)
	ctx.JSON(http.StatusOK, res)
}

// SuccessAutoData
// 根据数据自动判断是返回所有数据还是部分数据
func SuccessAutoData(ctx *gin.Context, req interface{}, data interface{}) {
	ref := reflect.ValueOf(req)
	// 判断 ref 是否是 ptr 类型
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}

	// 判断 req 是否有 Total 字段
	// 如果有则将count 赋值给 req.Total
	if !ref.FieldByName("Total").IsValid() {
		SuccessData(ctx, data)
		return
	}

	num := ref.FieldByName("PageNum").Int()
	size := ref.FieldByName("PageSize").Int()

	paging := false
	if num > 0 {
		paging = true
	}

	Total := ref.FieldByName("Total").Int()
	if paging {
		ctx.Set("RQ_PAGE_TOTAL", Total)
		ctx.Set("RQ_PAGE_NUM", num)
		ctx.Set("RQ_PAGE_SIZE", size)

		// 获取数据起点和终点
		begin, end := utils.Slice{}.SlicePage(int(num), int(size), int(Total))
		dataRef := reflect.ValueOf(data)
		if dataRef.Kind() == reflect.Ptr {
			dataRef = dataRef.Elem()
		}

		// 判断是否是切片，如果是切片则获取切片的对应数据
		if dataRef.Kind() == reflect.Slice {
			// 判断传入的切片长度，如果传入的长度和 size相等，则不需要再重新切
			if dataRef.Len() >= int(size) {
				data = dataRef.Slice(begin, end).Interface()
			}
		}
		SuccessPage(ctx, data)
		return
	}
	SuccessData(ctx, data)
}

// Success
// 返回成功，不包含数据
func Success(ctx *gin.Context) {
	res := new(Null)
	res.Code = OK
	res.RequestDatetime = ctx.GetString("RQ_DATETIME")
	res.RequestId = ctx.GetString("RQ_ID")
	res.ResponseDatetime = time.Now().Format(utils.FormatDateTimeMs)
	res.Message = CodeMessage[OK]

	ctx.Set("res", res)
	ctx.JSON(http.StatusOK, res)
}

// SuccessByMsg
// 返回成功, 并且包含消息
func SuccessByMsg(ctx *gin.Context, message string) {
	res := new(Base)
	res.Code = OK
	res.RequestDatetime = ctx.GetString("RQ_DATETIME")
	res.RequestId = ctx.GetString("RQ_ID")
	res.ResponseDatetime = time.Now().Format(utils.FormatDateTimeMs)
	res.Message = message
	ctx.Set("res", res)
	ctx.JSON(http.StatusOK, res)
}

// SuccessByMsgf
// 返回成功, 并且包含消息
func SuccessByMsgf(ctx *gin.Context, fmtStr string, args ...any) {
	message := fmt.Sprintf(fmtStr, args...)
	res := new(Base)
	res.Code = OK
	res.RequestDatetime = ctx.GetString("RQ_DATETIME")
	res.RequestId = ctx.GetString("RQ_ID")
	res.ResponseDatetime = time.Now().Format(utils.FormatDateTimeMs)
	res.Message = message
	ctx.Set("res", res)
	ctx.JSON(http.StatusOK, res)
}

// SuccessPage
// 返回成功，包含分页之后的数据
func SuccessPage(ctx *gin.Context, data interface{}) {
	res := new(Page)
	res.Code = OK
	res.RequestDatetime = ctx.GetString("RQ_DATETIME")
	res.RequestId = ctx.GetString("RQ_ID")
	res.ResponseDatetime = time.Now().Format(utils.FormatDateTimeMs)
	res.Message = CodeMessage[OK]
	res.Data = data
	res.Page.PageNum = ctx.GetInt64("RQ_PAGE_NUM")
	res.Page.PageSize = ctx.GetInt64("RQ_PAGE_SIZE")
	res.Page.Total = ctx.GetInt64("RQ_PAGE_TOTAL")

	ctx.Set("res", res)
	ctx.JSON(http.StatusOK, res)
}

// Fail
// 返回失败
func Fail(ctx *gin.Context, code RCode) {
	res := new(Base)
	res.Code = code
	res.RequestDatetime = ctx.GetString("RQ_DATETIME")
	res.RequestId = ctx.GetString("RQ_ID")
	res.ResponseDatetime = time.Now().Format(utils.FormatDateTimeMs)
	res.Message = CodeMessage[code]

	ctx.Set("res", res)
	ctx.JSON(200, res)
}

// FailByMsg
// 返回失败，并且返回自定义错误信息
func FailByMsg(ctx *gin.Context, msg string) {
	res := new(Base)
	res.Code = INTERNAL_SERVER_ERROR
	res.RequestDatetime = ctx.GetString("RQ_DATETIME")
	res.RequestId = ctx.GetString("RQ_ID")
	res.ResponseDatetime = time.Now().Format(utils.FormatDateTimeMs)
	res.Message = msg
	ctx.Set("res", res)
	ctx.JSON(http.StatusOK, res)
}

// FailByMsgf
// 返回失败，并且返回自定义错误信息
func FailByMsgf(ctx *gin.Context, fmtStr string, args ...any) {
	message := fmt.Sprintf(fmtStr, args...)
	res := new(Base)
	res.Code = INTERNAL_SERVER_ERROR
	res.RequestDatetime = ctx.GetString("RQ_DATETIME")
	res.RequestId = ctx.GetString("RQ_ID")
	res.ResponseDatetime = time.Now().Format(utils.FormatDateTimeMs)
	res.Message = message
	ctx.Set("res", res)
	ctx.JSON(http.StatusOK, res)
}

// FailByMsgAndCode
// 返回失败，并且返回自定义错误信息和状态码
func FailByMsgAndCode(ctx *gin.Context, msg string, code RCode) {
	res := new(Base)
	res.Code = code
	res.RequestDatetime = ctx.GetString("RQ_DATETIME")
	res.RequestId = ctx.GetString("RQ_ID")
	res.ResponseDatetime = time.Now().Format(utils.FormatDateTimeMs)
	res.Message = msg
	ctx.Set("res", res)
	ctx.JSON(http.StatusBadRequest, res)
}

// FailBind
// 表单验证失败
func FailBind(ctx *gin.Context, err error) {
	// 判断错误类型
	var errStr string
	switch t := err.(type) {
	case validator.ValidationErrors:
		errList := t.Translate(orange_validator.Trans)
		for _, element := range errList {
			errStr += element + " "
		}
	default:
		errStr = err.Error()
	}

	FailByMsgAndCode(ctx, errStr, BAD_REQUEST)
}
