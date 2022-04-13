package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func githubOauth(w http.ResponseWriter, r *http.Request) {
	redirectUrl := githubOauthConfig.AuthCodeURL("8888")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func githubOauthHandleReceive(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing github oauth receive...")

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
	defer res.Body.Close()

	u := &githubRes{}
	json.NewDecoder(res.Body).Decode(u)

	githubId := u.Data.Viewer.Id
	appUserId, ok := githubConnections[githubId]
	if !ok {
		// Create new account for example
		githubConnections[githubId] = uuid.New().String()
	}

	log.Println(appUserId)
	log.Println(githubConnections)
}
