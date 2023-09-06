package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

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

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("BeforeCreate")
	if len(u.UserName) == 0 {
		return fmt.Errorf("invalid username")
	}
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("AfterCreate")
	return
}

var db *gorm.DB

// gormLogger gorm 日志配置 https://gorm.io/docs/logger.html
func gormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
}

func init() {
	// MySQL 配置信息
	dsn := fmt.Sprintf("root:qwerasdf@tcp(127.0.0.1:3306)/demo?charset=utf8&parseTime=True&loc=Local&timeout=10s")
	// Open 连接
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		//CreateBatchSize: 2, // 批量插入时，每次插入的条数
	})
	//gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLogger()})
	if err != nil {
		panic("failed to connect mysql.")
	}
	// 调试单个操作，显示此操作的详细日志
	//err = gormDB.Debug().Where("username = ?", "jinzhu").Find(&[]User{}).Error
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println("init db success")
	db = gormDB
}

func jsonMarshal(obj interface{}) interface{} {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return obj
	}
	return string(jsonStr)
}

// https://gorm.io/docs/create.html#Create-Hooks
// insert 新增数据
func insert1() {
	user := User{
		UserName: "insert1",
		Password: "123456",
	}
	// 插入后返回的常用数据
	// user.UserID         插入的主键值
	// result.ERROR        返回的 error
	// result.RowsAffected 插入的条数
	result := db.Create(&user)
	// INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('insert1','123456',0)
	if result.Error != nil {
		fmt.Println("insert1 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("insert1 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

func insert2() {
	user := User{
		UserName: "insert2",
		Password: "123456",
	}
	// 只插入指定字段
	result := db.Select("UserName").Create(&user)
	// INSERT INTO `tbl_user` (`username`) VALUES ('insert2')
	if result.Error != nil {
		fmt.Println("insert2 failed. ", result.Error) //  Error 1364 (HY000): Field 'password' doesn't have a default value
		return
	}
	fmt.Println(fmt.Sprintf("insert2 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

func insert3() {
	user := User{
		UserName: "insert3",
		Password: "123456",
	}
	// 不插入指定字段
	result := db.Omit("UserName").Create(&user)
	// INSERT INTO `tbl_user` (`password`,`deleted`) VALUES ('123456',0)
	if result.Error != nil {
		fmt.Println("insert3 failed. ", result.Error) // Error 1364 (HY000): Field 'username' doesn't have a default value
		return
	}
	fmt.Println(fmt.Sprintf("insert3 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

func insert4() {
	users := []User{
		{
			UserName: "insert4_u1",
			Password: "123456",
		},
		{
			UserName: "insert4_u2",
			Password: "123456",
		},
		{
			UserName: "insert4_u3",
			Password: "123456",
		},
	}
	result := db.Create(&users)
	// INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('insert4_u1','123456',0),('insert4_u2','123456',0)
	if result.Error != nil {
		fmt.Println("insert4 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("insert4 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(users)))
}

func insert5() {
	users := []User{
		{
			UserName: "insert5_u1",
			Password: "123456",
		},
		{
			UserName: "insert5_u2",
			Password: "123456",
		},
		{
			UserName: "insert5_u3",
			Password: "123456",
		},
		{
			UserName: "insert5_u4",
			Password: "123456",
		},
	}
	result := db.CreateInBatches(&users, 2)
	// INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('insert5_u1','123456',0),('insert5_u2','123456',0)
	// INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('insert5_u3','123456',0),('insert5_u4','123456',0)
	if result.Error != nil {
		fmt.Println("insert4 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("insert5 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(users)))
}

// Upsert / On Conflict
func upsert1() {
	// Do nothing on conflict
	user := User{
		UserName: "insert1",
		Password: "123456",
	}
	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
	// [66.931ms] [rows:0] INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('insert1','123456',0) ON DUPLICATE KEY UPDATE `user_id`=`user_id`
	if result.Error != nil {
		fmt.Println("upsert1 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("upsert1 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

func upsert2() {
	user := User{
		UserName: "insert1",
		Password: "123456",
	}
	// 如果出现冲突，则更新指定字段Password值为"789abc"
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"Password": "789abc"}),
	}).Create(&user)
	// [102.821ms] [rows:2] INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('insert1','123456',0) ON DUPLICATE KEY UPDATE `Password`='789abc'

	if result.Error != nil {
		fmt.Println("upsert2 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("upsert2 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

func upsert3() {
	user := User{
		UserName: "insert1",
		Password: "aaaaaa",
		Deleted:  1,
	}
	// 如果出现冲突，则更新指定字段Password值为"aaaaaa", Deleted值为1
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"Password", "Deleted"}),
	}).Create(&user)
	// [88.379ms] [rows:2] INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('insert1','aaaaaa',1) ON DUPLICATE KEY UPDATE `Password`=VALUES(`Password`),`Deleted`=VALUES(`Deleted`)
	if result.Error != nil {
		fmt.Println("upsert3 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("upsert3 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

// update 数据更新
// https://gorm.io/docs/update.html

// update1 更新所有字段
func update1() {
	// 注意零值问题
	user := User{
		UserID:   2020062500070,
		UserName: "update1",
	}
	err := db.Save(&user).Error
	// [1.563ms] [rows:0] UPDATE `tbl_user` SET `username`='update1',`password`='',`deleted`=0,`create_time`='0000-00-00 00:00:00',`update_time`='0000-00-00 00:00:00' WHERE `user_id` = 2020062500070
	if err != nil {
		fmt.Println("update1 failed. ", err)
		return
	}
}

// update2 更新单列
func update2() {
	var user User
	user.UserID = 2020062500070
	result := db.Model(&user).Update("username", "update2")
	// [41.281ms] [rows:1] UPDATE `tbl_user` SET `username`='update2' WHERE `user_id` = 2020062500070
	if result.Error != nil {
		fmt.Println("update2 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("update2 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

// update3 更新多列
func update3() {
	user := User{
		UserID:   2020062500070,
		UserName: "update3",
		Password: "bbbbbb",
	}
	result := db.Model(&user).Updates(User{UserName: user.UserName, Password: user.Password})
	// [1.272ms] [rows:0] UPDATE `tbl_user` SET `username`='update3',`password`='bbbbbb' WHERE `user_id` = 2020062500070

	//result := db.Model(&user).Updates(user)
	// [82.348ms] [rows:1] UPDATE `tbl_user` SET `user_id`=2020062500070,`username`='update3',`password`='bbbbbb' WHERE `user_id` = 2020062500070

	//result := db.Updates(user)
	// [1.018ms] [rows:0] UPDATE `tbl_user` SET `user_id`=2020062500070,`username`='update3',`password`='bbbbbb' WHERE `user_id` = 2020062500070

	if result.Error != nil {
		fmt.Println("update3 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("update3 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

// 更新指定列 （推荐）
func update4() {
	var user User
	err := db.Where("user_id = ?", 2020062500070).First(&user).Error
	if err != nil {
		fmt.Println("update4 failed. ", err)
		return
	}

	// 更新Username列
	//result := db.Model(&user).Select("Username").Updates(map[string]interface{}{"username": user.UserName, "password": user.Password})
	// [0.546ms] [rows:0] UPDATE `tbl_user` SET `username`='update3' WHERE `user_id` = 2020062500070

	// 更新Username列
	//result := db.Model(&user).Select("Username").Updates(&user)
	// UPDATE `tbl_user` SET `username`='update3' WHERE `user_id` = 2020062500070

	// 更新除"CreateTime", "deleted"之外的列
	//result := db.Model(&user).Omit("CreateTime", "deleted").Updates(&user)
	// UPDATE `tbl_user` SET `username`='update3',`password`='bbbbbb',`update_time`='2023-08-29 18:05:32' WHERE `user_id` = 2020062500070

	// 更新所有除"CreateTime"之外的列
	result := db.Model(&user).Select("*").Omit("CreateTime").Updates(&user)
	// UPDATE `tbl_user` SET `username`='update3',`password`='bbbbbb',`deleted`=1,`update_time`='2023-08-29 18:05:32' WHERE `user_id` = 2020062500070

	if result.Error != nil {
		fmt.Println("update4 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("update4 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

// query https://gorm.io/docs/query.html
// query1 单条查询
func query1() {
	var user User

	result := db.Where("deleted = ?", 0).First(&user) // id 升序
	// SELECT * FROM `tbl_user` WHERE deleted = 0 ORDER BY `tbl_user`.`user_id` ASC LIMIT 1

	//result := db.Where("deleted = ?", 0).Take(&user) // 没有指定order
	// SELECT * FROM `tbl_user` WHERE deleted = 0 LIMIT 1

	//result := db.Where("deleted = ?", 0).Limit(1).Find(&user) // 同上
	// SELECT * FROM `tbl_user` WHERE deleted = 0 LIMIT 1

	//result := db.Where("deleted = ?", 0).Last(&user) // id 降序
	// SELECT * FROM `tbl_user` WHERE deleted = 0 ORDER BY `tbl_user`.`user_id` DESC LIMIT 1

	if result.Error != nil {
		fmt.Println("query1 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("query1 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

// query2 多条查询
func query2() {
	var users []User
	result := db.Where("deleted = ?", 0).Find(&users)
	// SELECT * FROM `tbl_user` WHERE deleted = 0
	if result.Error != nil {
		fmt.Println("query2 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("query2 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(users)))
}

// 主键查询
func query31() {
	var user User
	//result := db.First(&user, 2020062500070)
	// SELECT * FROM `tbl_user` WHERE `tbl_user`.`user_id` = 2020062500070 ORDER BY `tbl_user`.`user_id` ASC LIMIT 1

	result := db.First(&user, "2020062500070")
	// SELECT * FROM `tbl_user` WHERE `tbl_user`.`user_id` = '2020062500070' ORDER BY `tbl_user`.`user_id` LIMIT 1

	if result.Error != nil {
		fmt.Println("query31 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("query31 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

func query32() {
	user := User{UserID: 2020062500070}
	result := db.First(&user)
	// SELECT * FROM `tbl_user` WHERE `tbl_user`.`user_id` = '2020062500070' ORDER BY `tbl_user`.`user_id` LIMIT 1

	if result.Error != nil {
		fmt.Println("query32 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("query32 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(user)))
}

func query33() {
	var users []User
	result := db.Find(&users, []int{2020062500070, 2020062500017})
	// [1.512ms] [rows:2] SELECT * FROM `tbl_user` WHERE `tbl_user`.`user_id` IN (2020062500070,2020062500017)
	if result.Error != nil {
		fmt.Println("query33 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("query33 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(users)))
}

// 条件查询
func query4() {
	var users []User

	result := db.Where("deleted = ?", 0).Find(&users)
	// SELECT * FROM `tbl_user` WHERE deleted = 0

	//result := db.Where("username = ?", "zhangsan").Find(&users)
	//  SELECT * FROM `tbl_user` WHERE username = 'zhangsan'

	//result := db.Where("username <> ?", "zhangsan").Find(&users)
	//  SELECT * FROM `tbl_user` WHERE username <> 'zhangsan'

	//result := db.Where("username in ?", []string{"zhangsan"}).Find(&users)
	//  SELECT * FROM `tbl_user` WHERE username in ('zhangsan')

	//result := db.Where("username LIKE ?", "%zhang%").Find(&users)
	//  SELECT * FROM `tbl_user` WHERE username LIKE '%zhang%

	//result := db.Where("username = ? and deleted = ?", "zhangsan", 0).Find(&users)
	// SELECT * FROM `tbl_user` WHERE username = 'zhangsan' and deleted = 0

	//result := db.Where("username = ?", "zhangsan").Or("deleted = ?", 0).Find(&users)
	// SELECT * FROM `tbl_user` WHERE username = 'zhangsan' OR deleted = 0

	//result := db.Where("username = ?", "zhangsan").Order("user_id desc").Find(&users)
	// SELECT * FROM `tbl_user` WHERE username = 'zhangsan' ORDER BY user_id desc

	//result := db.Limit(2).Offset(2).Find(&users)
	// SELECT * FROM `tbl_user` LIMIT 2 OFFSET 2

	if result.Error != nil {
		fmt.Println("query33 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("query33 success. RowsAffected: %d, User:%v", result.RowsAffected, jsonMarshal(users)))

}

// delete
func delete1() {
	/*
		INSERT INTO `tbl_user` (`user_id`, `username`, `password`, `deleted`, `create_time`, `update_time`)
		VALUES
			(2020062500070, 'update3', 'bbbbbb', 1, '2023-08-29 17:32:56', '2023-08-29 18:05:32');
	*/

	//user := User{
	//	UserID: 2020062500070,
	//}
	//result := db.Delete(&user)
	// DELETE FROM `tbl_user` WHERE `tbl_user`.`user_id` = 2020062500070

	result := db.Where("user_id = ?", 2020062500070).Delete(&User{})
	// DELETE FROM `tbl_user` WHERE user_id = 2020062500070

	if result.Error != nil {
		fmt.Println("query33 failed. ", result.Error)
		return
	}
	fmt.Println(fmt.Sprintf("delete1 success. RowsAffected: %d", result.RowsAffected))
}

func main() {
	//insert1()
	//insert2()
	//insert3()
	//insert4()
	//insert5()
	//upsert1()
	//upsert2()
	//upsert3()

	//update1()
	//update2()
	//update3()
	//update4()

	//query1()
	//query2()
	//query31()
	//query32()
	//query33()
	//query4()

	//delete1()
}
