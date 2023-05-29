package main

import (
	"fmt"
	"net-disk-server/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.UpdateFileMetaHandler)
	http.HandleFunc("/file/delete", handler.DeleteFileHandler)
	http.HandleFunc("/file/fastupload", handler.TryFastUploadHandler)

	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	// 初始化分块信息
	http.HandleFunc("/file/mpupload/init", handler.InitMultipartUploadHandler)
	// 上传分块
	http.HandleFunc("/file/mpupload/uploadpart", handler.UploadPartHandler)
	// 通知分块上传完成
	http.HandleFunc("/file/mpupload/complete", handler.CompleteUploadHandler)
	// 取消分块上传

	// 查看分块上传整体状态

	fmt.Printf("Start server in localhost:8080\n")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server %s", err.Error())
	}
}
