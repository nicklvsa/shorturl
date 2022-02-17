package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/shared"
)

type ShortURLHandler struct {
	Config *shared.Config
}

func NewShortURLHandler(cfg *shared.Config) *ShortURLHandler {
	return &ShortURLHandler{
		Config: cfg,
	}
}

func (h ShortURLHandler) CreateShortURLHandler(c *gin.Context) {

}

func (h ShortURLHandler) GetShortURLMetricsHandler(c *gin.Context) {

}

func (h ShortURLHandler) DeleteShortURLHandler(c *gin.Context) {

}
