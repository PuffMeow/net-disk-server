package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net-disk-server/db"
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

	// todo: 更新用户文件表
	r.ParseForm()
	username := r.Form.Get("username")
	isSuccess := db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)

	if !isSuccess {
		w.Write([]byte("Upload failed"))
		return
	}

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

	username := r.Form.Get("username")
	// fileMetas := meta.GetLastFileMetas(limitCnt)
	userFiles, err := db.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
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

// 尝试秒传接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 解析请求参数
	username := r.Form.Get("username")
	filehash := r.Form.Get("fileahsh")
	filename := r.Form.Get("filename")
	filesize := r.Form.Get("filesize")

	// 从文件列表查询相同 hash 的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)

	// 查不到记录则返回秒传失败
	if fileMeta == nil || err != nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	convFilseSize, err := strconv.ParseInt(filesize, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 上传过则将文件信息写入用户文件表，返回成功
	isSuccess := db.OnUserFileUploadFinished(username, filehash, filename, convFilseSize)
	if !isSuccess {
		resp := util.RespMsg{
			Code: -2,
			Msg:  "秒传失败，请稍后重试",
		}
		w.Write(resp.JSONBytes())
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "秒传成功",
	}
	w.Write(resp.JSONBytes())
}
