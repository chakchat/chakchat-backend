package services

import (
	"context"

	"github.com/chakchat/chakchat/backend/identity/internal/jwt"
)

type IdentityIssuer struct {
	userConf     *jwt.Config
	internalConf *jwt.Config
}

func NewIdentityIssuer(userConf, internalConf *jwt.Config) *IdentityIssuer {
	return &IdentityIssuer{
		userConf:     userConf,
		internalConf: internalConf,
	}
}

func (i *IdentityIssuer) Idenitfy(ctx context.Context, token jwt.Token) (jwt.InternalToken, error) {
	claims, err := jwt.Parse(i.userConf, token)
	if err != nil {
		return "", err
	}

	internalClaims := extractInternal(claims)

	internalToken, err := jwt.Generate(i.internalConf, internalClaims)
	if err != nil {
		return "", err
	}
	return jwt.InternalToken(internalToken), nil
}

func extractInternal(claims jwt.Claims) jwt.Claims {
	return jwt.Claims{
		"sub":  claims["sub"],
		"name": claims["name"],
	}
}
