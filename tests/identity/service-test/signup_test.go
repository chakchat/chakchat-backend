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

// phone not verified
// no such sign up key
// username already exists
// validation failed:
// 		username
//  	name
// success

func Test_SignUp_Success(t *testing.T) {
	signUpKey := requestVerifiedSignUpKey(t)

	resp := doSignUpRequest(t, signUpKey, "John Dafu", "john")
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_SignUp_NotVerified(t *testing.T) {
	phone := common.NewPhone(common.PhoneNotFound)
	signUpKey := requestSignUpSendCode(t, phone)

	resp := doSignUpRequest(t, signUpKey, "John Dafu", "john")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body := common.GetBody(t, resp.Body)
	require.Equal(t, "phone_not_verified", body.ErrorType)
}

func Test_SignUp_UsernameAlreadyExists(t *testing.T) {
	signUpKey := requestVerifiedSignUpKey(t)

	resp := doSignUpRequest(t, signUpKey, "Already exists", "already_exists")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body := common.GetBody(t, resp.Body)
	require.Equal(t, "username_already_exists", body.ErrorType)
}

func Test_SignUp_NoSuchKey(t *testing.T) {
	resp := doSignUpRequest(t, uuid.New(), "John Dafu", "john")
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body := common.GetBody(t, resp.Body)
	require.Equal(t, "signup_key_not_found", body.ErrorType)
}

func Test_SignUp_ValidationFailed(t *testing.T) {
	cases := []struct {
		Username string
		Name     string
	}{
		{
			Username: "1word",
			Name:     "plumbux",
		},
		{
			Username: "_hero",
			Name:     "Heron water",
		},
		{
			Username: "data-saa",
			Name:     "Ok name",
		},
		{
			Username: "abc%62",
			Name:     "Name",
		},
		{
			Username: "it_seems_to_be_too_long_username_idk",
			Name:     "Ordinary name",
		},
	}

	signUpKey := requestVerifiedSignUpKey(t)

	for _, test := range cases {
		resp := doSignUpRequest(t, signUpKey, test.Name, test.Username)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body := common.GetBody(t, resp.Body)
		require.Equal(t, "validation_failed", body.ErrorType)
	}
}

func requestVerifiedSignUpKey(t *testing.T) (signUpKey uuid.UUID) {
	phone := common.NewPhone(common.PhoneNotFound)
	signUpKey = requestSignUpSendCode(t, phone)
	code := getPhoneCode(t, phone)

	resp := doSignUpVerifyRequest(t, signUpKey, code)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	return signUpKey
}

func doSignUpRequest(t *testing.T, signUpKey uuid.UUID, name, username string) *http.Response {
	type Req struct {
		SignUpKey uuid.UUID `json:"signup_key"`
		Username  string    `json:"username"`
		Name      string    `json:"name"`
	}
	reqBody, err := json.Marshal(Req{
		SignUpKey: signUpKey,
		Username:  username,
		Name:      name,
	})
	require.NoError(t, err, "reqBody was: %s", reqBody)

	req, err := http.NewRequest(http.MethodPost, common.ServiceUrl+"/v1.0/signup", bytes.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Add(common.HeaderIdemotencyKey, uuid.NewString())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}
