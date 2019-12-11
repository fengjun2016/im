package handler

import (
	"errors"
	"fmt"
	"imWebSocket/model"
	"imWebSocket/pkg/auth"
	"imWebSocket/pkg/util"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

//用户登录
func Login(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	req.ParseForm()
	mobile := req.PostForm.Get("mobile")
	passwd := req.PostForm.Get("passwd")
	logrus.Println("mobile", mobile)
	logrus.Println("passwd", passwd)

	user := model.User{}
	user.Mobile = mobile
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
	req.ParseForm()
	mobile := req.PostForm.Get("mobile")
	passwd := req.PostForm.Get("passwd")
	logrus.Println("mobile", mobile)
	logrus.Println("passwd", passwd)

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
	// user.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
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
