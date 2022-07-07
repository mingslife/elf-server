package model

type ListReq struct {
	Page   int
	Limit  int
	UserID uint
}

type ListRsp struct {
	Rows  []*Post `json:"rows"`
	Total int64   `json:"total"`
}

type GetReq struct {
	ID uint
}

type GetRsp *Post

type CreateReq struct {
	Row    *Post
	UserID uint
}

type CreateRsp struct{}

type UpdateReq struct {
	Row    *Post
	UserID uint
}

type UpdateRsp struct{}

type DeleteReq struct {
	Row    *Post
	UserID uint
}

type DeleteRsp struct{}

type GetContentReq struct {
	ID     uint
	UserID uint
}

type GetContentRsp struct {
	ID         uint   `json:"id"`
	Source     string `json:"source"`
	SourceType string `json:"sourceType"`
	Content    string `json:"content"`
}

type UpdateContentReq struct {
	ID         uint
	Source     string
	SourceType string
	UserID     uint
}

type UpdateContentRsp struct {
	ID         uint   `json:"id"`
	Source     string `json:"source"`
	SourceType string `json:"sourceType"`
	Content    string `json:"content"`
}
