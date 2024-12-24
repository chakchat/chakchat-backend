package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	storage := &MockIdempotencyStorage{
		m: map[string]*CapturedResponse{},
	}
	r.Use(CheckIdempotencyKey(storage))

	r.GET("/hello", func(c *gin.Context) {
		c.Header("idk", "idk")
		c.JSON(200, SuccessResponse{
			Data: struct {
				Time time.Time `json:"time"`
			}{
				Time: time.Now(),
			},
		})
	})

	r.Run()
}

type MockIdempotencyStorage struct {
	m map[string]*CapturedResponse
}

func (s *MockIdempotencyStorage) Get(key string) (*CapturedResponse, bool) {
	resp, ok := s.m[key]
	return resp, ok
}

func (s *MockIdempotencyStorage) Store(key string, resp *CapturedResponse) error {
	s.m[key] = resp
	return nil
}
