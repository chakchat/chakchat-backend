package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/chakchat/chakchat-backend/shared/go/jwt"
)

const (
	EnvHTTPAddr       = "USER_SERVICE_HTTP_ADDR"
	EnvPrivateKeyPath = "/app/keys/rsa"
)

// func TestHTTP_GetByID(t *testing.T) {
// 	addr := getURL(t)
// 	jwt := genJWT(t, JWTModeValid, jwt.Claims{
// 		jwt.ClaimName: "bob",
// 		jwt.ClaimSub:  "",
// 	})

// }

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
		key, err := os.ReadFile(EnvPrivateKeyPath)
		if err != nil {
			t.Fatalf("Cannot read %s file: %s", EnvPrivateKeyPath, err)
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
