package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const githubGraphqlApi = "https://api.github.com/graphql"
const jsonContentType = "application/json"

var githubOauthConfig = &oauth2.Config{
	ClientID:     "8ecab7653a8b3804c2c8",
	ClientSecret: "12722aeebc07090cf0e084f9fdb0fb0ac0cac44a",
	Endpoint:     github.Endpoint,
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/oauth2/github", githubOauth)
	http.HandleFunc("/oauth2/receive", githubOauthHandleReceive)
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
		<form action="/oauth2/github" method="post">
			<input type="email" name="email" />
			<input type="submit" value="Github Login"/>
		</form>
	</body>
	</html>`
	fmt.Fprint(w, html)
}

func githubOauth(w http.ResponseWriter, r *http.Request) {
	redirectUrl := githubOauthConfig.AuthCodeURL("8888")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func githubOauthHandleReceive(w http.ResponseWriter, r *http.Request) {
	log.Println("processing receive...")

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
	requestBody := strings.NewReader(`{
		"query": "query {viewer {id}}"
	}`)
	res, err := client.Post(githubGraphqlApi, jsonContentType, requestBody)
	if err != nil {
		http.Error(w, "Could not make request to github api", http.StatusInternalServerError)
		return
	}

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Could not read response body", http.StatusInternalServerError)
		return
	}

	log.Println(string(bs))
}
