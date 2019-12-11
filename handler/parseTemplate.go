package handler

import (
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func ParseLoginHtml(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//解析
	tpl, err := template.ParseFiles("web/view/user/login.html")
	if err != nil {
		//打印错误并且直接退出
		logrus.Fatal(err.Error())
	}

	tpl.ExecuteTemplate(w, "/user/login.shtml", nil)
}

func ParseRegisterHtml(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//解析
	tpl, err := template.ParseFiles("web/view/user/register.html")
	if err != nil {
		//打印错误并且直接退出
		logrus.Fatal(err.Error())
	}

	tpl.ExecuteTemplate(w, "/user/register.shtml", nil)
}
