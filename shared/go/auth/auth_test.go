package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chakchat/chakchat/backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func Test_AuthJwtMiddleware(t *testing.T) {
	conf := &JWTConfig{
		Conf: &jwt.Config{
			SigningMethod: "HS512",
			Lifetime:      2 * time.Minute,
			Issuer:        "iss",
			Audience:      []string{"aud"},
			Type:          "internal_access",
			SymmetricKey:  []byte("I_DONT_WANNA_GEN_A_KEY_SO_I_PASTE_UUID_2a7e232f79ef4a3caced0c504553afa9"),
		},
		Aud:           "",
		DefaultHeader: "",
	}

	t.Run("NoAuthorizationHeader", func(t *testing.T) {
		r := gin.New()

		r.Use(NewJWT(conf))
		r.GET("/", func(c *gin.Context) {
			t.Fatal("This code must be unreachable")
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", bytes.NewReader(nil))

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("InvalidJWT", func(t *testing.T) {
		r := gin.New()

		r.Use(NewJWT(conf))
		r.GET("/", func(c *gin.Context) {
			t.Fatal("This code must be unreachable")
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", bytes.NewReader(nil))
		// Generated at jwt.io
		req.Header.Add("Authorization", "Bearer "+`eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwidXNlcm5hbWUiOiJqb2huX2RvZSIsImFkbWluIjp0cnVlLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTYxNjIzOTAyMn0.qN_9N1TEJmIPfbZdfQn2MNt7t-it2iowCUxmMJHFp9chzvWMhdux0r1ryC0hk2xhB-iK8hU7Vilzhg6hTF_uGA`)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("Autorized", func(t *testing.T) {
		r := gin.New()

		r.Use(NewJWT(conf))
		r.GET("/", func(c *gin.Context) {
			claims := GetClaims(c)
			require.Contains(t, claims, ClaimName)
			require.Contains(t, claims, ClaimUsername)
			require.Contains(t, claims, ClaimId)
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", bytes.NewReader(nil))
		// Generated at jwt.io
		req.Header.Add("Authorization", "Bearer "+`eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwidXNlcm5hbWUiOiJqb2huX2RvZSIsImFkbWluIjp0cnVlLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTkxNjIzOTAyMiwiaXNzIjoiaXNzIiwidHlwIjoiaW50ZXJuYWxfYWNjZXNzIiwiYXVkIjoiYXVkIn0.8RDdDkCqp-ohTWbzHa_k0KV4tutG1Rfr85wQCwI2PfM8ARKr6zCz5qivJdfCIn3Q94F2cVQh6SvupFOZRFhBYA`)

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}
