package routes

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/actions"
	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/shared/errs"
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
		logger.Errorf(err.Error())

		msg := fmt.Sprintf("%s not found", shortID)
		http.HTTPResponse(
			404,
			false,
			nil,
			&msg,
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
		msg := errs.URLMustBeProvidedAPIError.Str()
		http.HTTPResponse(
			400,
			false,
			nil,
			&msg,
			c,
		)
		return
	}

	if len(expires_in) > 0 {
		// we should expire the url in x minutes
		mins, err := strconv.Atoi(expires_in)
		if err != nil {
			logger.Errorf(err.Error())

			msg := errs.FormatMismatchExpiresInAPIError.Str()
			http.HTTPResponse(
				400,
				false,
				nil,
				&msg,
				c,
			)
			return
		}

		expireMins = &mins
	}

	shortID, err := h.Actions.CreateURLMapping(longURL, employeeID, expireMins)
	if err != nil {
		msg := errs.SaveURLMappingFailedAPIError.Str()
		http.HTTPResponse(
			400,
			false,
			nil,
			&msg,
			c,
		)
		return
	}

	data := fmt.Sprintf("http://localhost:8080/v/%s", shortID)
	http.HTTPResponse(
		201,
		true,
		&data,
		nil,
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
			nil,
			&msg,
			c,
		)
		return
	}

	http.HTTPResponse(
		200,
		true,
		metrics,
		nil,
		c,
	)
	return
}

func (h ShortURLHandler) DeleteShortURLHandler(c *gin.Context) {
	shared.MustParams(c, "employee_id", "id")

	shortID := c.Param("id")
	employeeID := c.Param("employee_id")

	if err := h.Actions.DeleteShortURL(shortID, employeeID); err != nil {
		logger.Errorf(err.Error())

		statusCode := 400
		msg := errs.DeleteURLFailedAPIError.Str()

		if err == errs.UnauthorizedAPIError.Err() {
			statusCode = 401
		}

		http.HTTPResponse(
			statusCode,
			false,
			nil,
			&msg,
			c,
		)
		return
	}

	msg := fmt.Sprintf("%s has been deleted", shortID)
	http.HTTPResponse(
		200,
		true,
		nil,
		&msg,
		c,
	)
	return
}
