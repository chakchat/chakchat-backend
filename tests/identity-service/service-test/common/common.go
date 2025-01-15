package common

import (
	"encoding/json"
	"io"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	ServiceUrl = "http://identity-service:5000"

	HeaderIdemotencyKey = "Idempotency-Key"
)

var incNumber atomic.Uint64

const (
	PhoneExisting = "1"
	PhoneErroring = "2"
	PhoneNotFound = "0"
)

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type StdResp struct {
	ErrorType    string        `json:"error_type,omitempty"`
	ErrorMessage string        `json:"error_message,omitempty"`
	ErrorDetails []ErrorDetail `json:"error_details,omitempty"`

	Data json.RawMessage `json:"data,omitempty"`
}

func NewPhone(phoneState string) string {
	templ := "7900000000*"
	n := incNumber.Add(1)
	nStr := strconv.FormatUint(n, 10)
	return templ[:10-len(nStr)] + nStr + phoneState
}

func MatchUserId(phone string) string {
	id := uuid.Nil.String()
	id = id[:len(id)-11] + phone
	return id
}

func GetBody(t *testing.T, body io.ReadCloser) *StdResp {
	defer body.Close()
	raw, err := io.ReadAll(body)
	require.NoError(t, err)

	// log.Printf("body: \"%s\"", raw)
	resp := new(StdResp)
	err = json.Unmarshal(raw, resp)
	require.NoError(t, err)

	return resp
}
