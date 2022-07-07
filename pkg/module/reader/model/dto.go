package model

type ListReq struct {
	Page  int
	Limit int
}

type ListRsp struct {
	Rows  []*Reader `json:"rows"`
	Total int64     `json:"total"`
}

type GetReq struct {
	ID uint
}

type GetRsp *Reader

type CreateReq struct {
	Row *Reader
}

type CreateRsp struct{}

type UpdateReq struct {
	Row *Reader
}

type UpdateRsp struct{}

type DeleteReq struct {
	Row *Reader
}

type DeleteRsp struct{}
