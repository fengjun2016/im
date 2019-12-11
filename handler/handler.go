package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	// common errors
	OK                  = &Errno{Code: 0, Msg: "OK"}
	InternalServerError = &Errno{Code: 10001, Msg: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Msg: "Error occurred while binding the request body to the struct."}

	// login
	LoginFailed       = &Errno{Code: 11001, Msg: "Login failed. Please confirm and retry again."}
	LoginParamsFailed = &Errno{Code: 11002, Msg: "Login params wrong. Please confirm and retry again."}
	LoginNoUserExist  = &Errno{Code: 11003, Msg: "该手机号没有进行注册。"}

	// register
	RegisterFailed          = &Errno{Code: 12001, Msg: "Register failed. Please confirm and retry again."}
	RegisterParamsWrong     = &Errno{Code: 12002, Msg: "Register params wrong. Please confirm and retry again."}
	RegisterRepeated        = &Errno{Code: 12003, Msg: "Register failed for repeat mobile. Please confirm and retry again."}
	RegisterBindParamsError = &Errno{Code: 12004, Msg: "Register bind params failed."}
)

type Errno struct {
	Code int
	Msg  string
	Err  error
}

func (e *Errno) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", e.Code, e.Msg, e.Err)
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Msg
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Msg
	}

	return InternalServerError.Code, err.Error()
}

type Response struct {
	Code int         `json:"code"` // code为0时，表示成功
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"` // code != 0 时的错误信息
}

func ResponseJson(w http.ResponseWriter, data interface{}, err error) {
	code, msg := DecodeErr(err)
	re := Response{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	//json marshal
	jsonRes, _ := json.Marshal(re)

	//设置header
	w.Header().Set("Content-Type", "application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)
	//返回json数据格式
	w.Write([]byte(jsonRes))
}
