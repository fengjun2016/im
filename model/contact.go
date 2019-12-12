package model

import "imWebSocket/app"

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
