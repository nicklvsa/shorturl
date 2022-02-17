package routes

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/actions"
	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/shared/http"
)

type ShortURLHandler struct {
	Config  *shared.Config
	Actions *actions.Actions
}

func NewShortURLHandler(cfg *shared.Config) *ShortURLHandler {
	return &ShortURLHandler{
		Config:  cfg,
		Actions: actions.NewActions(cfg),
	}
}

func (h ShortURLHandler) VisitShortURL(c *gin.Context) {
	shared.MustParams(c, "id")

	shortID := c.Param("id")

	longURL, err := h.Actions.GetLongURL(shortID)
	if err != nil {
		msg := fmt.Sprintf("%s not found", shortID)
		http.HTTPResponse(
			404,
			false,
			shared.GetPointerToString(msg),
			c,
		)
	}

	http.HTTPResponse(
		200,
		true,
		shared.GetPointerToString(longURL),
		c,
	)
}

func (h ShortURLHandler) CreateShortURLHandler(c *gin.Context) {
	shared.MustParams(c, "employee_id", "url")

	var expireMins *int
	longURL := c.Param("url")
	employeeID := c.Param("employee_id")

	expires_in := c.Query("expires")
	if len(expires_in) >= 0 {
		// we should expire the url in x minutes
		mins, err := strconv.Atoi(expires_in)
		if err != nil {
			msg := "expires_in must be a number in minutes"
			http.HTTPResponse(
				400,
				false,
				shared.GetPointerToString(msg),
				c,
			)
		}

		expireMins = &mins
	}

	shortID, err := h.Actions.CreateURLMapping(longURL, employeeID, expireMins)
	if err != nil {
		if err != nil {
			msg := "unable to save url mapping, try again later"
			http.HTTPResponse(
				400,
				false,
				shared.GetPointerToString(msg),
				c,
			)
		}
	}

	msg := fmt.Sprintf("http://localhost:8080/v/%s", shortID)
	http.HTTPResponse(
		201,
		true,
		shared.GetPointerToString(msg),
		c,
	)
}

func (h ShortURLHandler) GetShortURLMetricsHandler(c *gin.Context) {
	// shared.MustParams(c, "id")

	// shortID := c.Param("id")
}

func (h ShortURLHandler) DeleteShortURLHandler(c *gin.Context) {
	shared.MustParams(c, "employee_id", "id")

	// shortID := c.Param("id")
	// employeeID := c.Param("employee_id")
}
