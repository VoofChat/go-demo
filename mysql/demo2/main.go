package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	UserID     int64  `db:"user_id"`     // 主键
	UserName   string `db:"username"`    // 用户名
	Password   string `db:"password"`    // 密码
	Deleted    int    `db:"deleted"`     // -1 删除
	CreateTime string `db:"create_time"` // 创建时间
	UpdateTime string `db:"update_time"` // 更新时间
}

var Db *sqlx.DB

func jsonMarshal(obj interface{}) interface{} {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return obj
	}
	return string(jsonStr)
}

func init() {
	database, err := sqlx.Open("mysql", "root:qwerasdf@tcp(127.0.0.1:3306)/demo")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database
}

func insert(username, password string) (insertId int64) {
	r, err := Db.Exec("INSERT INTO `tbl_user` (`username`, `password`) VALUES(?, ?)", username, password)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}
	fmt.Println("insert succ:", id)
	return id
}

func delete(userId int64) {
	res, err := Db.Exec("delete from `tbl_user` where user_id = ?", userId)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}

	row, err := res.RowsAffected()
	if err != nil {
		fmt.Println("rows failed, ", err)
	}

	fmt.Println("delete succ: ", row)
}

func query(userId int64) {
	var person []User
	err := Db.Select(&person, "select `user_id`, `username`, `password`, `deleted`, `create_time`, `update_time` from tbl_user where user_id=?", userId)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}
	fmt.Println("select succ:", jsonMarshal(person))
}

func update(username string, userId int64) {
	res, err := Db.Exec("update tbl_user set username=? where user_id=?", username, userId)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}
	row, err := res.RowsAffected()
	if err != nil {
		fmt.Println("rows failed, ", err)
	}
	fmt.Println("update succ:", row)
}

func main() {
	userId := insert("zhangsan", "123456")
	if userId > 0 {
		query(userId)
		update("lisi", userId)
		query(userId)
		delete(userId)
	}
}
