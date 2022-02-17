package http

import (
	"github.com/gin-gonic/gin"
)

type HTTPResponseData struct {
	Data    *string `json:"message,omitempty"`
	Success bool    `json:"success"`
}

func HTTPResponse(statusCode int, success bool, data *string, context *gin.Context) HTTPResponseData {
	resp := HTTPResponseData{
		Data:    data,
		Success: success,
	}

	if context != nil {
		context.JSON(statusCode, resp)
	}

	return resp
}
