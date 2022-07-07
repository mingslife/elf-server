package model

type ListReq struct {
	Page  int
	Limit int
}

type ListRsp struct {
	Rows  []*PostStatistics `json:"rows"`
	Total int64             `json:"total"`
}

type GetReq struct {
	ID uint
}

type GetRsp *PostStatistics

type CreateReq struct {
	Row *PostStatistics
}

type CreateRsp struct{}

type UpdateReq struct {
	Row *PostStatistics
}

type UpdateRsp struct{}

type DeleteReq struct {
	Row *PostStatistics
}

type DeleteRsp struct{}
