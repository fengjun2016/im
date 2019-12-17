package handler

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"
)

//程序运行初始化 创建一个本地存储图片文件夹
func init() {
	os.MkdirAll("./mnt", os.ModePerm)
}

func Upload(rw http.ResponseWriter, req *http.Request) {
	//UploadLocal(rw, req)
	UploadOss(rw, req)
}

//1.存储位置 ./mnt 需要确保已经创建好
//2.url格式 /mnt/xxxx.png 需要确保网络能够访问 /mnt/
func UploadLocal(rw http.ResponseWriter, req *http.Request) {
	//todo 获得上传的源文件
	sfile, head, err := req.FormFile("file")
	if err != nil {
		logrus.Println("get upload file failed.", err)
		ResponseJson(rw, "", err.Error())
		return
	}
	//todo 创建一个新文件
	suffix := ".png" //文件后缀
	//如果前端文件名称包含后缀 xxxx.png
	ofilename := head.Filename
	tmp := strings.Split(ofilename, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}

	//如果前端指定filetype
	fileType := req.FormValue("filetype")
	if len(fileType) > 0 {
		suffix = fileType
	}
	filename := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix) //拼凑文件名称
	dstfile, err := os.Create("./mnt/" + filename)
	if err != nil {
		logrus.Println("create new dst file failed.", err.Error())
		ResponseJson(rw, "", err.Error())
		return
	}
	//todo 将源文件内容copy到新文件
	if _, err := io.Copy(dstfile, sfile); err != nil {
		logrus.Println("copy source file to new dst file failed.", err.Error())
		ResponseJson(rw, "", err.Error())
		return
	}

	//拼接url返回给前端渲染
	url := "/mnt/" + filename
	ResponseJson(rw, url, nil)
}

//即将删掉 定期更新
const (
	AccessKeyId     = "5p2RZKnrUanMuQw9"
	AccessKeySecure = "bsNmjU8Au08axedV40TRPCS5XIFAkK"
	EndPoint        = "oss-cn-shenzhen.aliyuncs.com"
	Bucket          = "winliondev"
)

//权限设置为公共读状态
//需要安装
func UploadOss(rw http.ResponseWriter, req *http.Request) {
	//todo 获得上传的文件
	srcfile, header, err := req.FormFile("file")
	if err != nil {
		ResponseJson(rw, err.Error())
		logrus.Println("get source file from request failed.", err.Error())
		return
	}

	//todo 获得文件后缀.png / .mp3
	suffix := ".png"
	//如果前端文件名称包含后缀 xx.xx.png
	ofilename := head.Filename
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	//如果前端指定filetype
	//gormdata.append("filetype", ".png")
	filetype := req.FormValue("filetype")
	if len(filetype) > 0 {
		suffix = filetype
	}

	//todo 初始化ossclient
	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		ResponseJson(rw, err.Error())
		logrus.Println("init oss client failed error", err.Error())
		return
	}

	//todo 获得bucket
	bucket, err := client.Bucket(Bucket)
	if err != nil {
		ResponseJson(rw, err.Error())
		logrus.Println("get bucket from oss failed error", err.Error())
		return
	}

	//todo 设置文件名称
	// time.Now().Unix()
	filename := fmt.Sprintf("mnt/%d%04d%s",
		time.Now().Unix(), rand.Int31(),
		suffix)

	//todo 通过bucket上传
	err = bucket.PutObject(filename, srcfile)
	if err != nil {
		ResponseJson(rw, err.Error())
		logrus.Println("bucket upload file to online failed error", err.Error())
		return
	}
	//todo 获得url地址
	url := "http://" + Bucket + "." + EndPoint + "/" + filename

	//todo 响应到前端
	ResponseJson(rw, url, nil)
}
