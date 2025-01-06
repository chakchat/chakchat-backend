package send_code_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	serviceUrl = "http://identity-service:5000"

	headerIdemotencyKey = "Idempotency-Key"
)

var incNumber atomic.Uint64

const (
	phoneExisting = "1"
	phoneErroring = "2"
	phoneNotFound = "0"
)

func Test_SendCode(t *testing.T) {
	type TestCase struct {
		Name       string
		Phone      string
		StatusCode int
		ErrorType  string
	}
	cases := []TestCase{
		{
			Name:       "Success",
			Phone:      newPhone(phoneExisting),
			StatusCode: http.StatusOK,
		},
		{
			Name:       "UserNotFound",
			Phone:      newPhone(phoneNotFound),
			StatusCode: http.StatusNotFound,
			ErrorType:  "user_not_found",
		},
		{
			Name:       "InvalidPhone",
			Phone:      "it-is-not-phone-number",
			StatusCode: http.StatusBadRequest,
			ErrorType:  "validation_failed",
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			resp := sendRequest(t, test.Phone, uuid.New())
			body := getBody(t, resp.Body)

			require.Equal(t, test.StatusCode, resp.StatusCode)
			require.Equal(t, test.ErrorType, body.ErrorType)
		})
	}
}

func Test_SendCode_FreqExceeded(t *testing.T) {
	phone := newPhone(phoneExisting)

	resp1 := sendRequest(t, phone, uuid.New())
	require.Equal(t, http.StatusOK, resp1.StatusCode)

	resp2 := sendRequest(t, phone, uuid.New())
	require.Equal(t, http.StatusBadRequest, resp2.StatusCode)

	body := getBody(t, resp2.Body)
	require.Equal(t, "send_code_freq_exceeded", body.ErrorType)
}

type Request struct {
	Phone string `json:"phone"`
}

func newPhone(phoneState string) string {
	templ := "7900000000*"
	n := incNumber.Add(1)
	nStr := strconv.FormatUint(n, 10)
	return templ[:10-len(nStr)] + nStr + phoneState
}

func sendRequest(t *testing.T, phone string, idempotencyKey uuid.UUID) *http.Response {
	reqBody, _ := json.Marshal(Request{
		Phone: phone,
	})

	req, err := http.NewRequest(http.MethodPost, serviceUrl+"/v1.0/signin/send-phone-code", bytes.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Add(headerIdemotencyKey, idempotencyKey.String())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type StdResp struct {
	ErrorType    string        `json:"error_type,omitempty"`
	ErrorMessage string        `json:"error_message,omitempty"`
	ErrorDetails []ErrorDetail `json:"error_details,omitempty"`

	Data any `json:"data,omitempty"`
}

func getBody(t *testing.T, body io.ReadCloser) *StdResp {
	defer body.Close()
	raw, err := io.ReadAll(body)
	require.NoError(t, err)

	// log.Printf("body: \"%s\"", raw)
	resp := new(StdResp)
	err = json.Unmarshal(raw, resp)
	require.NoError(t, err)

	return resp
}
