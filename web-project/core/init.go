package core

import (
	"fmt"
	"gorm-web/conf"
	"gorm-web/pkg/log"
)

func PreInit() {
	// 1.init config
	conf.InitConf()

	// 2.init logger
	if err := log.InitLogger(&conf.BasicConf.Log); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	// 3. init mysql
	InitMysql()

	// 4. init redis
	//InitRedis()
}
