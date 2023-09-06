package po

import "time"

type User struct {
	UserID     int64     `gorm:"primary_key;AUTO_INCREMENT;column:user_id;type:bigint;" json:"user_id"`          // 主键
	UserName   string    `gorm:"column:username;type:varchar;size:255;" json:"username"`                         // 用户名
	Password   string    `gorm:"column:password;type:varchar;size:255;" json:"password"`                         // 密码
	Deleted    int       `gorm:"column:deleted;type:int;default:0;" json:"deleted"`                              // -1 删除
	CreateTime time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;" json:"update_time"` // 更新时间
}

const UserTableName = "tbl_user"

func (User) TableName() string {
	return UserTableName
}
