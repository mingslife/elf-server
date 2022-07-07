package model

type ListReq struct {
	Page  int
	Limit int
}

type ListRsp struct {
	Rows  []*Comment `json:"rows"`
	Total int64      `json:"total"`
}

type GetReq struct {
	ID uint
}

type GetRsp *Comment

type CreateReq struct {
	Row *Comment
}

type CreateRsp struct{}

type UpdateReq struct {
	Row *Comment
}

type UpdateRsp struct{}

type DeleteReq struct {
	Row *Comment
}

type DeleteRsp struct{}
