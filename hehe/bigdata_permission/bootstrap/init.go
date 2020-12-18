package bootstrap

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/pkg/logging"
)

//抽象这一层的原因 不只main.go需要初始化这些东西 可能cli和单元测试也需要 做一个封装聚合~

//可以做启动前的检查和变量初始化的工作
func Init() {
	pkg.GetEnv()
	logging.Init()
	ecode.Init()
	dao.Init()
}
