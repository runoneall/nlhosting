//go:build prod

package www

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var distFS embed.FS

func ServeWWW(r *gin.Engine) {
	staticFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatalln(err)
	}

	r.GET("/assets/*filepath", func(ctx *gin.Context) {
		filepath := strings.TrimPrefix(ctx.Request.URL.Path, "/")
		http.ServeFileFS(ctx.Writer, ctx.Request, staticFS, filepath)
	})

	r.GET("/favicon.ico", func(ctx *gin.Context) {
		http.ServeFileFS(ctx.Writer, ctx.Request, staticFS, "favicon.ico")
	})

	r.NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/") {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		http.ServeFileFS(ctx.Writer, ctx.Request, staticFS, "index.html")
	})
}
