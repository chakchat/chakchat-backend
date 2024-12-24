package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStoresResponse(t *testing.T) {
	r, storage := setUp()
	r.POST("/200-with-body", func(c *gin.Context) {
		c.JSON(http.StatusOK, SuccessResponse{
			Data: gin.H{
				"word": "Success Data",
			},
		})
	})

	const idempotencyKey = "d6f67723-cf79-46a2-9864-ab0d541cd434"
	respRecorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/200-with-body", strings.NewReader(""))
	req.Header[HeaderIdempotencyKey] = []string{idempotencyKey}

	r.ServeHTTP(respRecorder, req)

	captured, ok := storage.Get(idempotencyKey)
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, captured.StatusCode)
	assert.JSONEq(t, respRecorder.Body.String(), string(captured.Body))
}

func TestReturnsStored(t *testing.T) {
	r, storage := setUp()
	r.Use(gin.Logger())
	r.POST("/200", func(ctx *gin.Context) {
		assert.FailNow(t, "This code is not to re-execute")
	})
	const idempotencyKey = "2e89f9fc-5596-4a9c-8177-3b4ce3853b17"
	resp := &CapturedResponse{
		StatusCode: 200,
		Headers: map[string][]string{
			"Custom-Header": {"idk"},
		},
		Body: []byte{69},
	}
	storage.Store(idempotencyKey, resp)

	respRecorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/200", strings.NewReader(""))
	req.Header[HeaderIdempotencyKey] = []string{idempotencyKey}

	r.ServeHTTP(respRecorder, req)

	assert.Equal(t, resp.StatusCode, respRecorder.Code)
	assert.Equal(t, resp.Headers, respRecorder.Header())
	assert.Equal(t, resp.Body, respRecorder.Body.Bytes())
}

func setUp() (*gin.Engine, *mockIdempotencyStorage) {
	r := gin.New()
	mockStorage := &mockIdempotencyStorage{map[string]*CapturedResponse{}}
	r.Use(CheckIdempotencyKey(mockStorage))
	return r, mockStorage
}

type mockIdempotencyStorage struct {
	m map[string]*CapturedResponse
}

func (s *mockIdempotencyStorage) Get(key string) (*CapturedResponse, bool) {
	resp, ok := s.m[key]
	return resp, ok
}

func (s *mockIdempotencyStorage) Store(key string, resp *CapturedResponse) error {
	s.m[key] = resp
	return nil
}
