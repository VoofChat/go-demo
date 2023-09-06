package po

import "time"

type Blog struct {
	BlogID      int64     `gorm:"primary_key;AUTO_INCREMENT;column:blog_id;type:bigint;" json:"blog_id"`          // 主键
	UserID      int64     `gorm:"column:user_id;type:bigint;" json:"user_id"`                                     // 用户id
	Title       string    `gorm:"column:title;type:varchar;size:255;" json:"title"`                               // 标题
	Description string    `gorm:"column:description;type:text;size:65535;" json:"description"`                    // 博客描述
	Content     string    `gorm:"column:content;type:text;size:4294967295;" json:"content"`                       // 内容
	Recommend   bool      `gorm:"column:recommend;type:bit;default:b'0';" json:"recommend"`                       // 是否推荐
	Deleted     int       `gorm:"column:deleted;type:int;default:0;" json:"deleted"`                              // -1 删除
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;" json:"create_time"` // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;" json:"update_time"` // 更新时间
}

const BlogTableName = "tbl_blog"

func (Blog) TableName() string {
	return BlogTableName
}
