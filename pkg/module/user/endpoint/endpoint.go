package endpoint

import (
	"context"
	"errors"

	"github.com/mingslife/bone"

	"elf-server/pkg/models"
	readermodel "elf-server/pkg/module/reader/model"
	readerservice "elf-server/pkg/module/reader/service"
	"elf-server/pkg/module/user/model"
	"elf-server/pkg/module/user/service"
)

type UserEndpoint struct {
	Service       *service.UserService         `inject:""`
	ReaderService *readerservice.ReaderService `inject:""`
}

func (e *UserEndpoint) List(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.ListReq)
	rows, total, err := e.Service.List(ctx, req.Limit, req.Page)
	return &model.ListRsp{
		Rows: rows, Total: total,
	}, err
}

func (e *UserEndpoint) Get(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.GetReq)
	row, err := e.Service.Get(ctx, req.ID)
	return *model.GetRsp(row), err
}

func (e *UserEndpoint) Create(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.CreateReq)
	err = e.Service.Create(ctx, req.Row)
	if err != nil {
		return nil, err
	}
	err = e.ReaderService.Create(ctx, &readermodel.Reader{
		Nickname: req.Row.Nickname,
		Gender:   req.Row.Gender,
		Birthday: req.Row.Birthday,
		Email:    req.Row.Email,
		Phone:    req.Row.Phone,
		UserID:   &req.Row.ID,
		IsActive: req.Row.IsActive,
	})
	return &model.CreateRsp{}, err
}

func (e *UserEndpoint) Update(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.UpdateReq)
	err = e.Service.Update(ctx, req.Row)
	if err != nil {
		return nil, err
	}
	var reader *readermodel.Reader
	reader, err = e.ReaderService.GetByUserID(ctx, req.Row.ID)
	if err != nil {
		return nil, err
	}
	reader.Nickname = req.Row.Nickname
	reader.Gender = req.Row.Gender
	reader.Birthday = req.Row.Birthday
	reader.Email = req.Row.Email
	reader.Phone = req.Row.Phone
	reader.IsActive = req.Row.IsActive
	e.ReaderService.Update(ctx, reader)
	return &model.UpdateRsp{}, err
}

func (e *UserEndpoint) Delete(ctx context.Context, r any) (rsp any, err error) {
	req := r.(*model.DeleteReq)
	var user *model.User
	user, err = e.Service.Get(ctx, req.Row.ID)
	if err != nil {
		return nil, err
	}
	if user.Role == models.UserRoleAdmin {
		return nil, errors.New("cannot delete admin user")
	}
	err = e.Service.Delete(ctx, req.Row)
	return &model.DeleteRsp{}, err
}

var _ bone.Endpoint = (*UserEndpoint)(nil)
