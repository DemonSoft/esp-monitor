package rest

import (
	"fmt"
	"net/http"
	"remoteesp/internal/domain/env"

	"github.com/gin-gonic/gin"
)

const (
	XTOKEN = "xToken"
)

func (r *Rest) MiddlewareHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				PanicError[Empty](c, fmt.Sprintf("internal server error: %v", r))
			}
		}()

		extractedToken := c.GetHeader(XTOKEN)
		if extractedToken != env.XToken {
			c.IndentedJSON(http.StatusUnauthorized, restError[Empty]("Check your permissions", "-1"))
			c.Abort()
			return
		}

		c.Next()

		// If you see some errors (for example, c.AbortWithStatusJSON was called)
		status := c.Writer.Status()
		if status >= 400 && !c.IsAborted() {
			SomeError[Empty](c, status)
		}
	}
}
