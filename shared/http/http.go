package http

import (
	"github.com/gin-gonic/gin"
)

type HTTPResponseData struct {
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func HTTPResponse(statusCode int, success bool, data interface{}, msg *string, context *gin.Context) HTTPResponseData {
	resp := HTTPResponseData{
		Data:    data,
		Message: msg,
		Success: success,
	}

	if context != nil {
		context.JSON(statusCode, resp)
	}

	return resp
}
