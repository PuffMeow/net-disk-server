package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net-disk-server/meta"
	"net-disk-server/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 接收文件流及存储到本地目录
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Faild to get data %s", err.Error())
		return
	}

	defer file.Close()

	// 如果没有 tmp 文件夹则会去创建
	err = os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		fmt.Println("无法创建文件夹:", err)
		return
	}

	// 存放的路径
	filePath := filepath.Join("tmp", header.Filename)

	// 记录文件元信息
	fileMeta := meta.FileMeta{
		FileName: header.Filename,
		Location: filePath,
		UploadAt: time.Now().Format(time.DateTime),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Faild to create file %s", err.Error())
		return
	}
	// 延迟关闭，释放资源
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)

	if err != nil {
		fmt.Printf("Faild to save file %s", err.Error())
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)

	_ = meta.UpdateFileMetaDB(fileMeta)

	data, err := json.Marshal(fileMeta)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// 根据 filehash 入参查询文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	if filehash == "" {
		return
	}

	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fileMeta)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 返回文件元数据信息
	w.Write(data)
}

// 批量查询元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// Atoi 字符串转数字
	limitCnt, err := strconv.Atoi(r.Form.Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileMetas := meta.GetLastFileMetas(limitCnt)
	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 根据 filehash 下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filesha1 := r.Form.Get("filehash")
	filemeta := meta.GetFileMeta(filesha1)
	f, err := os.Open(filemeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// 将文件一次性存到内存并返回，如果文件大的话需要用流的形式返回
	data, err := io.ReadAll(f)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 添加响应头让浏览器可以识别到
	w.Header().Set("Content-Type:", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\""+filemeta.FileName+"\"")
	w.Write(data)
}

// 更新元信息（重命名）
func UpdateFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	// 0 表示删除
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 简单的删除元信息，多线程时需要加锁
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fileMeta.Location)

	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Delete Success")
}
