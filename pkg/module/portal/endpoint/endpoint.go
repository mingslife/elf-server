package endpoint

import (
	"context"
	"math"

	"github.com/gogf/gf/util/gconv"
	"github.com/mingslife/bone"

	categoryservice "elf-server/pkg/module/category/service"
	navigationservice "elf-server/pkg/module/navigation/service"
	"elf-server/pkg/module/portal/model"
	postservice "elf-server/pkg/module/post/service"
	poststatisticsservice "elf-server/pkg/module/poststatistics/service"
	settingservice "elf-server/pkg/module/setting/service"
	userservice "elf-server/pkg/module/user/service"
)

type PortalEndpoint struct {
	UserService           *userservice.UserService                     `inject:""`
	CategoryService       *categoryservice.CategoryService             `inject:""`
	PostService           *postservice.PostService                     `inject:""`
	PostStatisticsService *poststatisticsservice.PostStatisticsService `inject:""`
	SettingService        *settingservice.SettingService               `inject:""`
	NavigationService     *navigationservice.NavigationService         `inject:""`
}

func (e *PortalEndpoint) makeCommon(ctx context.Context) *model.Common {
	settings, _ := e.SettingService.ListAllMap(ctx)
	navigations, _ := e.NavigationService.Repo.ListAllActive(ctx)
	return &model.Common{
		Settings:    settings,
		Navigations: navigations,
	}
}

func (e *PortalEndpoint) Index(ctx context.Context, req any) (any, error) {
	common := e.makeCommon(ctx)
	rsp, err := &model.IndexRsp{Common: common}, error(nil)

	limit := gconv.Int(common.Settings["app.limit"])
	total, _ := e.PostService.CountForPortal(ctx, "", "")
	pages := int(math.Ceil(float64(total) / float64(limit)))
	posts, _ := e.PostService.ListForPortal(ctx, "", "", limit, 1)

	rsp.Template = "index.jet"
	rsp.Posts = posts
	rsp.Limit = limit
	rsp.Page = 1
	rsp.Total = total
	rsp.Pages = pages

	return rsp, err
}

func (e *PortalEndpoint) Post(ctx context.Context, req any) (any, error) {
	r := req.(*model.PostReq)

	common := e.makeCommon(ctx)
	rsp, err := &model.PostRsp{Common: common}, error(nil)

	post, err := e.PostService.GetForPortal(ctx, r.Route)
	if post == nil || err != nil {
		rsp.Template = "error.jet"
		return rsp, nil
	}
	if post.IsPrivate || post.Category.IsPrivate {
		post.Content = ""
	}
	go e.PostStatisticsService.UpdatePageView(context.Background(), post.UniqueID)

	rsp.Template = "post.jet"
	rsp.Title = post.Title
	rsp.Keywords = post.Keywords
	rsp.Description = post.Description
	rsp.Post = post

	return rsp, err
}

func (e *PortalEndpoint) User(ctx context.Context, req any) (any, error) {
	r := req.(*model.UserReq)

	common := e.makeCommon(ctx)
	rsp, err := &model.UserRsp{Common: common}, error(nil)

	user, _ := e.UserService.GetByUsername(ctx, r.Username)
	limit := gconv.Int(common.Settings["app.limit"])
	total, _ := e.PostService.CountForPortal(ctx, r.Username, "")
	pages := int(math.Ceil(float64(total) / float64(limit)))
	posts, _ := e.PostService.ListForPortal(ctx, r.Username, "", limit, r.Page)

	rsp.Template = "user.jet"
	rsp.Title = user.Nickname
	rsp.Keywords = user.Tags
	rsp.Description = user.Introduction
	rsp.User = user
	rsp.Posts = posts
	rsp.Limit = limit
	rsp.Page = r.Page
	rsp.Total = total
	rsp.Pages = pages

	return rsp, err
}

func (e *PortalEndpoint) Category(ctx context.Context, req any) (any, error) {
	r := req.(*model.CategoryReq)

	common := e.makeCommon(ctx)
	rsp, err := &model.CategoryRsp{Common: common}, error(nil)

	category, _ := e.CategoryService.GetByRoute(ctx, r.Route)
	limit := gconv.Int(common.Settings["app.limit"])
	total, _ := e.PostService.CountForPortal(ctx, "", r.Route)
	pages := int(math.Ceil(float64(total) / float64(limit)))
	posts, _ := e.PostService.ListForPortal(ctx, "", r.Route, limit, r.Page)

	rsp.Template = "category.jet"
	rsp.Title = category.CategoryName
	rsp.Keywords = category.Keywords
	rsp.Description = category.Description
	rsp.Category = category
	rsp.Posts = posts
	rsp.Limit = limit
	rsp.Page = r.Page
	rsp.Total = total
	rsp.Pages = pages

	return rsp, err
}

func (e *PortalEndpoint) Posts(ctx context.Context, req any) (any, error) {
	r := req.(*model.PostsReq)

	common := e.makeCommon(ctx)
	rsp, err := &model.PostsRsp{Common: common}, error(nil)

	limit := gconv.Int(common.Settings["app.limit"])
	total, _ := e.PostService.CountForPortal(ctx, "", "")
	pages := int(math.Ceil(float64(total) / float64(limit)))
	posts, _ := e.PostService.ListForPortal(ctx, "", "", limit, r.Page)

	rsp.Template = "posts.jet"
	// rsp.Title = ""
	// rsp.Keywords = ""
	// rsp.Description = ""
	rsp.Posts = posts
	rsp.Limit = limit
	rsp.Page = r.Page
	rsp.Total = total
	rsp.Pages = pages

	return rsp, err
}

func (e *PortalEndpoint) Reader(ctx context.Context, req any) (any, error) {
	common := e.makeCommon(ctx)
	rsp, err := &model.ReaderRsp{Common: common}, error(nil)

	rsp.Template = "reader.jet"

	return rsp, err
}

var _ bone.Endpoint = (*PortalEndpoint)(nil)
