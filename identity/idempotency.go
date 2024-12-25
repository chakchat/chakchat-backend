package main

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const HeaderIdempotencyKey = "Idempotency-Key"

type CapturedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type IdempotencyStorage interface {
	Get(ctx context.Context, key string) (*CapturedResponse, bool)
	Store(ctx context.Context, key string, resp *CapturedResponse) error
}

type IdempotencyMiddleware struct {
	storage IdempotencyStorage
	lock    *Locker
}

func NewIdempotencyMiddleware(storage IdempotencyStorage) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		storage: storage,
		lock:    NewLocker(),
	}
}

func (m *IdempotencyMiddleware) Handle(c *gin.Context) {
	key := c.GetHeader(HeaderIdempotencyKey)
	if key == "" {
		errResp := ErrorResponse{
			ErrorType:    "idempotency_key_missing",
			ErrorMessage: "No \"" + HeaderIdempotencyKey + "\" header provided",
		}
		c.JSON(http.StatusBadRequest, errResp)
		c.Abort()
		return
	}

	m.lock.Lock(key)
	defer m.lock.Unlock(key)

	cached, ok := m.storage.Get(c, key)
	if ok {
		writeCached(c, cached)
		c.Abort()
		return
	}
	// Check if storage (and this func too) was cancelled
	if err := c.Err(); err != nil {
		return
	}

	capturer := newResponseCapturer(c.Writer)
	c.Writer = capturer

	c.Next()

	if resp := capturer.ExtractResponse(); captureCondition(resp) {
		err := m.storage.Store(context.Background(), key, resp) // This operation shouldn't be stopped such easily
		if err != nil {
			// TODO: what to do then
			// I think I must guarantee that idempotent endpoint will NOT be re-executed
			// But in this scenario this concept goes wrong
			// Maybe `storage.Store()` retry?
			// I guess no, especially if Store() is deteministic
			// But what do I do?
		}
	}
}

func captureCondition(resp *CapturedResponse) bool {
	return resp.StatusCode < 500
}

func copyHeaders(src http.Header, dst http.Header) {
	for h, val := range src {
		dst.Set(h, strings.Join(val, "; "))
	}
}

func writeCached(c *gin.Context, cached *CapturedResponse) {
	copyHeaders(cached.Headers, c.Writer.Header())
	c.Status(cached.StatusCode)
	_, err := c.Writer.Write(cached.Body)
	if err != nil {
		// I don't know what case this error mean|

		// I remove headers added from response not to mix them
		for h := range cached.Headers {
			c.Writer.Header().Del(h)
		}
		writeInternalError(c)
		return
	}
}

func writeInternalError(c *gin.Context) {
	errResp := ErrorResponse{
		ErrorType:    ErrorTypeInternal,
		ErrorMessage: "Internal Server Error",
	}
	c.JSON(http.StatusInternalServerError, errResp)
}

func newResponseCapturer(writer gin.ResponseWriter) responseCapturer {
	return responseCapturer{
		ResponseWriter: writer,
		Body:           &bytes.Buffer{},
	}
}

type responseCapturer struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (c responseCapturer) Write(data []byte) (int, error) {
	_, _ = c.Body.Write(data) // Never returns error
	return c.ResponseWriter.Write(data)
}

func (c responseCapturer) ExtractResponse() *CapturedResponse {
	resp := &CapturedResponse{
		StatusCode: c.Status(),
		Headers:    make(http.Header, len(c.Header())),
		Body:       c.Body.Bytes(),
	}
	copyHeaders(c.Header(), resp.Headers)
	return resp
}
