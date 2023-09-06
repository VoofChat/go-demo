package main

import (
	"context"
	"demo4/dao"
	"demo4/dao/base"
	"demo4/po"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

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

// query1 查询示例1
func query1(ctx context.Context) {
	userDao := dao.NewUserDao(db)

	user, err := userDao.GetByUserId(ctx, 2020062500017)
	// SELECT * FROM `tbl_user` WHERE user_id = 2020062500017
	if err != nil {
		fmt.Println("GetByUserId failed. ", err)
		return
	}
	fmt.Println("GetByUserId success. ", jsonMarshal(user))

	user, err = userDao.GetByUsername(ctx, "zhangsan")
	// SELECT * FROM `tbl_user` WHERE username = 'zhangsan'
	if err != nil {
		fmt.Println("GetByUsername failed. ", err)
		return
	}
	fmt.Println("GetByUsername success. ", jsonMarshal(user))

	// 基础类提供的查询方法
	var user1 po.User
	err = userDao.GetByCond(ctx, &user1, base.Where("username = ? and password = ?", "zhangsan", 123456))
	// SELECT * FROM `tbl_user` WHERE username = 'zhangsan' and password = 123456
	if err != nil {
		fmt.Println("GetByCond failed. ", err)
		return
	}
	fmt.Println("GetByCond success. ", jsonMarshal(user1))
}

// update1 修改示例1
func update1(ctx context.Context) {
	userDao := dao.NewUserDao(db)

	rowAffects, err := userDao.UpdateUsernameById(ctx, 2020062500016, "lisi")
	// UPDATE `tbl_user` SET `username`='lisi' WHERE user_id = 2020062500016
	fmt.Println("UpdateUsernameById info. ", rowAffects, err) // 不存在 0 <nil>

	rowAffects, err = userDao.UpdateUsernameById(ctx, 2020062500017, "lisi")
	// UPDATE `tbl_user` SET `username`='lisi' WHERE user_id = 2020062500017
	fmt.Println("UpdateUsernameById info. ", rowAffects, err) // 存在 1 <nil>
}

// save1 保存示例1
func save1(ctx context.Context) {
	userDao := dao.NewUserDao(db)
	user := po.User{
		UserName: "wangwu",
		Password: "123456",
	}

	err := userDao.Save(ctx, &user) // save 参数注意是指针类型
	//  INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('wangwu','123456',0)
	if err != nil {
		fmt.Println("Save failed. ", err)
		return
	}
	fmt.Println("Save success. ", jsonMarshal(user))
	// {"user_id":2020062500077,"username":"wangwu","password":"123456","deleted":0,"create_time":"0001-01-01T00:00:00Z","update_time":"0001-01-01T00:00:00Z"}
}

// transaction1 事务示例1
func transaction1(ctx context.Context) {
	err := db.Transaction(func(tx *gorm.DB) error {
		userDao := dao.NewUserDao(tx) // 注意这里传入的是tx
		userInfo, err := userDao.GetByUserId(ctx, 2020062500017)
		// SELECT * FROM `tbl_user` WHERE user_id = 2020062500017
		if err != nil {
			fmt.Println("GetByUserId failed. ", err)
			return err
		}

		updateUser := po.User{
			UserID:   userInfo.UserID,
			UserName: "lurenjia",
		}

		rowAffects, err := userDao.Update(ctx, &updateUser)
		// UPDATE `tbl_user` SET `username`='lurenjia' WHERE `user_id` = 2020062500017
		if rowAffects == 0 || err != nil {
			fmt.Println("Update failed. ", err)
			return err
		}

		saveUser := po.User{
			UserName: "lisi",
			Password: "123123",
		}
		err = userDao.Save(ctx, &saveUser)
		// INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('lisi','123123',0)
		if err != nil {
			fmt.Println("Update failed. ", err)
			return err
		}
		return fmt.Errorf("手动触发事务失败") // 这里返回任何错误都会回滚事务
	})

	if err != nil {
		fmt.Println("transaction failed. ", err)
		return
	}
	fmt.Println("transaction success. ")
}

// transaction2 事务示例2
func transaction2(ctx context.Context) {
	// 开启事务
	tx := db.Begin()
	var err error
	defer func() {
		if err != nil {
			fmt.Println("transaction failed. ", err)
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	userDao := dao.NewUserDao(tx) // 注意这里传入的是tx
	userInfo, err := userDao.GetByUserId(ctx, 2020062500017)
	// SELECT * FROM `tbl_user` WHERE user_id = 2020062500017
	if err != nil {
		fmt.Println("GetByUserId failed. ", err)
	}

	updateUser := po.User{
		UserID:   userInfo.UserID,
		UserName: "lurenjia",
	}

	rowAffects, err := userDao.Update(ctx, &updateUser)
	// UPDATE `tbl_user` SET `username`='lurenjia' WHERE `user_id` = 2020062500017
	if rowAffects == 0 || err != nil {
		fmt.Println("Update failed. ", err)
	}

	saveUser := po.User{
		UserName: "lisi",
		Password: "123123",
	}
	err = userDao.Save(ctx, &saveUser)
	// INSERT INTO `tbl_user` (`username`,`password`,`deleted`) VALUES ('lisi','123123',0)
	if err != nil {
		fmt.Println("Update failed. ", err)
	}

	err = fmt.Errorf("手动触发事务失败") // 手动触发事务失败
}

func main() {
	ctx := context.Background()

	//query1(ctx)
	//update1(ctx)
	//save1(ctx)

	//transaction1(ctx)
	transaction2(ctx)
}
