package test

import (
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 测试接口
// @Id 1
// @Tags 测试
// @version 1.0
// @Accept application/x-json-stream
// @Success 200
// @Router /test/auth [get]
func Test(c *gin.Context) {
	//或err := errors.New("0")
	appG := &http2.Gin{c}
	appG.JSON(http.StatusOK, ecode.OK, "你拥有访问权限")
	return
}

func TestNoAuth(c *gin.Context)  {
	appG := &http2.Gin{c}
	appG.JSON(http.StatusOK, ecode.OK, "任何人都可访问")
	return
}
