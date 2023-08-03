package config

import (
	"net/url"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	PORT = 8080
)

func init() {
	ctx := gctx.GetInitCtx()
	port := g.Cfg().MustGetWithEnv(ctx, "PORT").Int()
	if port != 0 {
		PORT = port
	}

}

func PROXY(ctx g.Ctx) *url.URL {
	proxy := g.Cfg().MustGetWithEnv(ctx, "PROXY").String()
	// g.Log().Infof(ctx, "PROXY: %s", proxy)
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	return proxyURL
}
