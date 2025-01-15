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

// wrong sign up key
// wrong code
// success

func Test_SignUpVerify_NoSuchKey(t *testing.T) {
	resp := doSignUpVerifyRequest(t, uuid.New(), "123456")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body := common.GetBody(t, resp.Body)
	require.Equal(t, "signup_key_not_found", body.ErrorType)
}

func Test_SignUpVerify_WrongCode(t *testing.T) {
	phone := common.NewPhone(common.PhoneNotFound)

	signUpKey := requestSignUpSendCode(t, phone)

	resp := doSignUpVerifyRequest(t, signUpKey, "696969") // I believe it won't accidentally match true code
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body := common.GetBody(t, resp.Body)
	require.Equal(t, "wrong_code", body.ErrorType)
}

func Test_SignUpVerify_Success(t *testing.T) {
	phone := common.NewPhone(common.PhoneNotFound)
	signUpKey := requestSignUpSendCode(t, phone)
	code := getPhoneCode(t, phone)

	resp := doSignUpVerifyRequest(t, signUpKey, code)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func requestSignUpSendCode(t *testing.T, phone string) (signUpKey uuid.UUID) {
	resp := doSignUpSendCodeRequest(t, phone, uuid.New())
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := common.GetBody(t, resp.Body)

	type RespData struct {
		SignUpKey uuid.UUID `json:"signup_key"`
	}
	var data RespData
	err := json.Unmarshal(body.Data, &data)
	require.NoError(t, err)

	return data.SignUpKey
}

func doSignUpVerifyRequest(t *testing.T, signUpKey uuid.UUID, code string) *http.Response {
	type Req struct {
		SignUpKey uuid.UUID `json:"signup_key"`
		Code      string    `json:"code"`
	}
	reqBody, _ := json.Marshal(Req{
		SignUpKey: signUpKey,
		Code:      code,
	})

	req, err := http.NewRequest(http.MethodPost, common.ServiceUrl+"/v1.0/signup/verify-code", bytes.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Add(common.HeaderIdemotencyKey, uuid.NewString())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}
