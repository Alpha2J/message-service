package controller

import "message-service/internal/app/constant"

type ResponseContent struct {
	Code    int
	Message string
	Data    interface{}
}

func Failed() *ResponseContent {
	return Of(constant.COMMON_BIZ_ERROR.Code, constant.COMMON_BIZ_ERROR.Message, nil)
}

func Success(data interface{}) *ResponseContent {
	return Of(constant.SUCCESS.Code, constant.SUCCESS.Message, data)
}

func Of(code int, message string, data interface{}) *ResponseContent {
	return &ResponseContent{
		code,
		message,
		data,
	}
}
