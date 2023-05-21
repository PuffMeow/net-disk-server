package db

import (
	"database/sql"
	"fmt"
	mydb "net-disk-server/db/mysql"
)

// 文件上传完成
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	db := mydb.DbConnect()
	if db == nil {
		fmt.Println("Failed to connect to the database")
		return false
	}

	stmt, err := db.Prepare("INSERT IGNORE INTO table_file(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) VALUES (?, ?, ?, ?, 1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:", err)
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println("Exec err:", err)
		return false
	}

	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("File with hash %s has been uploaded before", filehash)
		}
		return true
	}

	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// 从 mysql 获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DbConnect().Prepare(
		"SELECT file_sha1, file_addr, file_name, file_size FROM table_file " +
			"WHERE file_sha1=? AND status=1 LIMIT 1",
	)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	tableFile := TableFile{}

	err = stmt.QueryRow(filehash).Scan(&tableFile.FileHash, &tableFile.FileAddr, &tableFile.FileName, &tableFile.FileSize)

	defer stmt.Close()

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &tableFile, nil

}
