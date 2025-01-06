package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	serviceUrl = "http://identity-service:5000"

	headerIdemotencyKey = "Idempotency-Key"
)

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

func TestMain(m *testing.M) {
	m.Run()
}

func Test_SendCode(t *testing.T) {
	type TestCase struct {
		Phone      string
		StatusCode int
		ErrorType  string
	}
	cases := []TestCase{
		{
			Phone:      "79000000001",
			StatusCode: http.StatusOK,
		},
		{
			Phone:      "79000000000",
			StatusCode: http.StatusNotFound,
			ErrorType:  "user_not_found",
		},
		{
			Phone:      "it-is-not-phone-number",
			StatusCode: http.StatusBadRequest,
			ErrorType:  "validation_failed",
		},
	}

	type Req struct {
		Phone string `json:"phone"`
	}

	for _, test := range cases {
		reqBody, _ := json.Marshal(Req{
			Phone: test.Phone,
		})

		req, err := http.NewRequest(http.MethodPost, serviceUrl+"/v1.0/signin/send-phone-code", bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Add(headerIdemotencyKey, uuid.NewString())

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		body := getBody(t, resp.Body)

		require.Equal(t, test.StatusCode, resp.StatusCode)
		require.Equal(t, test.ErrorType, body.ErrorType)
	}
}

func Test_SignIn(t *testing.T) {

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
