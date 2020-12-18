package routes

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger/example/basic/docs"
)

func routeListForWeb(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("templates/*")
	r.StaticFS("static", http.Dir("static"))
	r.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	url := ginSwagger.URL("http://localhost:8081/swag/swagger.json")//这样写死的目的是让线上无法访问swagger页面
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.StaticFile("/swag/swagger.json", "./docs/swagger.json")

	return r
}
