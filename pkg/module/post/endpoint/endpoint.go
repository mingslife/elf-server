package endpoint

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/post/model"
	"elf-server/pkg/module/post/service"
	poststatisticsmodel "elf-server/pkg/module/poststatistics/model"
	poststatisticsservice "elf-server/pkg/module/poststatistics/service"
)

type PostEndpoint struct {
	Service               *service.PostService                         `inject:""`
	PostStatisticsService *poststatisticsservice.PostStatisticsService `inject:""`
}

func (e *PostEndpoint) List(ctx context.Context, r any) (any, error) {
	req := r.(*model.ListReq)
	// rows, total, err := e.Service.List(ctx, req.Limit, req.Page)
	rows, total, err := e.Service.ListByUserID(ctx, req.UserID, req.Limit, req.Page)
	return &model.ListRsp{
		Rows: rows, Total: total,
	}, err
}

func (e *PostEndpoint) Get(ctx context.Context, r any) (any, error) {
	req := r.(*model.GetReq)
	row, err := e.Service.Get(ctx, req.ID)
	return *model.GetRsp(row), err
}

func (e *PostEndpoint) Create(ctx context.Context, r any) (any, error) {
	req, err := r.(*model.CreateReq), error(nil)
	postStatistics := &poststatisticsmodel.PostStatistics{}
	err = e.PostStatisticsService.Create(ctx, postStatistics)
	if err == nil {
		req.Row.UniqueID = postStatistics.UniqueID
		err = e.Service.Create(ctx, req.Row)
	}
	return &model.CreateRsp{}, err
}

func (e *PostEndpoint) Update(ctx context.Context, r any) (any, error) {
	req, err := r.(*model.UpdateReq), error(nil)
	err = e.Service.Update(ctx, req.Row)
	return &model.UpdateRsp{}, err
}

func (e *PostEndpoint) Delete(ctx context.Context, r any) (any, error) {
	req, err := r.(*model.DeleteReq), error(nil)
	err = e.Service.Delete(ctx, req.Row)
	return &model.DeleteRsp{}, err
}

func (e *PostEndpoint) GetContent(ctx context.Context, r any) (any, error) {
	req, err := r.(*model.GetContentReq), error(nil)
	var post *model.Post
	post, err = e.Service.GetContent(ctx, req.ID, int(req.UserID))
	if err != nil {
		return nil, err
	}
	return &model.GetContentRsp{
		ID:         post.ID,
		Source:     post.Source,
		SourceType: post.SourceType,
		Content:    post.Content,
	}, nil
}

func (e *PostEndpoint) UpdateContent(ctx context.Context, r any) (any, error) {
	req, err := r.(*model.UpdateContentReq), error(nil)
	err = e.Service.UpdateContent(ctx, req.ID, int(req.UserID), req.Source)
	if err != nil {
		return nil, err
	}
	var post *model.Post
	post, err = e.Service.GetContent(ctx, req.ID, int(req.UserID))
	return &model.UpdateContentRsp{
		ID:         post.ID,
		Source:     post.Source,
		SourceType: post.SourceType,
		Content:    post.Content,
	}, err
}

var _ bone.Endpoint = (*PostEndpoint)(nil)
