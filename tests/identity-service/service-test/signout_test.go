package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"test/common"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_SignsOut(t *testing.T) {
	symKey := getKey(t, "/app/keys/sym")

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"typ": "refresh",
		"jti": uuid.NewString(),
		"iat": time.Now().Add(-1 * time.Minute).Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		"iss": "identity_service", // watch identity-service-config.yml
		"aud": []string{"client"},
	})
	refreshToken, err := refreshJWT.SignedString(symKey)
	require.NoError(t, err)

	signOutResp := doSignOutRequest(t, refreshToken)
	require.Equal(t, http.StatusOK, signOutResp.StatusCode)
}

func doSignOutRequest(t *testing.T, refreshToken string) *http.Response {
	type Request struct {
		RefreshJWT string `json:"refresh_token"`
	}
	reqBody, _ := json.Marshal(Request{
		RefreshJWT: refreshToken,
	})
	req, err := http.NewRequest(http.MethodPut, common.ServiceUrl+"/v1.0/sign-out", bytes.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Add(common.HeaderIdemotencyKey, uuid.NewString())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}
