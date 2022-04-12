package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const jwtKey = "The Secret key for jwt"

var db = map[string]user{}

type user struct {
	Email    string
	Password []byte
}

type myClaims struct {
	jwt.StandardClaims
	Email string
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := buildIndexHtml("<token>")
	io.WriteString(w, html)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error generate password", http.StatusInternalServerError)
		return
	}

	db[email] = user{
		Email:    email,
		Password: hashed,
	}

	token, err := createToken(email)
	if err != nil {
		http.Error(w, "Error in register, could not create token", http.StatusInternalServerError)
		return
	}

	html := buildIndexHtml(token)
	io.WriteString(w, html)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	err := bcrypt.CompareHashAndPassword(db[email].Password, []byte(password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusInternalServerError)
		return
	}

	token, err := createToken(email)
	if err != nil {
		http.Error(w, "Error in login, could not create token", http.StatusInternalServerError)
		return
	}

	claims, err := parseToken(token)
	if err != nil {
		http.Error(w, "Error in login, could not parse token", http.StatusInternalServerError)
		return
	}

	c := http.Cookie{
		Name:  "sessionId",
		Value: claims.Email,
	}
	http.SetCookie(w, &c)

	html := buildIndexHtml(token)
	io.WriteString(w, html)
}

func buildIndexHtml(token string) string {
	return `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>Lv 2 - Excercise 1</title>
	</head>
	<body>
		<p>Token: ` + token + `</p>
		<h2>Register</h2>
		<form action="/register" method="post">
			<input type="email" name="email" required/></br>
			<input type="text" name="password" required/>
			<input type="submit" />
		</form>

		<h2>Login</h2>
		<form action="/login" method="post">
			<input type="email" name="email" required/></br>
			<input type="text" name="password" required/>
			<input type="submit" value="Login"/>
		</form>
	</body>
	</html>`
}

func createToken(email string) (string, error) {
	claims := myClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
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
