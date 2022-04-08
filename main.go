package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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
	pass := "pass"
	hashed, err := hashPassword(pass)
	if err != nil {
		log.Fatalln("Error hashing password")
	}

	log.Println(hashed)

	err = comparePassword(hashed, []byte("pass"))
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
	h := hmac.New(sha512.New, keys[currentKid].key)

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

func generateNewKey() error {
	newKey := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, newKey)
	if err != nil {
		return fmt.Errorf("Error in generateNewKey: %w", err)
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("Error in generateNewKey: %w", err)
	}

	keys[uid.String()] = key{
		key:     newKey,
		created: time.Now(),
	}
	currentKid = uid.String()

	return nil
}

type key struct {
	key     []byte
	created time.Time
}

var currentKid = ""
var keys = map[string]key{}

func createToken(c *UserClaims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	signedToken, err := jwtToken.SignedString(keys[currentKid].key)
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

		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("Invalid Key Id")
		}

		k, ok := keys[kid]
		if !ok {
			return nil, fmt.Errorf("Invalid Key Id")
		}

		return k.key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error in parseToken: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Error in parseToken: token is not valid")
	}

	return token.Claims.(*UserClaims), nil
}
