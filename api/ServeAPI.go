package api

import (
	"net/http"
	"nlhosting/api/auth"
	"nlhosting/api/host"
	"nlhosting/api/user"

	"github.com/gin-gonic/gin"
)

func ServeAPI(r *gin.Engine) {
	g := r.Group("/api")

	g.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	{
		r := g.Group("/auth")
		r.GET("/login", auth.LoginHandler)
		r.GET("/callback", auth.CallbackHandler)
		r.GET("/logout", auth.LogoutHandler)
	}

	{
		r := g.Group("/user")
		r.Use(auth.AuthMiddleware())
		r.GET("/info", user.UserInfo)
		r.GET("/domains", user.GetUserDomain)
	}

	{
		r := g.Group("/host")
		r.Use(auth.AuthMiddleware())
		r.GET("/available", host.GetAvailable)
		r.POST("/new", host.RequestNew)
	}
}
