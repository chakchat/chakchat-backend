package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"test/common"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_SignUpSendCode(t *testing.T) {
	type TestCase struct {
		Name       string
		Phone      string
		StatusCode int
		ErrorType  string
	}
	cases := []TestCase{
		{
			Name:       "Success",
			Phone:      common.NewPhone(common.PhoneNotFound),
			StatusCode: http.StatusOK,
		},
		{
			Name:       "UserExists",
			Phone:      common.NewPhone(common.PhoneExisting),
			StatusCode: http.StatusBadRequest,
			ErrorType:  "user_already_exists",
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
			resp := doSignUpSendCodeRequest(t, test.Phone, uuid.New())
			body := common.GetBody(t, resp.Body)

			require.Equal(t, test.StatusCode, resp.StatusCode)
			require.Equal(t, test.ErrorType, body.ErrorType)
		})
	}
}

func Test_SignUpSendCode_FreqExceeded(t *testing.T) {
	phone := common.NewPhone(common.PhoneNotFound)

	resp1 := doSignUpSendCodeRequest(t, phone, uuid.New())
	require.Equal(t, http.StatusOK, resp1.StatusCode)

	resp2 := doSignUpSendCodeRequest(t, phone, uuid.New())
	require.Equal(t, http.StatusBadRequest, resp2.StatusCode)

	body := common.GetBody(t, resp2.Body)
	require.Equal(t, "send_code_freq_exceeded", body.ErrorType)
}

func doSignUpSendCodeRequest(t *testing.T, phone string, idempotencyKey uuid.UUID) *http.Response {
	type Req struct {
		Phone string `json:"phone"`
	}
	reqBody, _ := json.Marshal(Req{
		Phone: phone,
	})

	req, err := http.NewRequest(http.MethodPost, common.ServiceUrl+"/v1.0/signup/send-phone-code", bytes.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Add(common.HeaderIdemotencyKey, idempotencyKey.String())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}
