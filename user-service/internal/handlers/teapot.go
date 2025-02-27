package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AmITeapot() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusTeapot, gin.H{
			"message": "I'm a teapot",
		})
	}
}
