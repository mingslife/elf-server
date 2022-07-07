package model

import (
	"elf-server/pkg/component"
	categorymodel "elf-server/pkg/module/category/model"
	navigationmodel "elf-server/pkg/module/navigation/model"
	postmodel "elf-server/pkg/module/post/model"
	usermodel "elf-server/pkg/module/user/model"
)

type Common struct {
	Template    string
	Settings    map[string]string
	Navigations []*navigationmodel.Navigation
	Title       string
	Keywords    string
	Description string
}

func (c *Common) TemplatePath() string {
	return c.Template
}

var _ component.PageData = (*Common)(nil)

type IndexRsp struct {
	*Common
	Posts []*postmodel.Post
	Limit int
	Page  int
	Total int64
	Pages int
}

type PostReq struct {
	Route string
}

type PostRsp struct {
	*Common
	Title       string
	Keywords    string
	Description string
	Post        *postmodel.Post
}

type UserReq struct {
	Username string
	Page     int
}

type UserRsp struct {
	*Common
	User  *usermodel.User
	Posts []*postmodel.Post
	Limit int
	Page  int
	Total int64
	Pages int
}

type CategoryReq struct {
	Route string
	Page  int
}

type CategoryRsp struct {
	*Common
	Category *categorymodel.Category
	Posts    []*postmodel.Post
	Limit    int
	Page     int
	Total    int64
	Pages    int
}

type PostsReq struct {
	Page int
}

type PostsRsp struct {
	*Common
	Posts []*postmodel.Post
	Limit int
	Page  int
	Total int64
	Pages int
}

type ReaderRsp struct {
	*Common
}
