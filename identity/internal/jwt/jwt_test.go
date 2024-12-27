package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GeneratesExact(t *testing.T) {
	// Arrange
	issTime, _ := time.Parse(time.DateTime, "2024-12-27 14:40:01")
	// Mocking this function to make Generate() deterministic
	nowFunc = func() time.Time { return issTime }

	claims := Claims{
		"sub":         "122f4915aa124492bd79539013819cd3",
		"name":        "Joshua Kimmich",
		"best_player": false, // Sorry, Joshua...
	}

	t.Run("HS512_Signing", func(t *testing.T) {
		config := Config{
			SigningMethod: "HS512",
			Lifetime:      3 * time.Minute,
			Issuer:        "idk_issuer",
			Audience:      []string{"idk_audience"},
			SymmetricKey:  []byte(RawHS512Key),
		}
		const trueToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiaWRrX2F1ZGllbmNlIl0sImJlc3RfcGxheWVyIjpmYWxzZSwiZXhwIjoxNzM1MzEwNTgxLCJpYXQiOjE3MzUzMTA0MDEsImlzcyI6Imlka19pc3N1ZXIiLCJuYW1lIjoiSm9zaHVhIEtpbW1pY2giLCJzdWIiOiIxMjJmNDkxNWFhMTI0NDkyYmQ3OTUzOTAxMzgxOWNkMyJ9.FZZaidtqFlNKEzGoOBwVt4OLo5AFp8TR7iK6GpK140HYp8vnoWeXRKxht64Tv7LGKRjZAgOzwUQNUY_HMSFN-A"

		// Act
		token, err := Generate(&config, claims)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, string(token), trueToken)
	})

	t.Run("RS256_Signing", func(t *testing.T) {
		config := Config{
			SigningMethod: "RS256",
			Lifetime:      3 * time.Minute,
			Issuer:        "idk_issuer",
			Audience:      []string{"idk_audience"},
		}
		config.RSAKeys([]byte(RawRSA2048PrivateKey))
		const trueToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiaWRrX2F1ZGllbmNlIl0sImJlc3RfcGxheWVyIjpmYWxzZSwiZXhwIjoxNzM1MzEwNTgxLCJpYXQiOjE3MzUzMTA0MDEsImlzcyI6Imlka19pc3N1ZXIiLCJuYW1lIjoiSm9zaHVhIEtpbW1pY2giLCJzdWIiOiIxMjJmNDkxNWFhMTI0NDkyYmQ3OTUzOTAxMzgxOWNkMyJ9.Br-JiNFgoLMS7Z2hjbC7bsjLxPfDIWYvXmkg53ikvg2zFE3DAxI_d9XC8em9yToMi5aFP4ELvs7f6jzGC25t0UYpR9ZtHXvFNzsN-UlHI343RAsJwn8UiQnMQYXbqiiunpbfl8wDlQKeh59umGw7qtpge3fjBnR8hDuvg88VXWa_j8Nv2QbsGyQFaP-N8x1prSWU0Tm7Tx7eHjTNype6q12oL40ofrMt9UIh2Vr2c4kzsaj4Bjuvx2H62Hp_E1WtBYqO7gWN_5M40w70rgBc8yP73ArckbvzuRjQCknUAlpOwUCl-s7_cnYFZqoHVsQ_u4d3EK0c1nwA4j-Nd8kJHA"

		// Act
		token, err := Generate(&config, claims)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, string(token), trueToken)
	})
}

const (
	RawHS512Key         = `bf566321d4b63ffc8b7211d491e2def3600f7d934f7328c2956bf996be06ea9b6ad5421ef233ce97d95664037b0e79cb10bcd6e8dfcfddace2a24e0cbf496f50`
	RawRSA2048PublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoM3ofM5rK+Nfa/sJsrQu
ooLxEc35G+vyIRJh/gpXDw+dxJKd/88KcR9eSw6e8e+Ut8eMShtufhKrvU9/533y
0033y08BDUpfq8EpsUFk1hhaIDva1l2cTpQWLfwWl25hlBijxyaYTTVF7U4YhTO2
jzB65QZ4s+UVOkicMMnJWGvry6JU9zhGj6zE28te/9v9hNe66/0X/+84iUs7dMmR
ZVyr3oV+iy8zIKJHLUupQ9/L+elgEJ7NznFFW4z/8YODEbKT8JLrEoJH+H6sl2bY
B6b2eWapQ82mBIyQe8I0u8LFkL2+v0szKB9xGu9fQFzME8Mxzb6YANhDDS/Ow6Bg
6QIDAQAB
-----END PUBLIC KEY-----
`
	RawRSA2048PrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEugIBADANBgkqhkiG9w0BAQEFAASCBKQwggSgAgEAAoIBAQCgzeh8zmsr419r
+wmytC6igvERzfkb6/IhEmH+ClcPD53Ekp3/zwpxH15LDp7x75S3x4xKG25+Equ9
T3/nffLTTffLTwENSl+rwSmxQWTWGFogO9rWXZxOlBYt/BaXbmGUGKPHJphNNUXt
ThiFM7aPMHrlBniz5RU6SJwwyclYa+vLolT3OEaPrMTby17/2/2E17rr/Rf/7ziJ
Szt0yZFlXKvehX6LLzMgokctS6lD38v56WAQns3OcUVbjP/xg4MRspPwkusSgkf4
fqyXZtgHpvZ5ZqlDzaYEjJB7wjS7wsWQvb6/SzMoH3Ea719AXMwTwzHNvpgA2EMN
L87DoGDpAgMBAAECgf9jPEjDqEhZ6sMV5BJ9MjJS2A3lo3CedNM0HVNzCOG96GPz
SQgH1zH6d0jwqjzlBIrvc/JMfLhV/kGkaFgkDRSmAvFavWVxwC3H1LNnWompC+sD
Q692IsdALLpVPZpcesXtzmq+63KAAAV6l2a7FN4tFqMnRVNJ0B7P7uM3cOlmZZ8g
prXO0IBWQ/PWDjuC/BFWrgkq4e6/ZJgEYuoNC1Gzvh3x4+Pc29oXlvWyK6cHcNxc
G0Y5sSI3haoLCI3dzJQwN8htsTtwYDdjeqMwhpUkat37gZ4Y7W8zUe82lCC73jsC
eAvIC72vMNSgf6JDC3efnuSeAT6F2YsOD3WQxpECgYEA1h9WXAtts1ueC6QBoUlm
jbrfPWEsvoO73TC9/OJ31BpGq24Yr0Om4BV1rig9YPm/fZUEGysqg5ILYvEJoNRE
Y2MVLoVeFOFpW2sHkoXyaZpW5qrB+LepZGKARLGM8IYnpFRTlSThFFkKEKmNAF1V
y8GmVTKqfuO50mLTPFwo5TECgYEAwEELq5B3Ow8Z7fOUg09Ce+VwxjB8p8RIlPFx
5dWTpjmCfUunjKr4eSuYj/f8nUfT8/EjVxtKToxSwwdGcMj65WZuUj9CtM8v/XAq
6JlzT0906Ss0GHuCESdcHhKq8Ez6kqTTl9z6UQZmxRkleXgGp3oAm4BSJjXZEZzl
+cS4qTkCgYAZQ+dXww11rWjPrNF4a4XLUXKH9pBmBntDVT4Fud8zysnt7nbBL3Vg
WYfiPeNILw/2TIAIiKZikff/+7sMHB/ZrlZQf/Ii+poI7G8fTejVpx176EgtBdbZ
/nluIZkkxF+nF0ApiAl68iqq3qbBlUHLYhUzVmAhytMhTQHpzGIS8QKBgDZR6o82
CUopkST3Xq3fNiS1hjCpQH9SaUOUGJ9cwhQEScdHGfcX047A76E16y0xP0S8jESv
VEZvRW8PXiq9zo4EbAVXFGzr4V5VU/pWaQsuoxTCfTyxoOVh3pgspBmzVlUatyJA
cIV2LpFf8oOokxC82vEUx6E+M6/TSfNRTu+ZAoGAJckxBLgqHpqhWXo5qiO59Efs
57c4JdBd8ScMPKaE5MoXSybLgVQyhcDcngeRyWYRsYHfqRwpdr4nDlieccL2nnHT
MCXxO7ZuKZYZblPPFQ59bcJRY69s982i6qw1VvvqRg3Y4lcNQAGE7RctQtF03C/B
LNmw9p4duLj8vXmdTts=
-----END PRIVATE KEY-----
`
)
