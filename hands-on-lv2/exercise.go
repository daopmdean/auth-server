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
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := buildIndexHtml("<email>", "<password>")
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
	html := buildIndexHtml(email, string(db[email]))
	io.WriteString(w, html)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	err := bcrypt.CompareHashAndPassword(db[email], []byte(password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusInternalServerError)
		return
	}

	html := buildIndexHtml(email, string(db[email]))
	io.WriteString(w, html)
}

func buildIndexHtml(email, password string) string {
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
