package handler

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
