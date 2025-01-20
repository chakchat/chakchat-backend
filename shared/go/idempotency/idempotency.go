package idempotency

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/chakchat/chakchat-backend/shared/go/internal/restapi"
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
	defer c.Abort()

	key := c.GetHeader(HeaderIdempotencyKey)
	if key == "" {
		errResp := restapi.ErrorResponse{
			ErrorType:    restapi.ErrTypeIdempotencyKeyMissing,
			ErrorMessage: "No \"" + HeaderIdempotencyKey + "\" header provided",
		}
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	m.lock.Lock(key)
	defer m.lock.Unlock(key)

	cached, ok, err := m.storage.Get(c, key)
	if err != nil {
		log.Printf("idempotency middleware: gettings cached response failed: %s", err)
		restapi.SendInternalError(c)
		return
	}
	if ok {
		writeCached(c, cached)
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
		// I think I must guarantee that idempotent endpoint will NOT be re-executed
		// So, some retries are performed if storing response fails
		// Store() looks idempotent so everything is okay
		var err error
		for range 3 {
			// This operation shouldn't be stopped such easily.
			// So, I pass Background() context to prevent cancellation.
			err = m.storage.Store(context.Background(), key, resp)
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Printf("idempotency middleware: all Store() retries failed. Last error: %s", err)
			// I guess no 500 response should be returned because c.Next() succeeded
			// We are just gonna pray that user will not re-execute with same key
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
		log.Printf("idempotency: writing cached body failed: %s", err)
		// I remove headers added from response not to mix them
		for h := range cached.Headers {
			c.Writer.Header().Del(h)
		}
		// TODO: it appends response, not overwrites!
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
