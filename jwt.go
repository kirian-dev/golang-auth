package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const key = "my secure jwt key"

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

func createToken(username string) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("Error signed token %w", err)
	}

	return tokenString, nil
}
func parseToken(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error parsing token %w", err)
	}

	if claim, ok := token.Claims.(*Claims); ok && token.Valid {
		return claim, nil
	}
	return nil, fmt.Errorf("invalid token")
}
