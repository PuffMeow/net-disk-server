package meta

import (
	"fmt"
	"net-disk-server/db"
	"sort"
	"time"
)

// 文件元信息结构体
type FileMeta struct {
	// 作为文件的唯一标识
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

// ByUploadTime 实现了 sort.Interface 接口，用于按照上传时间排序
type ByUploadTime []FileMeta

func (bt ByUploadTime) Len() int {
	return len(bt)
}

func (bt ByUploadTime) Less(i, j int) bool {
	// 使用 time.Parse 解析字符串形式的时间
	uploadTimeI, _ := time.Parse(time.DateTime, bt[i].UploadAt)
	uploadTimeJ, _ := time.Parse(time.DateTime, bt[j].UploadAt)

	// 按照上传时间进行比较
	return uploadTimeI.Before(uploadTimeJ)
}

func (bt ByUploadTime) Swap(i, j int) {
	bt[i], bt[j] = bt[j], bt[i]
}

var fileMetas map[string]FileMeta

// 初始化
func init() {
	fileMetas = make(map[string]FileMeta)
}

// 新增或更新文件元信息
func UpdateFileMeta(fileMeta FileMeta) {
	fileMetas[fileMeta.FileSha1] = fileMeta
}

// 新增/更新文件信息到 mysql 中
func UpdateFileMetaDB(fileMeta FileMeta) bool {
	return db.OnFileUploadFinished(fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize, fileMeta.Location)
}

// 根据 sha1 获取文件元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// 从 mysql 中获取文件元信息
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tableFile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		fmt.Println("GetFileMetaDB:", err.Error())
		return nil, err
	}

	fileMeta := FileMeta{
		FileSha1: tableFile.FileHash,
		FileName: tableFile.FileName.String,
		FileSize: tableFile.FileSize.Int64,
		Location: tableFile.FileAddr.String,
	}

	return &fileMeta, nil
}

// 获取批量元信息列表
func GetLastFileMetas(limitCnt int) []FileMeta {
	fileMetaArr := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fileMetaArr = append(fileMetaArr, v)
	}

	sort.Sort(ByUploadTime(fileMetaArr))
	return fileMetaArr[0:limitCnt]
}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
