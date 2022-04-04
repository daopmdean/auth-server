package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var key []byte

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
	h := hmac.New(sha256.New, key)

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
