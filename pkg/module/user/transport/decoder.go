package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"elf-server/pkg/module/user/model"
	"elf-server/pkg/utils"
)

type UserDecoder struct{}

func (*UserDecoder) List(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.ListReq{}, error(nil)
	req.Page, _ = utils.GetQueryDefault(r, "page", "1").Int()
	req.Limit, _ = utils.GetQueryDefault(r, "limit", "10").Int()
	return req, err
}

func (*UserDecoder) Get(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.GetReq{}, error(nil)
	req.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*UserDecoder) Create(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.CreateReq{
		Row: &model.User{},
	}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req.Row)
	req.Row.Role = model.UserRoleAuthor // only allowed to create author user
	return req, err
}

func (*UserDecoder) Update(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.UpdateReq{
		Row: &model.User{},
	}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req.Row)
	if err != nil {
		return nil, err
	}
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*UserDecoder) Delete(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.DeleteReq{
		Row: &model.User{},
	}, error(nil)
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}
