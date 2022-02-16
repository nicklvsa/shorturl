package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/shared/http"
)

func HealthcheckHandler(c *gin.Context) {
	http.HTTPResponse(200, true, shared.GetPointerToString("healthcheck"), c)
}
