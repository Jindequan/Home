package common

import (
	"bigdata_permission/pkg/dictionary"
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDictionary (c *gin.Context) {
	appG := &http2.Gin{c}
	dic := struct {
		RoleDic []*dictionary.GenericItem `json:"role_dic"`
		ModuleDic []*dictionary.GenericItem `json:"module_dic"`
		ModuleTree []*dictionary.GenericTreeNode `json:"module_tree"`
		InterfaceDic []*dictionary.GenericItem `json:"interface_dic"`
	}{
		RoleDic: dictionary.RoleList,
		ModuleDic: dictionary.ModuleList,
		ModuleTree: dictionary.ModuleTree,
		InterfaceDic: dictionary.InterfaceList,
	}
	appG.JSON(http.StatusOK, ecode.OK, dic)
}
