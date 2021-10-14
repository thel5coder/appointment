package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtCredential struct {
	TokenSecret         string
	ExpiredToken        int
	RefreshTokenSecret  string
	ExpiredRefreshToken int
}

type CustomClaims struct {
	jwt.StandardClaims
	Email    string `json:"email"`
	Session  string `json:"session"`
}

func (cred JwtCredential) GetToken(session, issuer, id string) (string, string, error) {
	expirationTime := time.Now().Add(time.Duration(cred.ExpiredToken) * time.Hour).Unix()
	unixTimeUtc := time.Unix(expirationTime, 0)
	unitTimeInRFC3339 := unixTimeUtc.UTC().Format(time.RFC3339)

	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Id:        id,
			Issuer:    issuer,
		},
		Session: session,
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(cred.TokenSecret))

	return token, unitTimeInRFC3339, err
}

func (cred JwtCredential) GetRefreshToken(session, issuer, id string) (string, string, error) {
	expirationTime := time.Now().Add(time.Duration(cred.ExpiredRefreshToken) * time.Hour).Unix()
	unixTimeUtc := time.Unix(expirationTime, 0)
	unitTimeInRFC3339 := unixTimeUtc.UTC().Format(time.RFC3339)

	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Id:        id,
			Issuer:    issuer,
		},
		Session: session,
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(cred.RefreshTokenSecret))

	return token, unitTimeInRFC3339, err
}
