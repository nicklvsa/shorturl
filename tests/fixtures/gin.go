package fixtures

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

func GetMockContext() (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Request = &http.Request{
		URL: &url.URL{},
	}

	return ctx, recorder
}
