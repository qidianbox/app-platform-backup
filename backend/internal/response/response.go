package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// 错误码定义
const (
	CodeSuccess          = 0
	CodeBadRequest       = 400
	CodeUnauthorized     = 401
	CodeForbidden        = 403
	CodeNotFound         = 404
	CodeConflict         = 409
	CodeTooManyRequests  = 429
	CodeInternalError    = 500
	CodeServiceUnavailable = 503
)

// 错误消息定义
var codeMessages = map[int]string{
	CodeSuccess:          "success",
	CodeBadRequest:       "Bad request",
	CodeUnauthorized:     "Unauthorized",
	CodeForbidden:        "Forbidden",
	CodeNotFound:         "Not found",
	CodeConflict:         "Conflict",
	CodeTooManyRequests:  "Too many requests",
	CodeInternalError:    "Internal server error",
	CodeServiceUnavailable: "Service unavailable",
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data: PageData{
			List:  list,
			Total: total,
			Page:  page,
			Size:  size,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeBadRequest]
	}
	Error(c, CodeBadRequest, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeUnauthorized]
	}
	Error(c, CodeUnauthorized, message)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeForbidden]
	}
	Error(c, CodeForbidden, message)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeNotFound]
	}
	Error(c, CodeNotFound, message)
}

// Conflict 409错误
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeConflict]
	}
	Error(c, CodeConflict, message)
}

// TooManyRequests 429错误
func TooManyRequests(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeTooManyRequests]
	}
	Error(c, CodeTooManyRequests, message)
}

// InternalError 500错误
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeInternalError]
	}
	Error(c, CodeInternalError, message)
}

// ServiceUnavailable 503错误
func ServiceUnavailable(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeServiceUnavailable]
	}
	Error(c, CodeServiceUnavailable, message)
}

// ParamError 参数错误
func ParamError(c *gin.Context, message string) {
	if message == "" {
		message = "参数错误"
	}
	BadRequest(c, message)
}

// DBError 数据库错误
func DBError(c *gin.Context, err error) {
	message := "数据库操作失败"
	if err != nil {
		// 生产环境不应该暴露详细错误信息
		// message = err.Error()
	}
	InternalError(c, message)
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, list interface{}, total int64, page, size int) {
	SuccessPage(c, list, total, page, size)
}

// getHTTPStatus 根据业务码获取HTTP状态码
func getHTTPStatus(code int) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	case CodeServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
