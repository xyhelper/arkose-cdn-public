package main

import (
	"arkose-cdn/config"
	"arkose-cdn/handel"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	client = g.Client()
)

func main() {
	ctx := gctx.New()

	s := g.Server()
	s.SetPort(config.PORT)
	s.SetServerRoot("/resource/public")
	s.BindHandler("/*", handel.Proxy)
	s.BindHandler("/pushtoken", func(r *ghttp.Request) {
		token := r.Get("token").String()
		if token == "" {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "token is empty",
			})
			return
		}
		forwordURL := g.Cfg().MustGetWithEnv(ctx, "FORWORD_URL").String()
		g.Log().Info(ctx, "forwordURL", forwordURL)

		if forwordURL != "" {
			result := client.PostVar(ctx, forwordURL, g.Map{
				"token": token,
			})
			g.Log().Info(ctx, getRealIP(r), "forwordURL", forwordURL, result)
			r.Response.WriteJson(g.Map{
				"code": 1,
				"msg":  "success",
			})
			return
		} else {
			r.Response.WriteJson(g.Map{
				"code": 0,
				"msg":  "FORWORD_URL is empty",
			})
			return
		}
	})
	s.BindHandler("/ping", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"code":    1,
			"msg":     "pong",
			"headers": r.Header,
		})
		return
	})
	s.Run()
}
func getRealIP(req *ghttp.Request) string {
	// 优先获取Cf-Connecting-Ip
	if ip := req.Header.Get("Cf-Connecting-Ip"); ip != "" {
		return ip
	}

	// 优先获取X-Real-IP
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// 其次获取X-Forwarded-For
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	// 最后获取RemoteAddr
	ip := req.RemoteAddr
	// 处理端口
	if index := strings.Index(ip, ":"); index != -1 {
		ip = ip[0:index]
	}
	if ip == "[" {
		ip = req.GetClientIp()
	}
	return ip
}
