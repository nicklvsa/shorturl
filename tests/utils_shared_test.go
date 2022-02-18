package tests

import (
	"encoding/json"
	"testing"

	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/tests/fixtures"
	"github.com/nicklvsa/shorturl/tests/utils"
)

func TestMustParams(t *testing.T) {
	ctx, data := fixtures.GetMockContext()
	shared.MustParams(ctx, "bad_param")

	out := make(map[string]interface{})
	json.NewDecoder(data.Body).Decode(&out)

	utils.AssertFalse(out["success"].(bool), t)
	utils.AssertStr(out["message"].(string), "bad_param parameter must be specified", t)
}

func TestShortenDBKey(t *testing.T) {
	key := shared.ShortenDBKey("cool_id")

	utils.AssertStr(key, "short::cool_id", t)
}

func TestShortenDBVal(t *testing.T) {
	key := shared.ShortenDBVal("123", "cool_id")

	utils.AssertStr(key, "123::cool_id", t)
}
