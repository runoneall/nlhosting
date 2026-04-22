package host

import (
	"net/http"
	"nlhosting/cfg"

	"github.com/gin-gonic/gin"
)

func GetAvailable(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, cfg.Config.Server.Available)
}
