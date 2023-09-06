package core

import (
	"gorm-web/conf"
	"gorm-web/pkg/log"
	"gorm-web/pkg/mysql"
	"gorm-web/web/dao"

	"gorm.io/gorm"
)

var (
	MysqlClientDemo *gorm.DB
)

func InitMysql() {
	var err error
	for name, dbConf := range conf.BasicConf.Mysql {
		switch name {
		case "demo":
			MysqlClientDemo, err = mysql.InitMysqlClient(dbConf)
			if err != nil {
				panic("mysql connect error: %v" + err.Error())
			}
			dao.InitDemoDao(MysqlClientDemo)
			log.Info("init mysql demo success")
		}
	}
}

func CloseMysql() {
}
