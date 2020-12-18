package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/gin-swagger/example/basic/docs"
	"bigdata_permission/conf"
	"bigdata_permission/middleware"
	"bigdata_permission/pkg/logging"
	"io"
	"net/http"
	"os"
	"time"
)

func SetRoutes(r *gin.Engine) *gin.Engine {
	defer func() { //捕获可能因为路由定义引起的错误
		if err := recover(); err != nil {
			logging.Fatal("路由定义错误：" + err.(string))
		}
	}()
	r = routeListForApi(r) //API路由定义
	r = routeListForWeb(r) //前端路由定义
	return r
}

func New() *http.Server {
	gin.DefaultWriter = io.MultiWriter(logging.F, os.Stdout)
	// init static file handler
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	r = SetRoutes(r) //路由定义
	//r.NoRoute(func(context *gin.Context) { //未知路由反404
	//	context.String(404, "No route")
	//})

	server := &http.Server{
		Addr:         conf.BaseConf.HTTPServerConfig.Addr,
		Handler:      r,
		ReadTimeout:  time.Duration(conf.BaseConf.HTTPServerConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.BaseConf.HTTPServerConfig.WriteTimeout) * time.Second,
	}

	return server
}
