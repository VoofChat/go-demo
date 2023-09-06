package po

import "time"

type Comment struct {
	CommentID  int64     `gorm:"primary_key;AUTO_INCREMENT;column:comment_id;type:bigint;" json:"comment_id"`    // 主键
	BlogID     int64     `gorm:"column:blog_id;type:bigint;" json:"blog_id"`                                     // 博客id
	UserID     int64     `gorm:"column:user_id;type:bigint;" json:"user_id"`                                     // 用户id
	Content    string    `gorm:"column:content;type:varchar;size:255;" json:"content"`                           // 评论内容
	Deleted    int       `gorm:"column:deleted;type:int;default:0;" json:"deleted"`                              // -1 删除
	CreateTime time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;" json:"update_time"` // 更新时间
}

const CommentTableName = "tbl_comment"

func (Comment) TableName() string {
	return CommentTableName
}
