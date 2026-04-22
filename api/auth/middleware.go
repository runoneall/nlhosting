package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userJSON := session.Get("user")

		if userJSON == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user UserInfo
		if err := json.Unmarshal([]byte(userJSON.(string)), &user); err != nil {
			session.Clear()
			_ = session.Save()

			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", &user)
		ctx.Next()
	}
}
