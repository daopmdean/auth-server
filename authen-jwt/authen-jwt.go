package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/submit", bar)
	http.ListenAndServe(":8080", nil)
}

var key = []byte("My Secret Key")

type myClaims struct {
	jwt.StandardClaims
	Email string
}

func getJwt(msg string) (string, error) {
	claims := myClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			Issuer:    "authen-jwt",
		},
		Email: msg,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	s, err := t.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Error in getJwt %w", err)
	}

	return s, nil
}

func bar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	code, err := getJwt(email)
	if err != nil {
		http.Error(w, "Could not get Jwt", http.StatusInternalServerError)
		return
	}

	c := http.Cookie{
		Name:  "session",
		Value: code,
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func foo(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{}
	}

	t, err := jwt.ParseWithClaims(c.Value, &myClaims{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	message := "Not Logged In"
	if err != nil {
		html := buildHtml(c.Value, message)
		io.WriteString(w, html)
		return
	}

	isLoggedIn := false
	if t.Valid {
		isLoggedIn = true
	} else {
		message = "Invalid Token"
	}

	if isLoggedIn {
		message = "Logged In"
	}

	html := buildHtml(c.Value, message)
	io.WriteString(w, html)
}

func buildHtml(value, message string) string {
	return `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>HMAC Example</title>
	</head>
	<body>
		<p>` + value + `</p>
		<p>` + message + `</p>
		<form action="/submit" method="post">
			<input type="email" name="email" />
			<input type="submit" />
		</form>
	</body>
	</html>`
}
