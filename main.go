package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/routes"
)

func main() {
	// setup main routing engine
	router := gin.Default()

	// GET routes
	router.GET("/healthcheck", routes.HealthcheckHandler)

	// start http server on port 8080
	router.Run(":8080")
}
