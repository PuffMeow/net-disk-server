package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/net-disk")
	if err != nil {
		fmt.Println("Failed to connect to MySQL:", err)
		return
	}

	db.SetMaxOpenConns(1000)
	err = db.Ping()

	if err != nil {
		fmt.Println("Failed to connect to database, err:" + err.Error())
		os.Exit(1)
	}

	fmt.Println("数据库连接成功")
}

// 返回数据库连接对象
func DbConnect() *sql.DB {
	return db
}
