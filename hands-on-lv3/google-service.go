package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func googleOauth(w http.ResponseWriter, r *http.Request) {
	redirectUrl := googleOauthConfig.AuthCodeURL("8888")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func googleOauthHandleReceive(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing google oauth receive...")

	state := r.FormValue("state")
	if state != "8888" {
		http.Error(w, "Invalid State", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Could not login", http.StatusInternalServerError)
		return
	}

	src := googleOauthConfig.TokenSource(ctx, token)
	client := oauth2.NewClient(ctx, src)

	res, err := client.Get(googleUserInfoApi)
	if err != nil {
		http.Error(w, "Could not get google info", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	gr := &googleRes{}
	json.NewDecoder(res.Body).Decode(gr)
	log.Println("ID:", gr.Id)
	log.Println("Email:", gr.Email)
	log.Println("Name:", gr.Name)
}
