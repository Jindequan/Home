package dao

import (
	"bigdata_permission/conf"
	"bigdata_permission/pkg/logging"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

type Base struct {
	CreatedAt int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
	//DeletedAt  gorm.DeletedAt
}

func InitConns() {

	db, err := gorm.Open(mysql.Open(conf.BaseConf.DBConfig.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "t_",   // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})
	if err != nil {
		logging.Fatal(errors.New("数据库连接错误: " + err.Error()))
	}

	DB = db
}

// 看情况决定是否使用DB迁移

func Init() {
	InitConns()
}

//func SwitchConn(name string) error {
//	conn, ok := Conns[name]
//	if !ok {
//		return errors.New("不存在的DB名称 请检查")
//	}
//
//	DB = conn
//	return nil
//}
//
//func GetConn(name string) *gorm.DB {
//	return Conns[name]
//}
