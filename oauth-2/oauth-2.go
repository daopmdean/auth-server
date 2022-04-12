package main

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOauthConfig = &oauth2.Config{
	ClientID:     "8ecab7653a8b3804c2c8",
	ClientSecret: "12722aeebc07090cf0e084f9fdb0fb0ac0cac44a",
	Endpoint:     github.Endpoint,
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/oauth/github", githubOauth)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>OAuth2 Example</title>
	</head>
	<body>
		<form action="/oauth/github" method="post">
			<input type="email" name="email" />
			<input type="submit" value="Github Login"/>
		</form>
	</body>
	</html>`
	fmt.Fprint(w, html)
}

func githubOauth(w http.ResponseWriter, r *http.Request) {
	redirectUrl := githubOauthConfig.AuthCodeURL("state8888")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func githubOauthHandleReceive(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != "8888" {
		http.Error(w, "Invalid State", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Could not login", http.StatusInternalServerError)
		return
	}

	src := githubOauthConfig.TokenSource(ctx, token)
	client := oauth2.NewClient(ctx, src)
}
