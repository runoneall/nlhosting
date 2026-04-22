package host

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"nlhosting/api/auth"
	"nlhosting/api/user"
	"nlhosting/cfg"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length  = 10
)

func RequestNew(ctx *gin.Context) {
	payload := struct {
		Domain string `json:"domain"`
		Server string `json:"server"`
	}{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userInfo, _ := ctx.Get("user")
	domains, err := user.GetUserDomains(userInfo.(*auth.UserInfo).Email)
	if err != nil {
		slog.Error(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !slices.Contains(domains, payload.Domain) || !slices.Contains(cfg.Config.Server.Available, payload.Server) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, ok := cfg.Config.Server.Config[payload.Server]; !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	createNew(
		payload.Server,
		userInfo.(*auth.UserInfo).Email,
		payload.Domain,
	)

	ctx.Status(http.StatusOK)
}

func createNew(server, email, domain string) {
	conf := cfg.Config.Server.Config[server]
	total := cfg.Config.Server.Limit[server]

	RunCmd(conf, "froxlor-cli froxlor:api-call admin Customers.listingCount '{}'", func(stdout, stderr string, err error) {
		if err != nil {
			slog.Error(err.Error())
			SendMail(email, "未能完成创建请求", "由于一个后端错误，未能完成您的创建请求")
			return
		}

		resp := struct {
			Total int `json:"data"`
		}{}

		if err := json.Unmarshal([]byte(stdout), &resp); err != nil {
			slog.Error(err.Error())
			SendMail(email, "未能完成创建请求", "由于一个后端错误，未能完成您的创建请求")
			return
		}

		if resp.Total > total {
			SendMail(email, "未能完成创建请求", "由于服务器名额已满，未能完成您的创建请求")
			return
		}

		username := strings.ReplaceAll(domain, ".", "")
		password, err := newPassword()
		if err != nil {
			slog.Error(err.Error())
			SendMail(email, "未能完成创建请求", "由于一个后端错误，未能完成您的创建请求")
			return
		}

		query := fmt.Sprintf(`'{
	"email": "%s",
	"name": "Panel",
	"firstname": "Froxlor",
	"new_loginname": "%s",
	"new_customer_password": "%s",
	"createstdsubdomain": 0,
	"phpenabled": 1,
	"allowed_phpconfigs": 1,
	"logviewenabled": 1,
	"hosting_plan_id": 1
}'`, email, username, password)
		RunCmd(conf, "froxlor-cli froxlor:api-call admin Customers.add "+query, func(stdout, stderr string, err error) {
			if strings.Contains(stdout, username) {
				SendMail(email, "未能完成创建请求", stdout)
				return
			}

			query := fmt.Sprintf(`'{
  "domain": "%s",
  "loginname": "%s",
  "adminid": 1,
  "selectserveralias": 2,
  "caneditdomain": 1,
  "phpenabled": 1,
  "openbasedir": 1
}'`, domain, username)
			RunCmd(conf, "froxlor-cli froxlor:api-call admin Domains.add "+query, func(stdout, stderr string, err error) {
				if err != nil {
					slog.Error(err.Error(), "stdout", stdout, "stderr", stderr)
					SendMail(email, "未能完成创建请求", "由于一个后端错误，未能完成您的创建请求")
					return
				}

				SendMail(email, "创建成功", fmt.Sprintf(`登录信息：
用户名：%s
密码：%s
`, username, password))
			})
		})
	})
}

func newPassword() (string, error) {
	password := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}

		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
}
