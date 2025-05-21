package util

import (
	"encoding/json"
	"net/http"
)

// Response 表示 API 响应
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error 表示 API 错误
type Error struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorResponse 创建一个错误响应
func ErrorResponse(w http.ResponseWriter, code string, message string, statusCode int, details interface{}) {
	response := Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// SuccessResponse 创建一个成功响应
func SuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	response := Response{
		Success: true,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// BadRequestError 返回 400 错误
func BadRequestError(w http.ResponseWriter, message string, details interface{}) {
	ErrorResponse(w, "BAD_REQUEST", message, http.StatusBadRequest, details)
}

// UnauthorizedError 返回 401 错误
func UnauthorizedError(w http.ResponseWriter, message string) {
	ErrorResponse(w, "UNAUTHORIZED", message, http.StatusUnauthorized, nil)
}

// ForbiddenError 返回 403 错误
func ForbiddenError(w http.ResponseWriter, message string) {
	ErrorResponse(w, "FORBIDDEN", message, http.StatusForbidden, nil)
}

// NotFoundError 返回 404 错误
func NotFoundError(w http.ResponseWriter, message string) {
	ErrorResponse(w, "NOT_FOUND", message, http.StatusNotFound, nil)
}

// InternalServerError 返回 500 错误
func InternalServerError(w http.ResponseWriter, message string) {
	ErrorResponse(w, "INTERNAL_SERVER_ERROR", message, http.StatusInternalServerError, nil)
}
