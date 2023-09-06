package dao

import (
	"gorm-web/entity/po"
	"gorm-web/pkg/base"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type (
	blogDao interface {
		base.BaseModel
		GetByUserId(ctx *gin.Context, userId int64) ([]po.Blog, error)
		Save(ctx *gin.Context, data po.Blog) (err error)
	}

	defaultBlogDao struct {
		base.BaseModel // 嵌入interface
	}
)

func newBlogDao(db *gorm.DB) blogDao {
	return &defaultBlogDao{
		base.NewBaseModel(db, po.BlogTableName),
	}
}

func (d defaultBlogDao) GetByUserId(ctx *gin.Context, userId int64) ([]po.Blog, error) {
	panic("implement me")
}

func (d defaultBlogDao) Save(ctx *gin.Context, data po.Blog) (err error) {
	panic("implement me")
}
