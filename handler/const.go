package handler

type LoginRequest struct {
}

type RegisterRequest struct {
	Mobile   string `json:"mobile" binding:"required"`
	Passwd   string `json:"passwd" binding:"required"`
	NickName string `json:"nick_name" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Sex      string `json:"sexs" binding:"required"`
}
