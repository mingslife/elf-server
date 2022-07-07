package endpoint

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/setting/model"
	"elf-server/pkg/module/setting/service"
)

type SettingEndpoint struct {
	Service *service.SettingService `inject:""`
}

func (e *SettingEndpoint) List(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.ListReq)
	rows, total, err := e.Service.List(ctx, req.Limit, req.Page)
	return &model.ListRsp{
		Rows: rows, Total: total,
	}, err
}

func (e *SettingEndpoint) Get(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.GetReq)
	row, err := e.Service.Get(ctx, req.ID)
	return *model.GetRsp(row), err
}

func (e *SettingEndpoint) Create(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.CreateReq)
	err = e.Service.Create(ctx, req.Row)
	return &model.CreateRsp{}, err
}

func (e *SettingEndpoint) Update(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.UpdateReq)
	err = e.Service.Update(ctx, req.Row)
	return &model.UpdateRsp{}, err
}

func (e *SettingEndpoint) Delete(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.DeleteReq)
	err = e.Service.Delete(ctx, req.Row)
	return &model.DeleteRsp{}, err
}

func (e *SettingEndpoint) ListAll(ctx context.Context, r any) (rsp any, err error) {
	return e.Service.ListAll(ctx)
}

var _ bone.Endpoint = (*SettingEndpoint)(nil)
