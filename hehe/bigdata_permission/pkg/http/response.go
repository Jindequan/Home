package http

import (
	"bigdata_permission/pkg/ecode"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) JSON(httpCode int, errCode error, data interface{}) {
	code := ecode.Cause(ecode.OK)
	if errCode != nil {
		code = ecode.Cause(errCode)
	}
	g.C.JSON(httpCode, Response{
		Code: code.Code(),
		Msg:  code.Message(),
		Data: data,
	})
	return
}
