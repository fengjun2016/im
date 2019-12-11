package app

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

var DB *gorm.DB

func InitDb() {
	var err error
	DB, err = gorm.Open("mysql", "root:@/chat?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logrus.Panic(err)
	}

	//mysql连接池设置
	DB.SingularTable(true)                       //全局禁用表名复数
	DB.DB().SetMaxOpenConns(300)                 //最大连接数
	DB.DB().SetMaxIdleConns(100)                 //最大空闲连接数
	DB.DB().SetConnMaxLifetime(30 * time.Second) //每个连接的过期时间

	logrus.Println("db database connected successfully.")
	DB.LogMode(true)
}
