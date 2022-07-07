package model

type ListReq struct {
	Month string
}

type ListRsp struct {
	Dates []string `json:"dates"`
}

type GetReq struct {
	Date string
}

type GetRsp struct {
	Data []byte
}
