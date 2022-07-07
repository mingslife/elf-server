package model

type ListReq struct {
	Page  int
	Limit int
}

type ListRsp struct {
	Rows  []*User `json:"rows"`
	Total int64   `json:"total"`
}

type GetReq struct {
	ID uint
}

type GetRsp *User

type CreateReq struct {
	Row *User
}

type CreateRsp struct{}

type UpdateReq struct {
	Row *User
}

type UpdateRsp struct{}

type DeleteReq struct {
	Row *User
}

type DeleteRsp struct{}
