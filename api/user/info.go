package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, user)
}
