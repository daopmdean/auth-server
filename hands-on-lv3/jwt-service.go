package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type myClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func createToken(email string) (string, error) {
	claims := myClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Second).Unix(),
		},
		Email: email,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	token, err := t.SignedString([]byte(jwtKey))
	if err != nil {
		return "", fmt.Errorf("Error in createToken %w", err)
	}

	return token, nil
}

func parseToken(token string) (*myClaims, error) {
	t, err := jwt.ParseWithClaims(token, &myClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("Wrong signing alg method")
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error parse token")
	}

	if !t.Valid {
		return nil, fmt.Errorf("Token Invalid")
	}

	return t.Claims.(*myClaims), nil
}
