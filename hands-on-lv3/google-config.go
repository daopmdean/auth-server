package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	ClientID:     "167611725265-is1oj5liqvtr5afpn6t9cfddb4bi5a8r.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8080/oauth2/google/receive",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

type googleRes struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
