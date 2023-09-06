package web

import (
	"gorm-web/web/controller"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Http(engine *gin.Engine) {
	router := engine.Group("/demo")

	router.GET("/hello", func(c *gin.Context) {
		// 4. 初始化dao
		// 假设你有一些数据需要记录到日志中
		var (
			name = "q1mi"
			age  = 18
		)
		// 记录日志并使用zap.Xxx(key, val)记录相关字段
		zap.L().Debug("this is hello func", zap.String("user", name), zap.Int("age", age))
		c.String(http.StatusOK, "hello")
	})

	router.POST("/user/register", controller.Register)
	router.POST("/user/query", controller.Query)
	router.POST("/user/update", controller.Update)
}
