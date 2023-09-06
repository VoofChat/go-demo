package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // init()
	"time"
)

type User struct {
	UserID     int64     // 主键
	UserName   string    // 用户名
	Password   string    // 密码
	Deleted    int       // -1 删除
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
}

// 声明全局变量
var db *sql.DB

func initMySQL() (err error) {
	dsn := "root:qwerasdf@tcp(localhost:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("connect to db failed ,err %v\n", err)
		return
	}
	// 数组需要业务具体情况来确定
	db.SetMaxOpenConns(100) //最大连接数
	db.SetMaxIdleConns(10)  //最大空闲连接数
	return
}

// 查询单条数据示例
func queryRowDemo(id int64) {
	sqlStr := "select user_id, username, password, deleted, create_time, update_time from tbl_user where user_id=?"
	var u1 User
	// 非常重要：确保 QueryRow() 调用的 Scan() 释放了连接
	// 否则持有的连接得不到释放，从而导致连接池被占用殆尽
	// 从而导致新的请求无法建立连接
	err := db.QueryRow(sqlStr, id).Scan(&u1.UserID, &u1.UserName, &u1.Password, &u1.Deleted, &u1.CreateTime, &u1.UpdateTime)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("u1:%v\n", jsonMarshal(u1))
}

// 查询多条数据
func queryMultiRowDemo() {
	sqlStr := "select user_id, username, password, deleted, create_time, update_time from tbl_user"
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	var users []User
	//循环读取结果中的数据
	for rows.Next() {
		var u1 User
		err := rows.Scan(&u1.UserID, &u1.UserName, &u1.Password, &u1.Deleted, &u1.CreateTime, &u1.UpdateTime)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		users = append(users, u1)
	}
	fmt.Printf("users:%v\n", jsonMarshal(users))
}

func jsonMarshal(obj interface{}) interface{} {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return obj
	}
	return string(jsonStr)
}

func main() {
	if err := initMySQL(); err != nil {
		fmt.Printf("init MySQL failed,err:%v\n", err)
		return
	}

	//queryRowDemo(2020062500017)
	queryMultiRowDemo()
}
