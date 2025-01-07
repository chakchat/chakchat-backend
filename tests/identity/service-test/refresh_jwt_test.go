package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"test/common"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_RefreshJWT(t *testing.T) {
	symKey := getKey(t, "/app/keys/sym")

	type Response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	t.Run("Refreshes", func(t *testing.T) {
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

		resp := doRefreshJWTRequest(t, refreshToken)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body := common.GetBody(t, resp.Body)
		require.Empty(t, body.ErrorType)

		var bodyResp Response
		err = json.Unmarshal(body.Data, &bodyResp)
		require.NoError(t, err)
	})

	t.Run("DeniesAccessToken", func(t *testing.T) {
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

		resp := doRefreshJWTRequest(t, accessToken)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body := common.GetBody(t, resp.Body)
		require.Equal(t, "invalid_token_type", body.ErrorType)
	})

	t.Run("DeniesWrongKey", func(t *testing.T) {
		refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"typ": "refresh",
			"jti": uuid.NewString(),
			"iat": time.Now().Add(-1 * time.Minute).Unix(),
			"exp": time.Now().Add(10 * time.Minute).Unix(),
			"iss": "identity_service", // watch identity-service-config.yml
			"aud": []string{"client"},
		})
		wrongKey := symKey[:len(symKey)-1]
		wrongToken, err := refreshJWT.SignedString(wrongKey)
		require.NoError(t, err)

		resp := doRefreshJWTRequest(t, wrongToken)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body := common.GetBody(t, resp.Body)
		require.Equal(t, "invalid_token", body.ErrorType)
	})

	t.Run("DeniesExpired", func(t *testing.T) {
		refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"typ": "refresh",
			"jti": uuid.NewString(),
			"iat": time.Now().Add(-10 * time.Minute).Unix(),
			"exp": time.Now().Add(-1 * time.Minute).Unix(),
			"iss": "identity_service", // watch identity-service-config.yml
			"aud": []string{"client"},
		})
		refreshToken, err := refreshJWT.SignedString(symKey)
		require.NoError(t, err)

		resp := doRefreshJWTRequest(t, refreshToken)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body := common.GetBody(t, resp.Body)
		require.Equal(t, "refresh_token_expired", body.ErrorType)
	})

	t.Run("DeniesInvalidated", func(t *testing.T) {
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

		resp := doRefreshJWTRequest(t, refreshToken)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		respBody := common.GetBody(t, resp.Body)
		require.Equal(t, "refresh_token_invalidated", respBody.ErrorType)
	})
}

func doRefreshJWTRequest(t *testing.T, token string) *http.Response {
	type Request struct {
		RefreshToken string `json:"refresh_token"`
	}
	reqBody, _ := json.Marshal(Request{
		RefreshToken: token,
	})

	req, err := http.NewRequest(http.MethodPost, common.ServiceUrl+"/v1.0/refresh-token", bytes.NewReader(reqBody))
	require.NoError(t, err)
	req.Header.Add(common.HeaderIdemotencyKey, uuid.NewString())

	signOutResp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return signOutResp
}

func getKey(t *testing.T, path string) []byte {
	key, err := os.ReadFile(path)
	require.NoError(t, err)
	return key
}
