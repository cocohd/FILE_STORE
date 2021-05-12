package util

import (
	"encoding/json"
	"log"
)

type RespMsg struct {
	Code int		`json:"code"`
	Msg string		`json:"msg"`
	Data interface{}`json:"data"'`
}

func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg: msg,
		Data: data,
	}
}

func (resp *RespMsg) JSONBytes() []byte {
	res, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return res
}
