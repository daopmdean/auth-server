package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func googleOauth(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New().String()
	oauthExp[uuid] = time.Now().Add(time.Hour)

	redirectUrl := googleOauthConfig.AuthCodeURL(uuid)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func googleOauthHandleReceive(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing google oauth receive...")

	state := r.FormValue("state")
	if time.Now().After(oauthExp[state]) {
		http.Error(w, "Login Time Expire", http.StatusInternalServerError)
		return
	}
	log.Println("Time still valid")

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

	jwtToken, err := createToken(gr.Email)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}
	log.Println(jwtToken)

	if _, ok := oauthConnections[gr.Id]; !ok {
		log.Println("user does not exist in system yet")
		oauthConnections[gr.Id] = *gr
	} else {
		log.Println("user exist in system")
	}
	c := http.Cookie{
		Name:  "jwtToken",
		Value: jwtToken,
		Path:  "/",
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/info", http.StatusSeeOther)
}
