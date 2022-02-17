package routes

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/actions"
	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/shared/http"
	"github.com/nicklvsa/shorturl/shared/logger"
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
		return
	}

	if err := h.Actions.IncrShortURLCount(shortID); err != nil {
		logger.Warnf("unable to increment url view count. Error: ", err.Error())
	}

	c.Redirect(301, longURL)
	return
}

func (h ShortURLHandler) CreateShortURLHandler(c *gin.Context) {
	shared.MustParams(c, "employee_id")

	var expireMins *int
	employeeID := c.Param("employee_id")

	longURL := c.Query("url")
	expires_in := c.Query("expires")

	if len(longURL) <= 0 {
		msg := "a url must be set to be shortened"
		http.HTTPResponse(
			400,
			false,
			shared.GetPointerToString(msg),
			c,
		)
		return
	}

	if len(expires_in) > 0 {
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
			return
		}

		expireMins = &mins
	}

	shortID, err := h.Actions.CreateURLMapping(longURL, employeeID, expireMins)
	if err != nil {
		msg := "unable to save url mapping, try again later"
		http.HTTPResponse(
			400,
			false,
			shared.GetPointerToString(msg),
			c,
		)
		return
	}

	msg := fmt.Sprintf("http://localhost:8080/v/%s", shortID)
	http.HTTPResponse(
		201,
		true,
		shared.GetPointerToString(msg),
		c,
	)
	return
}

func (h ShortURLHandler) GetShortURLMetricsHandler(c *gin.Context) {
	shared.MustParams(c, "id", "employee_id")

	shortID := c.Param("id")
	employeeID := c.Param("employee_id")

	metrics, err := h.Actions.GetShortURLMetrics(shortID, employeeID)
	if err != nil {
		logger.Errorf(err.Error())
		msg := "unable to retrieve metrics"
		http.HTTPResponse(
			400,
			false,
			shared.GetPointerToString(msg),
			c,
		)
		return
	}

	c.JSON(200, metrics)
	return
}

func (h ShortURLHandler) DeleteShortURLHandler(c *gin.Context) {
	shared.MustParams(c, "employee_id", "id")

	shortID := c.Param("id")
	employeeID := c.Param("employee_id")

	if err := h.Actions.DeleteShortURL(shortID, employeeID); err != nil {
		msg := "could not delete provided short url"
		statusCode := 400

		if err.Error() == "unauthorized" {
			statusCode = 401
		}

		http.HTTPResponse(
			statusCode,
			false,
			shared.GetPointerToString(msg),
			c,
		)
		return
	}

	msg := fmt.Sprintf("%s has been deleted", shortID)
	http.HTTPResponse(
		200,
		true,
		shared.GetPointerToString(msg),
		c,
	)
	return
}
