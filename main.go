package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/routes"
	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/shared/logger"
)

func main() {
	// setup main routing engine
	router := gin.Default()

	// initialize the global config struct
	config, err := shared.InitConfig()
	if err != nil {
		logger.Panicf(err.Error())
	}

	// initialize the short url routes
	shortHandler := routes.NewShortURLHandler(config)

	// healthcheck route for docker
	router.GET("/healthcheck", routes.HealthcheckHandler)

	// create a new short url based on a long url input
	router.GET("/shorten_url", shortHandler.CreateShortURLHandler)

	// delete a short url by its id
	router.GET("/delete_url", shortHandler.DeleteShortURLHandler)

	// retrieve metrics for a specific short url
	router.GET("/metrics", shortHandler.GetShortURLMetricsHandler)

	// start http server on port 8080
	if err := router.Run(":8080"); err != nil {
		logger.Panicf(err.Error())
	}
}
