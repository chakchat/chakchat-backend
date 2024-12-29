package jwt

import (
	"crypto/rsa"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type InternalToken string

const (
	ClaimType = "typ"
	ClaimName = "name"
	ClaimSub  = "sub"
)

var ErrTokenExpired = errors.New("token expired")

type Token string

type Pair struct {
	Access  Token
	Refresh Token
}

type Config struct {
	SigningMethod string

	Lifetime time.Duration
	Issuer   string
	Audience []string

	Type string

	publicKey  *rsa.PublicKey  // Used for asymmetric signing
	privateKey *rsa.PrivateKey // Used for asymmetric signing

	SymmetricKey []byte // Used for symmetric signing
}

func (c *Config) RSAPublicOnlyKey(key []byte) error {
	var err error
	c.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(key)
	return err
}

func (c *Config) RSAKeys(privateKey []byte) error {
	var err error
	c.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	c.publicKey = &c.privateKey.PublicKey
	return err
}

type Claims map[string]any

var nowFunc = func() time.Time {
	return time.Now()
}

func Generate(config *Config, claims Claims) (Token, error) {
	fillBasicClaims(config, claims)

	meth := jwt.GetSigningMethod(config.SigningMethod)
	if meth == nil {
		return "", errors.New("no such signing method")
	}

	token := jwt.NewWithClaims(meth, jwt.MapClaims(claims))

	var signed string
	var err error
	if isSymmetricAlg(meth.Alg()) {
		signed, err = token.SignedString(config.SymmetricKey)
	} else {
		signed, err = token.SignedString(config.privateKey)
	}

	return Token(signed), err
}

// If you need to parse token use only Parse function because it also verifies.
func Verify(config *Config, token Token) error {
	_, err := Parse(config, token)
	return err
}

// If you need to parse token use only ParseWithAud function because it also verifies.
func VerifyWithAud(config *Config, token Token, aud string) error {
	_, err := Parse(config, token)
	return err
}

// Parse also verifies token
func Parse(config *Config, token Token) (Claims, error) {
	claims, err := parse(config, token)
	return Claims(claims), err
}

// In comparison with just Parse it also checks audience
// Parse also verifies token
func ParseWithAud(config *Config, token Token, aud string) (Claims, error) {
	claims, err := parse(config, token)
	if err != nil {
		return nil, err
	}

	if !claims.VerifyAudience(aud, true) {
		return nil, errors.New("invalid audience")
	}

	return Claims(claims), nil
}

func fillBasicClaims(config *Config, claims Claims) {
	if config.Issuer != "" {
		claims["iss"] = config.Issuer
	}
	if len(config.Audience) != 0 {
		claims["aud"] = config.Audience
	}
	if config.Type != "" {
		claims[ClaimType] = config.Type
	}

	now := nowFunc()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(config.Lifetime).Unix()
}

func parse(config *Config, token Token) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(string(token), claims, func(t *jwt.Token) (interface{}, error) {
		if isSymmetricAlg(config.SigningMethod) {
			return config.SymmetricKey, nil
		}
		return config.publicKey, nil
	})
	if err != nil {
		if jwtErr, ok := err.(jwt.ValidationError); ok && jwtErr.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, ErrTokenExpired
		}
		return nil, err
	}
	if !parsed.Valid {

		return nil, errors.New("token invalid")
	}

	if config.Issuer != "" && !claims.VerifyIssuer(config.Issuer, true) {
		return nil, errors.New("invalid issuer")
	}

	// json unmarshals int to float64 so I don't wanna expose this specificty to the usage
	claims["iat"] = int64(claims["iat"].(float64))
	claims["exp"] = int64(claims["exp"].(float64))
	return claims, nil
}

func isSymmetricAlg(alg string) bool {
	return strings.HasPrefix(alg, "HS")
}
