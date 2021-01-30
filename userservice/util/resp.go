package util

import (
	"encoding/json"
	"log"
)

// RespMsg : http响应数据的通用结构
type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewRespMsg : 生成response对象
func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func JsonBytes(msg *RespMsg) []byte {
	r, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return r
}

func JsonString(msg *RespMsg) string {
	return string(JsonBytes(msg))
}
