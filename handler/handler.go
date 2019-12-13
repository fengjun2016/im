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
	LoginMust         = &Errno{Code: 11004, Msg: "请先登录账户。"}

	// register
	RegisterFailed          = &Errno{Code: 12001, Msg: "Register failed. Please confirm and retry again."}
	RegisterParamsWrong     = &Errno{Code: 12002, Msg: "Register params wrong. Please confirm and retry again."}
	RegisterRepeated        = &Errno{Code: 12003, Msg: "Register failed for repeat mobile. Please confirm and retry again."}
	RegisterBindParamsError = &Errno{Code: 12004, Msg: "Register bind params failed."}

	// contact
	AddFriendsFailed             = &Errno{Code: 13001, Msg: "添加好友失败。"}
	AddFriendsRepeatedInSameCate = &Errno{Code: 13002, Msg: "已经添加过该好友，请勿重复添加。"}
	AddFriendsCanNotSelf         = &Errno{Code: 13003, Msg: "无法添加自己作为好友。"}
	GetCommunityUserIdsFailed    = &Errno{Code: 13004, Msg: "获取用户群友信息失败。"}
	SearchMyFriendFailed         = &Errno{Code: 13005, Msg: "查找我的好友失败。"}

	// contact groups
	CreateContactGroupsFailed    = &Errno{Code: 14001, Msg: "创建群聊失败。"}
	JoinContactGroupsFailed      = &Errno{Code: 14002, Msg: "加入群聊失败。"}
	CreateContactGroupsNoSetName = &Errno{Code: 14003, Msg: "创建群聊没有设置名称。"}
	CreateContatcGroupsTooMuch   = &Errno{Code: 14004, Msg: "一个用户最多只能创见5个群。"}
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
