package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var GithubOauthConfig = &oauth2.Config{
	ClientID:     "8ecab7653a8b3804c2c8",
	ClientSecret: "12722aeebc07090cf0e084f9fdb0fb0ac0cac44a",
	Endpoint:     github.Endpoint,
}

type githubRes struct {
	Data struct {
		Viewer struct {
			Id string `json:"id"`
		} `json:"viewer"`
	} `json:"data"`
}
