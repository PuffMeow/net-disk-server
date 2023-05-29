package handler

import (
	"fmt"
	"math"
	rPool "net-disk-server/cache/redis"
	"net-disk-server/db"
	"net-disk-server/util"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

// 初始化分块的信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// 初始化分块上传
func InitMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析用户请求
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	intFileSize, err := strconv.ParseInt(filesize, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 获取 redis 连接
	redisConn := rPool.RedisPool().Get()
	defer rPool.RedisPool().Close()

	// 生成初始化信息
	uploadInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   int(intFileSize),
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(intFileSize) / (5 * 1024 * 1024))),
	}

	// 初始化信息写入到 redis 缓存
	redisConn.Do("HSET", "MP_"+uploadInfo.UploadID, "chunkount", uploadInfo.ChunkCount)
	redisConn.Do("HSET", "MP_"+uploadInfo.UploadID, "filehash", uploadInfo.FileHash)
	redisConn.Do("HSET", "MP_"+uploadInfo.UploadID, "filesize", uploadInfo.FileSize)

	// 准备要设置的字段和值
	fieldsAndValues := map[string]interface{}{
		"chunkcount": uploadInfo.ChunkCount,
		"filehash":   uploadInfo.FileHash,
		"filesize":   uploadInfo.FileSize,
	}

	// 执行 HMSET 命令
	redisConn.Send("HMSET", redis.Args{}.Add("MP_"+uploadInfo.UploadID).AddFlat(fieldsAndValues)...)

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: uploadInfo,
	}
	// 将初始化信息同时返回到客户端
	w.Write(resp.JSONBytes())
}

// 上传分块文件
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 解析用户请求参数
	r.ParseForm()
	// username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadID")
	chundIndex := r.Form.Get("index")

	// 获得 redis 连接池中的一个连接
	redisConn := rPool.RedisPool().Get()
	defer rPool.RedisPool().Close()

	// 获得文件句柄，用于存储分块内容
	fpath := "/data/" + uploadID + "/" + chundIndex
	// 除了当前用户可写，其它用户都只读
	os.Mkdir(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	// 一个优化的点就是在客户端算好文件 hash，到服务端再算一次，比较这两次是否相同，如果相同则说明没被篡改
	// 1M 的 buffer
	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}
	// 更新 redis 缓存数据
	redisConn.Do("HSET", "MP_"+uploadID, "chunkindex_"+chundIndex, 1)
	// 返回处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uploadID := r.Form.Get("uploadID")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	redisConn := rPool.RedisPool().Get()
	defer rPool.RedisPool().Close()

	// 通过 uploadid 查询 redis 并判断是否上传完成所有分块
	data, err := redis.Values(redisConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Complete upload failed", nil).JSONBytes())
		return
	}

	totalCount := 0
	chunkCount := 0
	// 通过 HGETALL 查出来的 data ，key 和 value 都在同一个数组里，所有加2
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chunkindex_") && v == "1" {
			chunkCount += 1
		}
	}

	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "Invalid request", nil).JSONBytes())
		return
	}
	// 合并分块

	// 更新唯一文件表以及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	db.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	db.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	// 响应处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// 通知取消上传
func CancelUploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 删除存在的分块文件

	// 删除 redis 缓存状态

	// 更新 mysql 文件 status
}

// 查看分块上传状态
func MultipartUploadStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 检查上传分块是否有效

	// 获得分块初始化信息

	// 获取已上传的分块信息
}
