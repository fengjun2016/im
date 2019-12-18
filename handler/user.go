package handler

import (
	"errors"
	"fmt"
	"imWebSocket/model"
	"imWebSocket/pkg/auth"
	"imWebSocket/pkg/request"
	"imWebSocket/pkg/util"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

//用户登录
func Login(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//绑定参数
	loginRequest := LoginRequest{}
	if err := request.Bind(req, &loginRequest); err != nil {
		logrus.Println("login request bind params failed.", ErrBind)
		ResponseJson(rw, "", ErrBind)
		return
	}

	mobile := loginRequest.Mobile
	passwd := loginRequest.Passwd
	if mobile == "" || passwd == "" {
		logrus.Println("login request bind params failed.", LoginParamsFailed)
		ResponseJson(rw, "", LoginParamsFailed)
		return
	}

	user := model.User{}
	user.Mobile = mobile
	user.Passwd = passwd
	//检测手机号是否存在
	if user.Mobile == "" {
		ResponseJson(rw, "", LoginParamsFailed)
		return
	}
	//如果存在则返回提示已经注册
	//如果错误为空 则代表用户表里面有相关手机注册的记录
	if err := user.GetByMobile(); err != nil {
		ResponseJson(rw, "", LoginNoUserExist)
		return
	}

	//查询到了比对密码
	err := auth.Compare(user.Passwd, passwd)
	if err != nil {
		logrus.Info(err)
		err = errors.New("密码不正确")
		ResponseJson(rw, "", err)
		return
	}

	//刷新token,安全
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := util.MD5Encode(str)

	//跟新token 并返回token
	updateUserToken := model.User{}
	updateUserToken.ID = user.ID
	updateUserToken.Token = token
	if err := updateUserToken.UpdateToken(); err != nil {
		ResponseJson(rw, "", LoginFailed)
		return
	}

	ResponseJson(rw, updateUserToken, nil)
}

func Register(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	//解析参数
	registerRequest := RegisterRequest{}
	if err := request.Bind(req, &registerRequest); err != nil {
		logrus.Println("register bind rquest params failed.", err.Error())
		ResponseJson(rw, "", ErrBind)
		return
	}

	mobile := registerRequest.Mobile
	passwd := registerRequest.Passwd

	user := model.User{}
	user.Mobile = mobile
	//检测手机号是否存在
	if user.Mobile == "" {
		ResponseJson(rw, "", RegisterParamsWrong)
		return
	}
	//如果存在则返回提示已经注册
	//如果错误为空 则代表用户表里面有相关手机注册的记录
	if err := user.GetByMobile(); err == nil {
		ResponseJson(rw, "", RegisterRepeated)
		return
	}
	//否则插入新建数据
	pwdHash, err := auth.Encrypt(passwd)
	if err != nil {
		ResponseJson(rw, "", err)
		return
	}
	user.Passwd = pwdHash

	if err := user.Create(); err != nil {
		ResponseJson(rw, "", RegisterFailed)
		return
	}
	//最后返回新用户信息
	ResponseJson(rw, user, nil)
	return
}
