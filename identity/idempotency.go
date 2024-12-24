package main

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

const HeaderIdempotencyKey = "Idempotency-Key"

type CapturedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type IdempotencyStorage interface {
	Get(key string) (*CapturedResponse, bool)
	Store(key string, resp *CapturedResponse) error
}

func captureCondition(resp *CapturedResponse) bool {
	return resp.StatusCode < 500
}

func CheckIdempotencyKey(storage IdempotencyStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader(HeaderIdempotencyKey)
		cached, ok := storage.Get(key)
		if ok {
			writeCached(c, cached)
			c.Abort()
			return
		}

		capturer := newResponseCapturer(c.Writer)
		c.Writer = capturer

		c.Next()

		if resp := capturer.ExtractResponse(); captureCondition(resp) {
			err := storage.Store(key, resp)
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
}

func copyHeaders(src http.Header, dst http.Header) {
	for h, val := range src {
		dst[h] = val
	}
}

func writeCached(c *gin.Context, cached *CapturedResponse) {
	_, err := c.Writer.Write(cached.Body)
	if err != nil {
		// I don't know what case this error mean
		errResp := ErrorResponse{
			ErrorType:    ErrorTypeInternal,
			ErrorMessage: "Internal Server Error",
		}
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	copyHeaders(cached.Headers, c.Writer.Header())
	c.Status(cached.StatusCode)
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
