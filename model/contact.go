package model

import (
	"imWebSocket/app"
	"imWebSocket/model"

	"github.com/sirupsen/logrus"
)

//好友和群都存在这个表里面
//可根据具体业务做拆分
type Contact struct {
	Model
	OwnerId   string `form:"owner_id" json:"owner_id"`
	DstUserId string `form:"dst_user_id" json:"dst_user_id"`
	Cate      int    `form:"cate" json:"cate"`
	Memo      string `form:"memo" json:"memo" gorm:"type:text"`
}

const (
	SingleFriends = 1 //用户加用户的好友关系
	GroupFriends  = 2 //加群的好友关系 也即是群友关系
)

func (c *Contact) Create() error {
	//开启事务
	tx := app.DB.Begin()

	//插入自己的好友记录
	if err := tx.Create(c).Error; err != nil {
		tx.Rollback()
		return err
	}

	//新建对方的好友添加记录
	otherContact := Contact{
		OwnerId:   c.DstUserId,
		DstUserId: c.OwnerId,
		Cate:      c.Cate,
		Memo:      c.Memo,
	}

	if err := tx.Create(&otherContact).Error; err != nil {
		tx.Rollback()
		return err
	}

	//提交
	tx.Commit()
	return nil
}

func (c *Contact) CheckIsRepeatedAddFriends() error {
	return app.DB.Model(c).Where("owner_id = ?", c.OwnerId).Where("dst_user_id = ?", c.DstUserId).Where("cate = ?", c.Cate).Find(c).Error
}

func (c *Contact) SearchCommunityIds(userId string) ([]Contact, []string, int, error) {
	//todo 获取用户全部群id
	contacts := make([]Contact, 0)
	commIds := make([]string, 0)
	count := 0

	if err := app.DB.Model(c).Where("owner_id = ? and cate = ?", userId, GroupFriends).Count(&count).Find(&contacts).Error; err != nil {
		logrus.Println("get community contacts failed.", err.Error())
		return contacts, commIds, count, err
	}

	for _, c := range contacts {
		commIds = append(commIds, c.DstUserId)
	}

	return contacts, commIds, count, nil
}

//查找群
func (c *Contact) SearchCommunity() ([]model.User, int, error) {
	contacts := make([]model.Contact, 0)
	objIds := make([]string, 0)
	count := 0
	coms := make([]User, 0)
	if err := app.DB.Model(Contact{}).Where("owner_id = ? and cate = ?", c.OwnerId, GroupFriends).Find(&contacts).Error; err != nil {
		logrus.Println("get my owner contact groups failed.", err.Error())
		return coms, count, err
	}

	for _, v := range contacts {
		objIds = append(objIds, v.DstUserId)
	}

	if len(objIds) == 0 {
		return coms
	}

	if err := app.DB.Model(User{}).Where("id in (?)", objIds).Count(&count).Find(&coms).Error; err != nil {
		logrus.Println("get my contact friends user list failed.", err.Error())
		return
	}

	return coms, count, nil
}

//查找好友
func (c *Contact) SearchFriend() ([]model.User, int, error) {
	contacts := make([]model.Contact, 0)
	friends := make([]model.User, 0)
	objIds := make([]string, 0)
	count := 0
	if err := app.DB.Model(c).Where("owner_id = ? and cate = ?", c.UserId, SingleFriends).Find(&contacts).Error; err != nil {
		logrus.Println("get my friend contatcs list failed error.", err.Error())
		return friends, count, err
	}

	for _, v := range contacts {
		objIds = append(objIds, v.DstUserId)
	}

	if len(objIds) == 0 {
		logrus.Println("no friends exists")
		return firends, count, nil
	}

	if err := app.DB.Model(model.User).Where("id in (?)", objIds).Find(&friends).Count(&count).Error; err != nil {
		logrus.Println("select user form db failed", err.Error())
		return friends, count, err
	}
	return friedns, count, nil
}
