package main

import (
	"flag"
	"fmt"
	"todo/app/user/api/internal/config"
	"todo/app/user/api/internal/handler"
	"todo/app/user/api/internal/middleware"
	"todo/app/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 全局安全中间件：安全响应头 + CSRF 双重提交校验
	// 当前 go-zero 版本的 Use 只接收一个 middleware 参数，因此分两次调用
	server.Use(middleware.SecurityHeadersMiddleware)
	server.Use(middleware.CSRFMiddleware)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}
