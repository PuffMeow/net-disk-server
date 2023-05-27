package db

import (
	"database/sql"
	"fmt"
	mydb "net-disk-server/db/mysql"
	"net-disk-server/util"
	"time"
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

func UserSignIn(username string, encodePassword string) bool {
	stmt, err := mydb.DbConnect().Prepare("SELECT user_pwd FROM table_user WHERE user_name=? LIMIT 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	var password string
	err = stmt.QueryRow(username).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			// 处理未找到用户的情况
			fmt.Println("User not found")
		} else {
			// 处理其他查询错误
			fmt.Println("Query error:", err)
		}
		return false
	}

	if encodePassword == password {
		return true
	}

	return false
}

func GenerateToken(username string) string {
	// 抽成 40 位的 token 字符
	// md5(username + timestamp + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPreffix := util.MD5([]byte(username + ts + "_token_salt"))
	return tokenPreffix + ts[:8]
}

// 刷新用户登录 token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DbConnect().Prepare(
		"REPLACE INTO table_token (user_name, user_token) VALUES (?, ?)",
	)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DbConnect().Prepare("SELECT user_name, signup_at FROM table_user WHERE user_name = ? LIMIT 1")

	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)

	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	return user, nil
}
