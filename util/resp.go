package util

import (
	"encoding/json"
	"log"
)

type RespMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func NewRespMsg(code int, msg string, data any) *RespMsg {
	return &RespMsg{Code: code, Msg: msg, Data: data}
}

// 对象转 json 格式的二进制数组
func (resp *RespMsg) JSONBytes() []byte {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return r
}

// 对象转 json string
func (resp *RespMsg) JSONString() string {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}

	return string(r)
}
