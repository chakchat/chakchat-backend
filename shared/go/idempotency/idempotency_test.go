package idempotency

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chakchat/chakchat/backend/shared/go/restapi"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TODO: Write comprehensive tests

func TestStoresResponse(t *testing.T) {
	// Arrange
	r, storage := setUp()
	r.POST("/200-with-body", func(c *gin.Context) {
		c.JSON(http.StatusOK, restapi.SuccessResponse{
			Data: gin.H{
				"word": "Success Data",
			},
		})
	})

	const idempotencyKey = "d6f67723-cf79-46a2-9864-ab0d541cd434"

	// Act
	resp := execute(r, "/200-with-body", idempotencyKey)
	captured, ok, err := storage.Get(context.Background(), idempotencyKey)

	// Assert
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusOK, captured.StatusCode)
	assert.JSONEq(t, resp.Body.String(), string(captured.Body))
}

func TestReturnsStored(t *testing.T) {
	// Arrange
	r, storage := setUp()
	r.POST("/200", func(_ *gin.Context) {
		assert.FailNow(t, "This code is not to re-execute")
	})
	const idempotencyKey = "2e89f9fc-5596-4a9c-8177-3b4ce3853b17"
	cachedResp := &CapturedResponse{
		StatusCode: 200,
		Headers: map[string][]string{
			"Custom-Header": {"idk"},
		},
		Body: []byte{69},
	}
	storage.Store(context.Background(), idempotencyKey, cachedResp)

	// Act
	resp := execute(r, "/200", idempotencyKey)

	// Assert
	assert.Equal(t, cachedResp.StatusCode, resp.Code)
	assert.Equal(t, cachedResp.Headers, resp.Header())
	assert.Equal(t, cachedResp.Body, resp.Body.Bytes())
}

func TestSlowExecutionFastRetry(t *testing.T) {
	// Arrange
	r, _ := setUp()
	r.POST("/200-slow", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		c.String(200, "%s", time.Now())
	})
	const idempotencyKey = "2e89f9fc-5596-4a9c-8177-3b4ce3853b17"

	// Act
	wg := sync.WaitGroup{}
	var resp1, resp2, resp3 *httptest.ResponseRecorder
	wg.Add(3)
	go func() {
		resp1 = execute(r, "/200-slow", idempotencyKey)
		wg.Done()
	}()
	go func() {
		resp2 = execute(r, "/200-slow", idempotencyKey)
		wg.Done()
	}()
	go func() {
		resp3 = execute(r, "/200-slow", idempotencyKey)
		wg.Done()
	}()
	wg.Wait()

	// Assert
	assert.Equal(t, resp1, resp2)
	assert.Equal(t, resp1, resp3)
}

func execute(r *gin.Engine, path, idempotencyKey string) *httptest.ResponseRecorder {
	respRecorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(""))
	req.Header[HeaderIdempotencyKey] = []string{idempotencyKey}

	r.ServeHTTP(respRecorder, req)
	return respRecorder
}

func setUp() (*gin.Engine, *mockIdempotencyStorage) {
	r := gin.New()
	mockStorage := &mockIdempotencyStorage{map[string]*CapturedResponse{}}
	r.Use(New(mockStorage))
	return r, mockStorage
}

type mockIdempotencyStorage struct {
	m map[string]*CapturedResponse
}

func (s *mockIdempotencyStorage) Get(_ context.Context, key string) (*CapturedResponse, bool, error) {
	resp, ok := s.m[key]
	return resp, ok, nil
}

func (s *mockIdempotencyStorage) Store(_ context.Context, key string, resp *CapturedResponse) error {
	s.m[key] = resp
	return nil
}
