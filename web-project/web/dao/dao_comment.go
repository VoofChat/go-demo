package dao

import (
	"gorm-web/entity/po"
	"gorm-web/pkg/base"

	"gorm.io/gorm"
)

type (
	commentDao interface {
		base.BaseModel
		GetByUserId(userId int64) ([]po.Comment, error)
		Save(data po.Comment) (err error)
	}

	defaultCommentDao struct {
		base.BaseModel // 嵌入interface
	}
)

func newCommentDao(db *gorm.DB) commentDao {
	return &defaultCommentDao{
		base.NewBaseModel(db, po.CommentTableName),
	}
}

func (d defaultCommentDao) GetByUserId(userId int64) ([]po.Comment, error) {
	panic("implement me")
}

func (d defaultCommentDao) Save(data po.Comment) (err error) {
	panic("implement me")
}
