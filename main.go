package main

import (
	"net/http"
	"nlhosting/api"
	"nlhosting/cfg"
	"nlhosting/www"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte(cfg.Config.Cookie.Secret))
	store.Options(sessions.Options{
		MaxAge:   0,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("session", store))

	api.ServeAPI(r)
	www.ServeWWW(r)

	r.RunTLS(cfg.Config.Serve, "cert.pem", "key.pem")
}
