package endpoint

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/category/model"
	"elf-server/pkg/module/category/service"
)

type CategoryEndpoint struct {
	Service *service.CategoryService `inject:""`
}

func (e *CategoryEndpoint) List(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.ListReq)
	rows, total, err := e.Service.List(ctx, req.Limit, req.Page)
	return &model.ListRsp{
		Rows: rows, Total: total,
	}, err
}

func (e *CategoryEndpoint) Get(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.GetReq)
	row, err := e.Service.Get(ctx, req.ID)
	return *model.GetRsp(row), err
}

func (e *CategoryEndpoint) Create(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.CreateReq)
	err = e.Service.Create(ctx, req.Row)
	return &model.CreateRsp{}, err
}

func (e *CategoryEndpoint) Update(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.UpdateReq)
	err = e.Service.Update(ctx, req.Row)
	return &model.UpdateRsp{}, err
}

func (e *CategoryEndpoint) Delete(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.DeleteReq)
	err = e.Service.Delete(ctx, req.Row)
	return &model.DeleteRsp{}, err
}

func (e *CategoryEndpoint) ListAll(ctx context.Context, r any) (rsp any, err error) {
	rsp, err = e.Service.ListAll(ctx)
	return rsp, err
}

var _ bone.Endpoint = (*CategoryEndpoint)(nil)
