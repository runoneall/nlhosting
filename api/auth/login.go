package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func genState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func LoginHandler(ctx *gin.Context) {
	state, err := genState()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "生成 state 失败")
		return
	}

	session := sessions.Default(ctx)
	session.Set("oauth_state", state)

	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, "保存 session 失败")
		return
	}

	url := config.AuthCodeURL(state)
	ctx.Redirect(http.StatusFound, url)
}
