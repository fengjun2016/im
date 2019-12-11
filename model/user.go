package model

import (
	"imWebSocket/app"
)

const (
	SEX_WOMEN = "W"
	SEX_MEN   = "M"
	//
	SEX_UNKNOW = "U"
)

type User struct {
	Model
	Mobile   string `json:"mobile" gorm:"not null;unique"`
	Passwd   string `json:"passwd"` // 什么角色
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`       // 什么角色
	NickName string `json:"nick_name"` // 什么角色
	//加盐随机字符串6
	Salt   string `json:"-"`      // 什么角色
	Online int    `json:"online"` //是否在线
	//前端鉴权因子,
	Token string `json:"token"` // 什么角色
	Memo  string `json:"memo"`  // 什么角色
}

func (u *User) Create() error {
	return app.DB.Create(u).Error
}

func (u *User) Get() error {
	return app.DB.Model(u).Find(u).Error
}

func (u *User) Update() error {
	return app.DB.Model(u).Updates(u).Error
}

func (u *User) Delete() error {
	return app.DB.Delete(u).Error
}

func (u *User) GetByMobile() error {
	return app.DB.Model(u).Where("mobile = ?", u.Mobile).Find(u).Error
}

func (u *User) UpdateToken() error {
	return app.DB.Model(u).Updates(u).Find(u).Error
}
