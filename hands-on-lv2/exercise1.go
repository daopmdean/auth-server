package main

import (
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var db = map[string][]byte{}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := buildHtml("<email>", "<password>")
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
	}

	db[email] = hashed
	html := buildHtml(email, string(db[email]))
	io.WriteString(w, html)
}

func buildHtml(email, password string) string {
	return `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>Lv 2 - Excercise 1</title>
	</head>
	<body>
		<p>` + email + `</p>
		<p>` + password + `</p>
		<form action="/register" method="post">
			<input type="email" name="email" /></br>
			<input type="text" name="password" />
			<input type="submit" />
		</form>
	</body>
	</html>`
}
