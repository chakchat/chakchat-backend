package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func InternalError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		status := c.Writer.Status()
		if status >= 500 {
			c.Error(errors.New("internal error middleware: got status code >=500"))
		}
	}
}
