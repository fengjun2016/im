package main

import (
	"imWebSocket/app"
	"imWebSocket/handler"
	"imWebSocket/migrate"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func RegisterView(router *httprouter.Router) {
	//通配符解决模板渲染的文件
	tpl, err := template.ParseGlob("web/view/**/*.html")
	//如果报错则停止
	if err != nil {
		logrus.Fatal(err.Error())
	}

	//遍历循环所有的模板文件进行相应的模板文件执行
	for _, v := range tpl.Templates() {
		tplname := v.Name()
		logrus.Println("模板名称:", tplname)
		if strings.HasPrefix(tplname, "/") {
			router.GET(tplname, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
				logrus.Println("is execute", tplname)
				tpl.ExecuteTemplate(w, tplname, nil)
			})
		}
	}

}

func main() {
	app.InitDb()
	//自动建表
	migrate.CreateTable()

	router := httprouter.New()
	router.POST("/user/login", handler.Login)
	router.POST("/user/register", handler.Register)

	//contact
	router.POST("/user/addfriend", handler.AddFriend)

	//提供静态资源目录支持
	router.ServeFiles("/asset/*filepath", http.Dir("web/asset"))

	//解析所有的模板文件
	RegisterView(router)

	log.Printf("http.ListenAndServer ip:port: %s", "127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
