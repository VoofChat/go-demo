package main

import (
	"fmt"
	"gorm-web/conf"
	"gorm-web/core"
	"gorm-web/pkg/log"
	"gorm-web/web"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1.初始化基本资源
	core.PreInit()

	// 2.gin
	//gin.SetMode("debug")
	r := gin.New()
	r.Use(log.GinLogger(), log.GinRecovery(true)) // 注册zap相关中间件

	// 3 注册路由
	web.Http(r)

	// 4.启动服务
	addr := fmt.Sprintf(":%v", conf.BasicConf.Server.Port)
	fmt.Printf("start http server at %s\n", addr)
	if err := r.Run(addr); err != nil {
		fmt.Printf("init http server failed, err:%v\n", err)
		return
	}

}
