package endpoint

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/comment/model"
	"elf-server/pkg/module/comment/service"
)

type CommentEndpoint struct {
	Service *service.CommentService `inject:""`
}

func (e *CommentEndpoint) List(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.ListReq)
	rows, total, err := e.Service.List(ctx, req.Limit, req.Page)
	return &model.ListRsp{
		Rows: rows, Total: total,
	}, err
}

func (e *CommentEndpoint) Get(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.GetReq)
	row, err := e.Service.Get(ctx, req.ID)
	return *model.GetRsp(row), err
}

func (e *CommentEndpoint) Create(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.CreateReq)
	err = e.Service.Create(ctx, req.Row)
	return &model.CreateRsp{}, err
}

func (e *CommentEndpoint) Update(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.UpdateReq)
	err = e.Service.Update(ctx, req.Row)
	return &model.UpdateRsp{}, err
}

func (e *CommentEndpoint) Delete(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.DeleteReq)
	err = e.Service.Delete(ctx, req.Row)
	return &model.DeleteRsp{}, err
}

var _ bone.Endpoint = (*CommentEndpoint)(nil)
