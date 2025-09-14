package models

// Response 统一响应结构
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// 常见错误代码
const (
	CodeSuccess      = "SUCCESS"
	CodeInvalidParam = "INVALID_PARAM"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_ERROR"
	CodeUnauthorized = "UNAUTHORIZED"
)

// 成功响应构造函数
func SuccessResponse(data interface{}) Response {
	return Response{
		Code:    CodeSuccess,
		Message: "操作成功",
		Data:    data,
	}
}

// 错误响应构造函数
func ErrorResponse(code, message string, details ...interface{}) Response {
	resp := Response{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		resp.Details = details[0]
	}
	return resp
}

// 分页响应
type PaginationResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}
