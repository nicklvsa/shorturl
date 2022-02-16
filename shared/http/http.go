package http

import "github.com/gin-gonic/gin"

type HTTPResponseData struct {
	Message *string `json:"message,omitempty"`
	Success bool    `json:"success"`
}

func HTTPResponse(statusCode int, success bool, message *string, context *gin.Context) HTTPResponseData {
	data := HTTPResponseData{
		Message: message,
		Success: success,
	}

	if context != nil {
		context.JSON(statusCode, data)
	}

	return data
}
