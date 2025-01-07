package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	m := sync.Map{}

	r.POST("/", func(c *gin.Context) {
		type Req struct {
			Phone   string `json:"phone"`
			Message string `json:"message"`
		}
		req := new(Req)
		if err := c.ShouldBindBodyWithJSON(req); err != nil {
			c.String(http.StatusBadRequest, "it is not valid json")
			return
		}
		m.Store(req.Phone, req.Message)
		c.Status(http.StatusOK)
	})

	r.GET("/:phone", func(c *gin.Context) {
		phone := c.Param("phone")
		if code, ok := m.Load(phone); ok {
			c.String(http.StatusOK, "%s", code)
		} else {
			c.String(http.StatusNotFound, "not found")
		}
	})

	r.Run(":5023")
}
