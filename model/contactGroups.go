package model

import "imWebSocket/app"

type ContactGroups struct {
	Model
	Name    string `form:"name" json:"name"`        //群名称
	OwnerId int64  `form:"ownerid" json:"owner_id"` //群主ID
	Icon    string `form:"icon" json:"icon"`        //群logo
	Cate    int    `form:"cate" json:"cate"`        //群类型 默认是common 普通类型的聊天群
	Memo    string `form:"memo" json:"memo"`        //群描述或者群备注
}

const (
	CommonGroups = 1
)

func (c *ContactGroups) Create() error {
	//创建群聊 首先自己得是这个群聊的成员 所以这里需要一个事务操作才行
	tx := app.DB.Begin()
	if err := tx.Model(c).Create(c).Error; err != nil {
		tx.Rollback()
		return err
	}

	contact := Contact{
		OwnerId:   c.OwnerId,
		DstUserId: c.ID,
		Cate:      GroupFriends,
	}

	if err := tx.Model(contact).Create(&contact).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
