package shared

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nicklvsa/shorturl/shared/http"
)

func MustParams(c *gin.Context, params ...string) {
	msg := "%s parameter must be specified"

	for _, param := range params {
		if val, found := c.Params.Get(param); !found || len(val) <= 0 {
			msg := fmt.Sprintf(msg, param)
			http.HTTPResponse(
				400,
				false,
				nil,
				&msg,
				c,
			)
		}
	}
}

func ShortenDBKey(shortID string) string {
	return fmt.Sprintf("short::%s", shortID)
}

func ShortenDBVal(employeeID string, data string) string {
	return fmt.Sprintf("%s::%s", employeeID, data)
}
