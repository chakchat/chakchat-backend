package main_test

import (
	"bytes"
	"net/http"
	"test/common"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_Identifies(t *testing.T) {
	symKey := getKey(t, "/app/keys/sym")

	t.Run("Success", func(t *testing.T) {
		accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"typ": "access",
			"jti": uuid.NewString(),
			"iat": time.Now().Add(-1 * time.Minute).Unix(),
			"exp": time.Now().Add(10 * time.Minute).Unix(),
			"iss": "identity_service", // watch identity-service-config.yml
			"aud": []string{"client"},
		})
		accessToken, err := accessJWT.SignedString(symKey)
		require.NoError(t, err)

		resp := doIdentityRequest(t, accessToken)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		internalToken := resp.Header.Get("X-Internal-Token")
		require.NotEmpty(t, internalToken)
		// TODO: check the token
	})

	t.Run("NoAuthHeader", func(t *testing.T) {
		resp, err := http.Get(common.ServiceUrl + "/v1.0/identity")
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		body := common.GetBody(t, resp.Body)
		require.Equal(t, "unauthorized", body.ErrorType)
	})

	t.Run("InvalidTokenType", func(t *testing.T) {
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

		resp := doIdentityRequest(t, refreshToken)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		body := common.GetBody(t, resp.Body)
		require.Equal(t, "invalid_token_type", body.ErrorType)
	})

	t.Run("TokenExpired", func(t *testing.T) {
		accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"typ": "access",
			"jti": uuid.NewString(),
			"iat": time.Now().Add(-10 * time.Minute).Unix(),
			"exp": time.Now().Add(-1 * time.Minute).Unix(),
			"iss": "identity_service", // watch identity-service-config.yml
			"aud": []string{"client"},
		})
		accessToken, err := accessJWT.SignedString(symKey)
		require.NoError(t, err)

		resp := doIdentityRequest(t, accessToken)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		body := common.GetBody(t, resp.Body)
		require.Equal(t, "access_token_expired", body.ErrorType)
	})

	t.Run("InvalidKey", func(t *testing.T) {
		accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"typ": "access",
			"jti": uuid.NewString(),
			"iat": time.Now().Add(-1 * time.Minute).Unix(),
			"exp": time.Now().Add(10 * time.Minute).Unix(),
			"iss": "identity_service", // watch identity-service-config.yml
			"aud": []string{"client"},
		})
		wrongKey := symKey[:len(symKey)-1]
		accessToken, err := accessJWT.SignedString(wrongKey)
		require.NoError(t, err)

		resp := doIdentityRequest(t, accessToken)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		body := common.GetBody(t, resp.Body)
		require.Equal(t, "invalid_token", body.ErrorType)
	})
}

func doIdentityRequest(t *testing.T, accessToken string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, common.ServiceUrl+"/v1.0/identity", bytes.NewReader(nil))
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}
