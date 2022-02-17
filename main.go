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

	// base route for redirecting to a short url's long url mapping
	router.GET("/v/:id", shortHandler.VisitShortURL)

	// create a router group to nest all shortener routes under the /short prefix
	shortGroup := router.Group("/short")

	// create a new short url based on a long url input
	shortGroup.GET("/new/:employee_id", shortHandler.CreateShortURLHandler)

	// delete a short url by its id
	shortGroup.GET("/delete/:employee_id/:id", shortHandler.DeleteShortURLHandler)

	// retrieve metrics for a specific short url
	shortGroup.GET("/metrics/:employee_id/:id", shortHandler.GetShortURLMetricsHandler)

	// start http server on port 8080
	if err := router.Run(":8080"); err != nil {
		logger.Panicf(err.Error())
	}
}
