package tests

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v8"
	"github.com/nicklvsa/shorturl/routes"
	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/tests/fixtures"
	"github.com/nicklvsa/shorturl/tests/utils"
)

func TestHealthCheckHandler(t *testing.T) {
	ctx, data := fixtures.GetMockContext()
	routes.HealthcheckHandler(ctx)

	out := make(map[string]interface{})
	json.NewDecoder(data.Body).Decode(&out)

	utils.AssertTrue(out["success"].(bool), t)
	utils.AssertInt(ctx.Writer.Status(), 200, t)
	utils.AssertStr(out["message"].(string), "healthcheck", t)
}

func TestVisitShortURL(t *testing.T) {
	ctx, _ := fixtures.GetMockContext()
	client, mock := redismock.NewClientMock()

	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "id",
		Value: "someshortid",
	})

	cfg := &shared.Config{
		DB:            client,
		MetricsConfig: &shared.MetricsConfig{},
	}

	mock.ExpectGet("short::someshortid").SetVal("abc::example.com")

	handler := routes.NewShortURLHandler(cfg)
	handler.VisitShortURL(ctx)

	utils.AssertInt(ctx.Writer.Status(), 301, t)
}

func TestCreateShortURLHandler(t *testing.T) {
	ctx, _ := fixtures.GetMockContext()
	client, mock := redismock.NewClientMock()

	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "employee_id",
		Value: "abc",
	})

	parsedURL, err := url.Parse("http://localhost:8080?url=http://example.com")
	if err != nil {
		t.Fatal(err)
	}

	ctx.Request.URL = parsedURL

	cfg := &shared.Config{
		DB:            client,
		MetricsConfig: &shared.MetricsConfig{},
	}

	mock.Regexp().ExpectSet("(short::*)\\w+", "abc::http://example.com", 0).SetVal("1")

	handler := routes.NewShortURLHandler(cfg)
	handler.CreateShortURLHandler(ctx)

	utils.AssertInt(ctx.Writer.Status(), 201, t)
}

func TestGetShortURLMetricsHandler(t *testing.T) {
	ctx, _ := fixtures.GetMockContext()
	client, mock := redismock.NewClientMock()

	ctx.Params = append(ctx.Params,
		[]gin.Param{
			{
				Key:   "id",
				Value: "someshortid",
			},
			{
				Key:   "employee_id",
				Value: "abc",
			},
		}...,
	)

	cfg := &shared.Config{
		DB: client,
		MetricsConfig: &shared.MetricsConfig{
			Periods: map[string]string{
				"12h": "helloworld",
			},
		},
	}

	mock.ExpectGet("short::someshortid").SetVal("abc::abc")

	handler := routes.NewShortURLHandler(cfg)
	handler.GetShortURLMetricsHandler(ctx)

	utils.AssertInt(ctx.Writer.Status(), 200, t)
}

func TestDeleteShortURLHandler(t *testing.T) {
	ctx, _ := fixtures.GetMockContext()
	client, mock := redismock.NewClientMock()

	ctx.Params = append(ctx.Params,
		[]gin.Param{
			{
				Key:   "id",
				Value: "someshortid",
			},
			{
				Key:   "employee_id",
				Value: "abc",
			},
		}...,
	)

	cfg := &shared.Config{
		DB:            client,
		MetricsConfig: &shared.MetricsConfig{},
	}

	mock.ExpectGet("short::someshortid").SetVal("abc::abc")
	mock.ExpectDel("short::someshortid").SetVal(1)

	handler := routes.NewShortURLHandler(cfg)
	handler.DeleteShortURLHandler(ctx)

	utils.AssertInt(ctx.Writer.Status(), 200, t)
}
