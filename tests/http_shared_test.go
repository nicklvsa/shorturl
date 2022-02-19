package tests

import (
	"encoding/json"
	"testing"

	"github.com/nicklvsa/shorturl/shared/http"
	"github.com/nicklvsa/shorturl/tests/fixtures"
	"github.com/nicklvsa/shorturl/tests/utils"
)

func TestHTTPResponse(t *testing.T) {
	ctx, data := fixtures.GetMockContext()

	message := "error"
	http.HTTPResponse(401, false, nil, &message, ctx)

	out := make(map[string]interface{})
	json.NewDecoder(data.Body).Decode(&out)

	utils.AssertFalse(out["success"].(bool), t)
	utils.AssertInt(ctx.Writer.Status(), 401, t)
	utils.AssertStr(out["message"].(string), "error", t)
}
