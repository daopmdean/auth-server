package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var userState = map[string]time.Time{}

// key - user id in github; value - user id in our app
var githubConnections = map[string]string{}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/oauth2/github", githubOauth)
	http.HandleFunc("/oauth2/receive", githubOauthHandleReceive)
	http.HandleFunc("/oauth2/google", googleOauth)
	http.HandleFunc("/oauth2/google/receive", googleOauthHandleReceive)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "dsa",
		Value: "dscaascws",
	})
	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>OAuth2 Example</title>
	</head>
	<body>
		<form action="/oauth2/github" method="post">
			<input type="submit" value="Github Login"/>
		</form>
		</br>
		<form action="/oauth2/google" method="post">
			<input type="submit" value="Google Login"/>
		</form>
	</body>
	</html>`
	fmt.Fprint(w, html)
}

func printResponseBody(body io.Reader) error {
	bs, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	log.Println(string(bs))
	return nil
}
