package user

import (
	"fmt"
	"nlhosting/cfg"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	tokenCache     string
	tokenCacheTime time.Time
	cacheMutex     sync.RWMutex
)

func GetToken() (string, error) {
	cacheMutex.RLock()
	if tokenCache != "" && time.Since(tokenCacheTime) < 5*time.Minute {
		cached := tokenCache
		cacheMutex.RUnlock()
		return cached, nil
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if tokenCache != "" && time.Since(tokenCacheTime) < 5*time.Minute {
		return tokenCache, nil
	}

	client := resty.New()
	result := struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{}

	resp, err := client.R().
		SetBody(map[string]string{"email": cfg.Config.APIToken.Email, "password": cfg.Config.APIToken.Password}).
		SetResult(&result).
		Post("https://domain.nodeloc.com/api/auth/login")

	if err != nil {
		return "", err
	}

	if result.Message != "success.login" {
		return "", fmt.Errorf("api token 请求失败 (状态码 %d): %s", resp.StatusCode(), resp.String())
	}

	tokenCache = result.Token
	tokenCacheTime = time.Now()
	return tokenCache, nil
}
