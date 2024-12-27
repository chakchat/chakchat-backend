package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type InternalToken string

type Token string

type Pair struct {
	Access  Token
	Refresh Token
}

type Config struct {
	SecureKey []byte
	Lifetime  time.Duration
	Issuer    string
	Audience  string
}

type Claims map[string]any

var nowUTCFunc = func() time.Time {
	return time.Now().UTC()
}

func Generate(config *Config, claims Claims) (Token, error) {
	if config.Issuer != "" {
		claims["iss"] = config.Issuer
	}
	if config.Audience != "" {
		claims["aud"] = config.Audience
	}

	now := nowUTCFunc()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(config.Lifetime).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims(claims))
	signed, err := token.SignedString(config.SecureKey)

	return Token(signed), err
}

// If you need to parse token use only Parse function because it also verifies.
func Verify(config *Config, token Token) error {
	_, err := Parse(config, token)
	return err
}

// Parse also verifies token
func Parse(config *Config, token Token) (Claims, error) {
	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(string(token), claims, func(t *jwt.Token) (interface{}, error) {
		return config.SecureKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("token invalid")
	}
	// json unmarshals int to float64 so I don't wanna expose this specificty to the usage
	claims["iat"] = int64(claims["iat"].(float64))
	claims["exp"] = int64(claims["exp"].(float64))
	return Claims(claims), nil
}
