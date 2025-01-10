package idempotency

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/chakchat/chakchat/backend/shared/go/internal/restapi"
	"github.com/gin-gonic/gin"
)

const HeaderIdempotencyKey = "Idempotency-Key"

type CapturedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type IdempotencyStorage interface {
	Get(ctx context.Context, key string) (*CapturedResponse, bool, error)
	Store(ctx context.Context, key string, resp *CapturedResponse) error
}

func New(storage IdempotencyStorage) gin.HandlerFunc {
	m := &idempotencyMiddleware{
		storage: storage,
		lock:    NewLocker(),
	}
	return m.Handle
}

type idempotencyMiddleware struct {
	storage IdempotencyStorage
	lock    *Locker
}

func (m *idempotencyMiddleware) Handle(c *gin.Context) {
	key := c.GetHeader(HeaderIdempotencyKey)
	if key == "" {
		errResp := restapi.ErrorResponse{
			ErrorType:    restapi.ErrTypeIdempotencyKeyMissing,
			ErrorMessage: "No \"" + HeaderIdempotencyKey + "\" header provided",
		}
		c.JSON(http.StatusBadRequest, errResp)
		c.Abort()
		return
	}

	m.lock.Lock(key)
	defer m.lock.Unlock(key)

	cached, ok, err := m.storage.Get(c, key)
	if err != nil {
		restapi.SendInternalError(c)
		return
	}
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
		restapi.SendInternalError(c)
		return
	}
}

type responseCapturer struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func newResponseCapturer(writer gin.ResponseWriter) responseCapturer {
	return responseCapturer{
		ResponseWriter: writer,
		Body:           &bytes.Buffer{},
	}
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
