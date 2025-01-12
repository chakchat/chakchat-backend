package services

import (
	"context"
	"errors"
	"log"

	"github.com/chakchat/chakchat/backend/shared/go/jwt"
)

var ErrAccessTokenExpired = errors.New("access token expired")

type IdentityService struct {
	userConf     *jwt.Config
	internalConf *jwt.Config
}

func NewIdentityService(userConf, internalConf *jwt.Config) *IdentityService {
	return &IdentityService{
		userConf:     userConf,
		internalConf: internalConf,
	}
}

func (i *IdentityService) Idenitfy(ctx context.Context, token jwt.Token) (jwt.InternalToken, error) {
	claims, err := jwt.Parse(i.userConf, token)
	if err != nil {
		log.Printf("jwt validation failed: %s", err)
		if err == jwt.ErrTokenExpired {
			return "", ErrAccessTokenExpired
		}
		if err == jwt.ErrInvalidTokenType {
			return "", ErrInvalidTokenType
		}
		return "", ErrInvalidJWT
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
		jwt.ClaimSub:      claims[jwt.ClaimSub],
		jwt.ClaimName:     claims[jwt.ClaimName],
		jwt.ClaimUsername: claims[jwt.ClaimUsername],
	}
}
