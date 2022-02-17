package tests

import (
	"testing"

	"github.com/nicklvsa/shorturl/routes"
	"github.com/nicklvsa/shorturl/tests/fixtures"
	"github.com/nicklvsa/shorturl/tests/utils"
)

func TestHealthCheckHandler(t *testing.T) {
	ctx, _ := fixtures.GetMockContext()
	routes.HealthcheckHandler(ctx)

	utils.AssertInt(ctx.Writer.Status(), 200, t)
}
