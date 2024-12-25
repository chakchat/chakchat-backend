package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	storage := &MockIdempotencyStorage{
		m: map[string]*CapturedResponse{},
	}
	r.Use(NewIdempotencyMiddleware(storage).Handle)

	r.GET("/hello", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
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

func (s *MockIdempotencyStorage) Get(_ context.Context, key string) (*CapturedResponse, bool) {
	resp, ok := s.m[key]
	delete(s.m, key)
	return resp, ok
}

func (s *MockIdempotencyStorage) Store(_ context.Context, key string, resp *CapturedResponse) error {
	s.m[key] = resp
	return nil
}
