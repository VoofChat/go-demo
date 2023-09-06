package dao

import (
	"context"
	"demo4/dao/base"
	"demo4/po"
	"fmt"
	"gorm.io/gorm"
)

type (
	userDao interface {
		base.BaseModel // interface接口嵌套
		GetByUserId(ctx context.Context, id uint64) (po.User, error)
		GetByUsername(ctx context.Context, username string) (po.User, error)
		Save(ctx context.Context, data *po.User) (err error)
		UpdateUsernameById(ctx context.Context, id uint64, username string) (rowAffects int64, err error)
	}

	defaultUserDao struct {
		base.BaseModel // 嵌入interface
	}
)

func NewUserDao(db *gorm.DB) userDao {
	return &defaultUserDao{
		base.NewBaseModel(db, po.UserTableName),
	}
}

func (d *defaultUserDao) GetByUserId(ctx context.Context, id uint64) (user po.User, err error) {
	err = d.GetByCond(ctx, &user, base.Where("user_id = ?", id))
	if err != nil {
		return
	}
	return
}

func (d *defaultUserDao) GetByUsername(ctx context.Context, username string) (user po.User, err error) {
	err = d.GetByCond(ctx, &user, base.Where("username = ?", username))
	if err != nil {
		return
	}
	return
}

func (d *defaultUserDao) UpdateUsernameById(ctx context.Context, id uint64, username string) (rowAffects int64, err error) {
	if len(username) == 0 {
		return 0, fmt.Errorf("username is empty")
	}

	return d.Update(ctx, po.User{UserName: username}, base.Where("user_id = ?", id))
}

func (d *defaultUserDao) Save(ctx context.Context, data *po.User) (err error) {
	_, err = d.Create(ctx, data)
	if err != nil {
		return
	}
	return
}
