package model

type ListReq struct {
	Page  int
	Limit int
}

type ListRsp struct {
	Rows  []*Category `json:"rows"`
	Total int64       `json:"total"`
}

type GetReq struct {
	ID uint
}

type GetRsp *Category

type CreateReq struct {
	Row *Category
}

type CreateRsp struct{}

type UpdateReq struct {
	Row *Category
}

type UpdateRsp struct{}

type DeleteReq struct {
	Row *Category
}

type DeleteRsp struct{}
