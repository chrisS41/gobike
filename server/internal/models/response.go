package models

import "github.com/chrisS41/gobike-server/internal/errors"

// SuccessResponse 성공 응답을 위한 구조체
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// NewSuccessResponse 성공 응답 생성
func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Code:    errors.SUCCESS,
		Data:    data,
		Message: errors.GetErrorMessage(errors.SUCCESS),
	}
}

// NewErrorResponse 에러 응답 생성
func NewErrorResponse(errorCode int) *Response {
	return &Response{
		Code:    errorCode,
		Data:    nil,
		Message: errors.GetErrorMessage(errorCode),
	}
}

// NewErrorResponseWithMessage 에러 응답 생성
func NewErrorResponseWithMessage(errorCode int, message string) *Response {
	return &Response{
		Code:    errorCode,
		Data:    nil,
		Message: message,
	}
}
