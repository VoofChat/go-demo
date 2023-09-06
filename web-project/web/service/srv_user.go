package service

import (
	"fmt"
	"gorm-web/entity/dto"
	"gorm-web/entity/po"
	"gorm-web/web/dao"

	"github.com/gin-gonic/gin"
)

var User userSrv = new(defaultUserSrv)

type (
	userSrv interface {
		// Register 用户注册
		Register(ctx *gin.Context, req *dto.UserRegister) (err error)
		UserQuery(ctx *gin.Context, req *dto.UserQuery) (UserQuery po.User, err error)
		UserUpdate(ctx *gin.Context, req *dto.UserUpdate) (err error)
	}
	defaultUserSrv struct{}
)

// Register 用户注册
func (d *defaultUserSrv) Register(ctx *gin.Context, req *dto.UserRegister) (err error) {
	//zap.Info("Register Success, req:%#v", req)

	UserQuery, err := dao.UserDao.GetByUsername(ctx, req.Username)
	if err != nil {
		return
	}

	//UserQuery1, _ := dao.UserDao.GetByUserId(ctx, UserQuery.UserID)
	//zap.Info("UserQuery1: %#v", UserQuery1)

	if UserQuery.UserID > 0 {
		err = fmt.Errorf("user already exists")
		return
	}

	user := po.User{
		UserName: req.Username,
		Password: req.Password,
	}
	err = dao.UserDao.Save(ctx, user)
	if err != nil {
		return
	}
	return nil
}

func (d *defaultUserSrv) UserQuery(ctx *gin.Context, req *dto.UserQuery) (UserQuery po.User, err error) {
	UserQuery, err = dao.UserDao.GetByUserId(ctx, req.UserId)
	if err != nil {
		return
	}
	if UserQuery.UserID == 0 {
		err = fmt.Errorf("user not exists")
		return
	}
	return
}

func (d *defaultUserSrv) UserUpdate(ctx *gin.Context, req *dto.UserUpdate) (err error) {
	userInfo, err := dao.UserDao.GetByUserId(ctx, req.UserId)
	if err != nil {
		return
	}
	if userInfo.UserID == 0 {
		err = fmt.Errorf("user not exists")
		return
	}

	userInfo.UserName = req.Username
	_, err = dao.UserDao.Update(ctx, userInfo)
	if err != nil {
		return
	}
	return
}
