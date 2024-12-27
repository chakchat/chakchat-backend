package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func Test_GeneratesExact(t *testing.T) {
	// Assert
	issTime, _ := time.Parse(time.DateTime, "2024-12-27 14:40:01")
	// Mocking this function to make Generate() deterministic
	nowUTCFunc = func() time.Time { return issTime }
	config := Config{
		SecureKey: []byte("c057130fa78d415e9bcb5ab081e7963b"),
		Lifetime:  3 * time.Minute,
		Issuer:    "idk_issuer",
		Audience:  "idk_audience",
	}
	claims := Claims{
		"sub":         "122f4915aa124492bd79539013819cd3",
		"name":        "Joshua Kimmich",
		"best_player": false, // Sorry, Joshua...
	}

	// Act
	token, err := Generate(&config, claims)

	// Assert
	assert.NoError(t, err)
	// Created in https://jwt.io/
	const trueToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpZGtfYXVkaWVuY2UiLCJiZXN0X3BsYXllciI6ZmFsc2UsImV4cCI6MTczNTMxMDU4MSwiaWF0IjoxNzM1MzEwNDAxLCJpc3MiOiJpZGtfaXNzdWVyIiwibmFtZSI6Ikpvc2h1YSBLaW1taWNoIiwic3ViIjoiMTIyZjQ5MTVhYTEyNDQ5MmJkNzk1MzkwMTM4MTljZDMifQ.AqvA6FuQNhTG2mKAn878dFP9I2m4BWwktPEaxZVynMRn4fRyhKXbNfaIlCsniuxWa8BKsf15Rhs_t3zDsx8wvw"
	assert.Equal(t, string(token), trueToken)
}

func Test_ParsesExact(t *testing.T) {
	// Arrange
	issTime, _ := time.Parse(time.DateTime, "2024-12-27 14:40:01")
	// Mocking it to make jwt package functions deterministic
	jwt.TimeFunc = func() time.Time { return issTime }
	config := Config{
		SecureKey: []byte("c057130fa78d415e9bcb5ab081e7963b"),
		Lifetime:  3 * time.Minute,
		Issuer:    "idk_issuer",
		Audience:  "idk_audience",
	}
	// Created in https://jwt.io/
	const token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpZGtfYXVkaWVuY2UiLCJiZXN0X3BsYXllciI6ZmFsc2UsImV4cCI6MTczNTMxMDU4MSwiaWF0IjoxNzM1MzEwNDAxLCJpc3MiOiJpZGtfaXNzdWVyIiwibmFtZSI6Ikpvc2h1YSBLaW1taWNoIiwic3ViIjoiMTIyZjQ5MTVhYTEyNDQ5MmJkNzk1MzkwMTM4MTljZDMifQ.AqvA6FuQNhTG2mKAn878dFP9I2m4BWwktPEaxZVynMRn4fRyhKXbNfaIlCsniuxWa8BKsf15Rhs_t3zDsx8wvw"

	// Act
	claims, err := Parse(&config, token)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, Claims{
		"aud":         "idk_audience",
		"best_player": false,
		"exp":         int64(1735310581),
		"iat":         int64(1735310401),
		"iss":         "idk_issuer",
		"name":        "Joshua Kimmich",
		"sub":         "122f4915aa124492bd79539013819cd3",
	}, claims)
}

func Test_InvalidSecureKey(t *testing.T) {
	// Arrange
	issTime, _ := time.Parse(time.DateTime, "2024-12-27 14:40:01")
	// Mocking it to make jwt package functions deterministic
	jwt.TimeFunc = func() time.Time { return issTime }
	config := Config{
		SecureKey: []byte("it-is-invalid-key"),
		Lifetime:  3 * time.Minute,
		Issuer:    "idk_issuer",
		Audience:  "idk_audience",
	}
	// Created in https://jwt.io/
	const token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpZGtfYXVkaWVuY2UiLCJiZXN0X3BsYXllciI6ZmFsc2UsImV4cCI6MTczNTMxMDU4MSwiaWF0IjoxNzM1MzEwNDAxLCJpc3MiOiJpZGtfaXNzdWVyIiwibmFtZSI6Ikpvc2h1YSBLaW1taWNoIiwic3ViIjoiMTIyZjQ5MTVhYTEyNDQ5MmJkNzk1MzkwMTM4MTljZDMifQ.AqvA6FuQNhTG2mKAn878dFP9I2m4BWwktPEaxZVynMRn4fRyhKXbNfaIlCsniuxWa8BKsf15Rhs_t3zDsx8wvw"

	// Act
	_, err := Parse(&config, token)

	// Assert
	assert.Error(t, err)
}
