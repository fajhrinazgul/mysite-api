package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type userAuth struct{}

type credential struct {
	UserID int64 `json:"user_id"`
}

type claims struct {
	credential
	jwt.RegisteredClaims
}

var secretKey = []byte("secretkey")

func getToken(credential credential) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := claims{
		credential: credential,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(secretKey)
	return token, err
}

func verifyToken(token string) (claims, error) {
	var claims claims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) { return secretKey, nil })
	if err != nil {
		return claims, err
	}
	if !jwtToken.Valid {
		return claims, errors.New("your token is expired.")
	}
	return claims, nil
}
