package controller

import (
	"gorm-web/entity/dto"
	"gorm-web/pkg/base"
	"gorm-web/pkg/log"
	"gorm-web/web/service"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(ctx *gin.Context) {
	params := &dto.UserRegister{}
	if err := ctx.ShouldBind(params); err != nil {
		log.Warn("controller[Antispam] param bind error[%s]", err.Error())
		base.RenderJsonFail(ctx, err)
		return
	}

	err := service.User.Register(ctx, params)
	if err != nil {
		base.RenderJsonFail(ctx, err)
		return
	}

	base.RenderJsonSucc(ctx, nil)
}

// Query 查询用户信息
func Query(ctx *gin.Context) {
	params := &dto.UserQuery{}
	if err := ctx.ShouldBind(params); err != nil {
		log.Warn("controller[Antispam] param bind error[%s]", err.Error())
		base.RenderJsonFail(ctx, err)
		return
	}

	data, err := service.User.UserQuery(ctx, params)
	if err != nil {
		base.RenderJsonFail(ctx, err)
		return
	}

	base.RenderJsonSucc(ctx, data)
}

// Update 修改用户信息
func Update(ctx *gin.Context) {
	params := &dto.UserUpdate{}
	if err := ctx.ShouldBind(params); err != nil {
		log.Warn("controller[Antispam] param bind error[%s]", err.Error())
		base.RenderJsonFail(ctx, err)
		return
	}

	err := service.User.UserUpdate(ctx, params)
	if err != nil {
		base.RenderJsonFail(ctx, err)
		return
	}

	base.RenderJsonSucc(ctx, nil)
}
