package fixtures

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func GetMockContext() (*gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()
	return gin.CreateTestContext(w)
}
