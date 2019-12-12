package model

type ContactGroups struct {
	Model
	Name    string `form:"name" json:"name"`       //群名称
	Ownerid int64  `form:"ownerid" json:"ownerid"` //群主ID
	Icon    string `form:"icon" json:"icon"`       //群logo
	Cate    int    `form:"cate" json:"cate"`       //群类型 默认是common 普通类型的聊天群
	Memo    string `form:"memo" json:"memo"`       //群描述或者群备注
}

const (
	CommonGroups = 1
)
