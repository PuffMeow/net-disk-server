package db

import (
	"fmt"
	mydb "net-disk-server/db/mysql"
	"time"
)

type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

// 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := mydb.DbConnect().Prepare("INSERT INTO table_user_file (`user_name`, `file_sha1`, `file_name`, `file_size`, `upload_at`) values (?,?,?,?,?)")

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

// 批量获取用户和文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DbConnect().Prepare("SELECT file_sha1,file_name,file_size,upload_at,last_update FROM table_user_file WHERE user_name = ? limit ?")

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(username, limit)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		userFile := UserFile{}
		err := rows.Scan(&userFile.FileHash, &userFile.FileName, &userFile.FileSize, &userFile.UploadAt, &userFile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		userFiles = append(userFiles, userFile)
	}

	return userFiles, nil
}
