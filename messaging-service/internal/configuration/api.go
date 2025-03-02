package configuration

import (
	"log"
	"os"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
)

func GinEngine(service *Service, conf *Config) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(auth.NewJWT(&auth.JWTConfig{
		Conf: &jwt.Config{
			SigningMethod: conf.JWT.SigningMethod,
			Issuer:        conf.JWT.Issuer,
			Audience:      conf.JWT.Audience,
			Type:          "internal_access",
			SymmetricKey:  readKey(conf.JWT.KeyFilePath),
		},
		DefaultHeader: "X-Internal-Token",
	}))

	return r
}

func readKey(path string) []byte {
	key, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Cannot read file: %s", err)
	}
	return key
}
