package user

import (
	"log/slog"
	"net/http"
	"nlhosting/api/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func GetUserDomains(userEmail string) ([]string, error) {
	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	result := struct {
		Domains []struct {
			FullDomain string `json:"full_domain"`
			Status     string `json:"status"`
		} `json:"domains"`
	}{}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetQueryParam("email", userEmail).
		SetResult(&result).
		Get("https://domain.nodeloc.com/api/admin/domains")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		slog.Error(resp.String())
		return nil, err
	}

	var activeDomains []string
	for _, domain := range result.Domains {
		if domain.Status == "active" {
			activeDomains = append(activeDomains, domain.FullDomain)
		}
	}

	return activeDomains, nil
}

func GetUserDomain(ctx *gin.Context) {
	userInfo, _ := ctx.Get("user")
	domains, err := GetUserDomains(userInfo.(*auth.UserInfo).Email)
	if err != nil {
		slog.Error(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, domains)
}
