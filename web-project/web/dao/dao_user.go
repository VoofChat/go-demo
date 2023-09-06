package dao

import (
	"gorm-web/entity/po"
	"gorm-web/pkg/base"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type (
	userDao interface {
		base.BaseModel
		GetByUserId(ctx *gin.Context, id uint64) (po.User, error)
		GetByUsername(ctx *gin.Context, username string) (po.User, error)
		Save(ctx *gin.Context, data po.User) (err error)
	}

	defaultUserDao struct {
		base.BaseModel // 嵌入interface
	}
)

func newUserDao(db *gorm.DB) userDao {
	return &defaultUserDao{
		base.NewBaseModel(db, po.UserTableName),
	}
}

func (d *defaultUserDao) GetByUserId(ctx *gin.Context, id uint64) (user po.User, err error) {
	err = d.GetByCond(ctx, &user, base.Where("user_id = ?", id))
	if err != nil {
		return
	}
	return
}

func (d *defaultUserDao) GetByUsername(ctx *gin.Context, username string) (user po.User, err error) {
	err = d.GetByCond(ctx, &user, base.Where("username = ?", username))
	if err != nil {
		return
	}
	return
}

func (d *defaultUserDao) Save(ctx *gin.Context, data po.User) (err error) {
	_, err = d.Create(ctx, &data)
	if err != nil {
		return
	}
	return
}
