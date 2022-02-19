package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/shared/http"
)

func HealthcheckHandler(c *gin.Context) {
	msg := "healthcheck"
	http.HTTPResponse(200, true, &msg, &msg, c)
}
