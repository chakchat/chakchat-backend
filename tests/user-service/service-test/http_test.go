package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"test/userservice"
	"testing"
	"time"

	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	EnvHTTPAddr    = "USER_SERVICE_HTTP_ADDR"
	PrivateKeyPath = "/app/keys/rsa"
)

func TestHTTP_GetUser(t *testing.T) {
	baseUrl := getURL(t)

	type Response struct {
		ErrorType string `json:"error_type"`
		Data      *struct {
			Id          uuid.UUID `json:"id"`
			Name        string    `json:"name"`
			Username    string    `json:"username"`
			Phone       string    `json:"phone"`
			DateOfBirth string    `json:"date_of_birth"`
		} `json:"data"`
	}

	requestFunc := func(path string) (*http.Response, Response) {
		username, _ := getUniqueUser()
		jwt := genJWT(t, JWTModeValid, jwt.Claims{
			jwt.ClaimName:     "Bob",
			jwt.ClaimSub:      uuid.NewString(),
			jwt.ClaimUsername: username,
		})

		req, err := http.NewRequest(
			http.MethodGet,
			baseUrl+path,
			bytes.NewReader(nil),
		)
		require.NoError(t, err)
		req.Header.Add("Authorzation", "Bearer "+string(jwt))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("HTTP request failed: %s", err)
		}

		defer resp.Body.Close()
		bodyB, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var body Response
		if err := json.Unmarshal(bodyB, &body); err != nil {
			t.Fatalf("Response body unmarshalling failed: %s", err)
		}

		return resp, body
	}

	t.Run("ByID", func(t *testing.T) {
		t.Run("NotFound", func(t *testing.T) {
			resp, body := requestFunc("/users/" + uuid.NewString())
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
			require.Equal(t, "user_not_found", body.ErrorType)
		})

		t.Run("Success", func(t *testing.T) {
			username, phone := getUniqueUser()

			grpcClient, closeFunc := connectGRPC(t)
			defer closeFunc()

			createResp, err := grpcClient.CreateUser(context.Background(), &userservice.CreateUserRequest{
				PhoneNumber: phone,
				Name:        username,
				Username:    username,
			})
			require.NoError(t, err)
			require.Equal(t, userservice.CreateUserStatus_CREATED, createResp.Status)

			resp, body := requestFunc("/users/" + createResp.UserId.Value)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, phone, body.Data.Phone)
			require.Equal(t, username, body.Data.Username)
			require.Equal(t, username, body.Data.Name)
			require.Equal(t, createResp.UserId.Value, body.Data.Id.String())
		})
	})

	t.Run("ByUsername", func(t *testing.T) {
		t.Run("NotFound", func(t *testing.T) {
			resp, body := requestFunc("/users/username/" + uuid.NewString())
			require.Equal(t, http.StatusNotFound, resp.StatusCode)
			require.Equal(t, "user_not_found", body.ErrorType)
		})

		t.Run("Success", func(t *testing.T) {
			username, phone := getUniqueUser()

			grpcClient, closeFunc := connectGRPC(t)
			defer closeFunc()

			createResp, err := grpcClient.CreateUser(context.Background(), &userservice.CreateUserRequest{
				PhoneNumber: phone,
				Name:        username,
				Username:    username,
			})
			require.NoError(t, err)
			require.Equal(t, userservice.CreateUserStatus_CREATED, createResp.Status)

			resp, body := requestFunc("/users/username/" + username)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, phone, body.Data.Phone)
			require.Equal(t, username, body.Data.Username)
			require.Equal(t, username, body.Data.Name)
			require.Equal(t, createResp.UserId.Value, body.Data.Id.String())
		})
	})
}

func TestHTTP_Teapot(t *testing.T) {
	addr := getURL(t)
	resp, err := http.Get(addr + "/v1.0/are-you-a-real-teapot")
	if err != nil {
		t.Fatalf("HTTP request failed: %s", err)
	}
	if resp.StatusCode != 418 {
		t.Fatal(`418 "I'm a teapot" status expected`)
	}
}

const (
	JWTModeValid = iota
	JWTModeInvalidIssuer
	JWTModeInvalidTokenType
	JWTModeInvalidKey
)

var jwtPrivateKey []byte

func genJWT(t *testing.T, mode int, claims jwt.Claims) jwt.Token {
	key := make([]byte, len(jwtPrivateKey))
	if jwtPrivateKey != nil {
		copy(jwtPrivateKey, key)
	} else {
		key, err := os.ReadFile(PrivateKeyPath)
		if err != nil {
			t.Fatalf("Cannot read %s file: %s", PrivateKeyPath, err)
		}
		jwtPrivateKey = key
	}

	var (
		issuer    = "identity_service"
		tokenType = "internal_access"
	)
	switch mode {
	case JWTModeInvalidIssuer:
		issuer = "invalid_issuer"
	case JWTModeInvalidTokenType:
		tokenType = "invalid_token_type"
	case JWTModeInvalidKey:
		key = append(key, 69) // Just mess it
	}

	config := &jwt.Config{
		SigningMethod: "RS256",
		Lifetime:      time.Minute,
		Issuer:        issuer,
		Audience:      []string{"identity_service", "user_service"},
		Type:          tokenType,
		SymmetricKey:  key,
	}

	token, err := jwt.Generate(config, claims)
	if err != nil {
		t.Fatalf("Generating JWT failed: %s", err)
	}
	return token
}

func getURL(t *testing.T) string {
	addr := os.Getenv(EnvHTTPAddr)
	if addr == "" {
		t.Fatalf("You should pass %s environment variable", EnvHTTPAddr)
	}
	return addr
}

var uniqueCounter = atomic.Int32{}

func getUniqueUser() (username, phone string) {
	i := int(uniqueCounter.Add(1))

	phone = "79000000000"
	suff := strconv.Itoa(i)
	phone = phone[:len(phone)-len(suff)] + suff

	username = "user_with_phone_" + phone
	return username, phone
}
