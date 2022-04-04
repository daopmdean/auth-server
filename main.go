package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
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
