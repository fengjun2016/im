package handler

import (
	"fmt"
	"time"
)

type LoginRequest struct {
	Mobile string `json:"mobile" binding:"required"`
	Passwd string `json:"passwd" binding:"required"`
}

type RegisterRequest struct {
	Mobile   string `json:"mobile" binding:"required"`
	Passwd   string `json:"passwd" binding:"required"`
	NickName string `json:"nick_name" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Sex      string `json:"sexs" binding:"required"`
}

type AddFriendsRequest struct {
	OwnerId   string `json:"owner_id"`
	DstUserId string `json:"dst_user_id"`
	Cate      int    `json:"cate"`
	Memo      string `json:"memo"`
}

type ContactArg struct {
	PageArg
	Userid int64 `json:"userid" form:"userid"`
	Dstid  int64 `json:"dstid" form:"dstid"`
}

type PageArg struct {
	//从哪页开始
	Pagefrom int `json:"pagefrom" form:"pagefrom"`
	//每页大小
	Pagesize int `json:"pagesize" form:"pagesize"`
	//关键词
	Kword string `json:"kword" form:"kword"`
	//asc：“id”  id asc
	Asc  string `json:"asc" form:"asc"`
	Desc string `json:"desc" form:"desc"`
	//
	Name string `json:"name" form:"name"`
	//
	Userid int64 `json:"userid" form:"userid"`
	//dstid
	Dstid int64 `json:"dstid" form:"dstid"`
	//时间点1
	Datefrom time.Time `json:"datafrom" form:"datafrom"`
	//时间点2
	Dateto time.Time `json:"dateto" form:"dateto"`
	//
	Total int64 `json:"total" form:"total"`
}

func (p *PageArg) GetPageSize() int {
	if p.Pagesize == 0 {
		return 100
	} else {
		return p.Pagesize
	}

}
func (p *PageArg) GetPageFrom() int {
	if p.Pagefrom < 0 {
		return 0
	} else {
		return p.Pagefrom
	}
}

func (p *PageArg) GetOrderBy() string {
	if len(p.Asc) > 0 {
		return fmt.Sprintf(" %s asc", p.Asc)
	} else if len(p.Desc) > 0 {
		return fmt.Sprintf(" %s desc", p.Desc)
	} else {
		return ""
	}
}
