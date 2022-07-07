package model

type ListReq struct {
	Page  int
	Limit int
}

type ListRsp struct {
	Rows  []*Setting `json:"rows"`
	Total int64      `json:"total"`
}

type GetReq struct {
	ID uint
}

type GetRsp *Setting

type CreateReq struct {
	Row *Setting
}

type CreateRsp struct{}

type UpdateReq struct {
	Row *Setting
}

type UpdateRsp struct{}

type DeleteReq struct {
	Row *Setting
}

type DeleteRsp struct{}
