package dao

import (
	"gorm.io/gorm"
)

var (
	BlogDao    blogDao
	CommentDao commentDao
	UserDao    userDao
)

// InitDemoDao 初始化demo dao
func InitDemoDao(db *gorm.DB) {
	BlogDao = newBlogDao(db)
	CommentDao = newCommentDao(db)
	UserDao = newUserDao(db)
}
