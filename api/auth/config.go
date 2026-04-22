package auth

import (
	"nlhosting/cfg"

	"golang.org/x/oauth2"
)

var config = &oauth2.Config{
	ClientID:     cfg.Config.OAuth2.Client.ID,
	ClientSecret: cfg.Config.OAuth2.Client.Secret,
	RedirectURL:  cfg.Config.OAuth2.Redirect,
	Scopes:       cfg.Config.OAuth2.Scopes,
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://www.nodeloc.com/oauth-provider/authorize",
		TokenURL: "https://www.nodeloc.com/oauth-provider/token",
	},
}
