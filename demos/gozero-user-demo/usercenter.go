package main

import (
	"flag"
	"fmt"

	"gozero-user-demo/internal/config"
	"gozero-user-demo/internal/handler"
	"gozero-user-demo/internal/responsex"
	"gozero-user-demo/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/usercenter-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	httpx.SetOkHandler(responsex.OkHandler)
	httpx.SetErrorHandlerCtx(responsex.ErrorHandler)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
