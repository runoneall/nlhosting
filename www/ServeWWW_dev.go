//go:build !prod

package www

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func ServeWWW(r *gin.Engine) {
	viteURL, _ := url.Parse("http://localhost:5173")
	proxy := httputil.NewSingleHostReverseProxy(viteURL)

	proxy.Director = func(req *http.Request) {
		req.Host = viteURL.Host
		req.URL.Scheme = viteURL.Scheme
		req.URL.Host = viteURL.Host
	}

	r.NoRoute(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		if strings.HasPrefix(path, "/api/") {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
}
