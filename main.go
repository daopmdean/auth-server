package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var key []byte

type UserClaims struct {
	jwt.StandardClaims
	SessionId int64
}

func (u *UserClaims) Valid() error {
	if !u.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("Token expired")
	}

	if u.SessionId == 0 {
		return fmt.Errorf("Invalid Session Id")
	}

	return nil
}

func main() {
	for i := 1; i <= 64; i++ {
		key = append(key, byte(i))
	}
	pass := "somePassword"
	hashed, err := hashPassword(pass)
	if err != nil {
		log.Fatalln("Error hashing password")
	}

	err = comparePassword(hashed, []byte("assd"))
	if err != nil {
		log.Fatalln("Wrong password")
	}

	log.Println("Password correct")
}

func hashPassword(password string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error hashing password: %w", err)
	}
	return hashed, nil
}

func comparePassword(hashedPassword, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return fmt.Errorf("Invalid Password: %w", err)
	}
	return nil
}

func signMessage(msg []byte) ([]byte, error) {
	h := hmac.New(sha512.New, key)

	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("Error in signMessage while hashing: %w", err)
	}

	sig := h.Sum(nil)
	return sig, nil
}

func checkSig(msg, sig []byte) (bool, error) {
	newSig, err := signMessage(msg)
	if err != nil {
		return false, fmt.Errorf("Error in checkSig: %w", err)
	}

	result := hmac.Equal(newSig, sig)
	return result, nil
}

func createToken(c *UserClaims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	signedToken, err := jwtToken.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Errror in createToken while signing token: %w", err)
	}

	return signedToken, nil
}

func parseToken(signedToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("Invalid Signing Algorithm")
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error in parseToken: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Error in parseToken: token is not valid")
	}

	return token.Claims.(*UserClaims), nil
}
