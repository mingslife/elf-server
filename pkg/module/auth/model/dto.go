package model

import "time"

type LoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginRsp struct {
	Token string `json:"token"`
}

type RefreshReq struct {
	UserID   uint
	UserRole uint8
}

type RefreshRsp struct {
	Token string `json:"token"`
}

type RegisterReq struct {
	InviteCode string
	Username   string
	Password   string
	Nickname   string
	Email      string
	Phone      string
	Gender     uint8
	Birthday   *time.Time
}

type RegisterRsp struct{}

type GetInfoReq struct {
	UserID   uint
	UserRole uint8
}

type GetInfoRsp struct {
	UserID   uint  `json:"userId"`
	UserRole uint8 `json:"userRole"`
}
