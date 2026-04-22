package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type UserInfo struct {
	Sub               string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
	Picture           string `json:"picture"`
	ID                int    `json:"id"`
	Username          string `json:"username"`
	AvatarURL         string `json:"avatar_url"`
	TrustLevel        int    `json:"trust_level"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
}

func fetchUserInfo(accessToken string) (*UserInfo, error) {
	var userInfo UserInfo

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		SetResult(&userInfo).
		Get("https://www.nodeloc.com/oauth-provider/userinfo")

	if err != nil {
		return nil, fmt.Errorf("userinfo 请求失败 (状态码 %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	return &userInfo, nil
}

func CallbackHandler(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	errorParam := ctx.Query("error")

	session := sessions.Default(ctx)
	savedState := session.Get("oauth_state")

	if errorParam != "" {
		ctx.String(http.StatusBadRequest, "授权错误: %s", ctx.Query("error_description"))
		return
	}

	if savedState == nil || state != savedState.(string) {
		ctx.String(http.StatusBadRequest, "Invalid state parameter")
		return
	}

	session.Delete("oauth_state")
	_ = session.Save()

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "换取 token 失败: %v", err)
		return
	}

	userInfo, err := fetchUserInfo(token.AccessToken)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "获取用户信息失败: %v", err)
		return
	}

	userJSON, _ := json.Marshal(userInfo)
	session.Set("user", string(userJSON))

	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, "保存用户信息失败")
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}
