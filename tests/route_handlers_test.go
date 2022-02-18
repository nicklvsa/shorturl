package tests

import (
	"encoding/json"
	"testing"

	"github.com/nicklvsa/shorturl/routes"
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
