package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/shared/http"
)

// healthcheck handler exists so the container can run healthchecks
// against our server
func HealthcheckHandler(c *gin.Context) {
	msg := "healthcheck"
	http.HTTPResponse(200, true, &msg, &msg, c)
}
