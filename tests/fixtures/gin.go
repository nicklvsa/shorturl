package fixtures

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func GetMockContext() (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	return ctx, recorder
}
