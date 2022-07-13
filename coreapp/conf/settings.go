package conf

import (
	"golang.org/x/oauth2"
)

type AuthApiSocial struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	RedirectURL  string
}

var url_server = "http://localhost:9000/"

var EndpointFacebook = oauth2.Endpoint{
	AuthURL:  "https://www.facebook.com/v12.0/dialog/oauth",
	TokenURL: "https://graph.facebook.com/v12.0/oauth/access_token",
}

var FACEBOOK = &oauth2.Config{
	ClientID:     "1964421813758740",
	ClientSecret: "e7abc5a6a834badc94cfc15c10240663",
	Scopes:       []string{},
	Endpoint:     EndpointFacebook,
	RedirectURL:  url_server + "Application/Auth",
}

var FacebookAuth = AuthApiSocial{
	ClientID:     "1964421813758740",
	ClientSecret: "e7abc5a6a834badc94cfc15c10240663",
	AuthURL:      "https://www.facebook.com/v12.0/dialog/oauth",
	TokenURL:     "https://graph.facebook.com/v12.0/oauth/access_token",
	RedirectURL:  url_server + "Application/Auth",
}
