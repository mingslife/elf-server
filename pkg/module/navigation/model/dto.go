package model

type ListReq struct {
	Page  int
	Limit int
}

type ListRsp struct {
	Rows  []*Navigation `json:"rows"`
	Total int64         `json:"total"`
}

type GetReq struct {
	ID uint
}

type GetRsp *Navigation

type CreateReq struct {
	Row *Navigation
}

type CreateRsp struct{}

type UpdateReq struct {
	Row *Navigation
}

type UpdateRsp struct{}

type DeleteReq struct {
	Row *Navigation
}

type DeleteRsp struct{}
