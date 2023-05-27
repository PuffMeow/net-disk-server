package handler

import (
	"fmt"
	"math"
	rPool "net-disk-server/cache/redis"
	"net-disk-server/util"
	"net/http"
	"strconv"
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
	redisConn.Do("HSET", "MP_"+uploadInfo.UploadID, "chunkcount", uploadInfo.ChunkCount)
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
