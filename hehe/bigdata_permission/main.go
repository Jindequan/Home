package main

import (
	"github.com/facebookgo/grace/gracehttp"
	"bigdata_permission/dao"
	"bigdata_permission/pkg"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/pkg/logging"
	"bigdata_permission/routes"
	"net/http"
)

func init() {
	pkg.GetEnv()
	logging.Init()
	ecode.Init()
	dao.Init()

	//util.Init()
}

// @title xxx项目
// @version x.x
// @description xxx项目描述

// @contact.name 联系人名称
// @contact.email 联系人电话
// @host localhost:8081
func main() {

	server := routes.New()

	if err := gracehttp.Serve(&http.Server{Addr: server.Addr, Handler: server.Handler}); err != nil {
		logging.Fatal("http服务启动失败：" + err.Error())
	}
}
