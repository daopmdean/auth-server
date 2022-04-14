package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var oauthExp = map[string]time.Time{}
var oauthConnections = map[string]googleRes{}

// key - user id in github; value - user id in our app
var githubConnections = map[string]string{}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/oauth2/github", githubOauth)
	http.HandleFunc("/oauth2/receive", githubOauthHandleReceive)
	http.HandleFunc("/oauth2/google", googleOauth)
	http.HandleFunc("/oauth2/google/receive", googleOauthHandleReceive)
	http.HandleFunc("/info", info)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Hands on lv3</title>
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

func info(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("jwtToken")
	if err != nil {
		html := buildNotLoggedInHtml()
		fmt.Fprint(w, html)
		return
	}

	claims, err := parseToken(c.Value)
	if err != nil {
		html := buildErrorHtml(err.Error())
		fmt.Fprint(w, html)
		return
	}

	html := `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Hands on lv3</title>
		</head>
		<body>
			<p>You logged in as: ` + claims.Email + `</p>
		</body>
		</html>`
	fmt.Fprint(w, html)
}

func buildNotLoggedInHtml() string {
	return `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Hands on lv3</title>
	</head>
	<body>
		<p>You are not logged in yet</p>
	</body>
	</html>`
}

func buildErrorHtml(err string) string {
	return `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Hands on lv3</title>
	</head>
	<body>
		<p>` + err + `</p>
	</body>
	</html>`
}

func printResponseBody(body io.Reader) error {
	bs, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	log.Println(string(bs))
	return nil
}
