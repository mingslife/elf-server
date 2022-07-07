package endpoint

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/navigation/model"
	"elf-server/pkg/module/navigation/service"
)

type NavigationEndpoint struct {
	Service *service.NavigationService `inject:""`
}

func (e *NavigationEndpoint) List(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.ListReq)
	rows, total, err := e.Service.List(ctx, req.Limit, req.Page)
	return &model.ListRsp{
		Rows: rows, Total: total,
	}, err
}

func (e *NavigationEndpoint) Get(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.GetReq)
	row, err := e.Service.Get(ctx, req.ID)
	return *model.GetRsp(row), err
}

func (e *NavigationEndpoint) Create(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.CreateReq)
	err = e.Service.Create(ctx, req.Row)
	return &model.CreateRsp{}, err
}

func (e *NavigationEndpoint) Update(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.UpdateReq)
	err = e.Service.Update(ctx, req.Row)
	return &model.UpdateRsp{}, err
}

func (e *NavigationEndpoint) Delete(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.DeleteReq)
	err = e.Service.Delete(ctx, req.Row)
	return &model.DeleteRsp{}, err
}

func (e *NavigationEndpoint) ListAll(ctx context.Context, r any) (rsp any, err error) {
	return e.Service.ListAll(ctx)
}

var _ bone.Endpoint = (*NavigationEndpoint)(nil)
