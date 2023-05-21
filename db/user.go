package db

import (
	"fmt"
	mydb "net-disk-server/db/mysql"
)

// 通过用户名和密码完成 user 表的注册
func UserSignup(username string, password string) bool {
	stmt, err := mydb.DbConnect().Prepare("INSERT IGNORE INTO table_user(`user_name`, `user_pwd`) VALUES (?, ?)")

	if err != nil {
		fmt.Println("Failed to insert user, err:" + err.Error())
		return false
	}

	defer stmt.Close()

	res, err := stmt.Exec(username, password)

	if err != nil {
		fmt.Println("Failed to insert user, err:" + err.Error())
		return false
	}

	if rowsAffected, err := res.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}

	return false

}
