package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"elf-server/pkg/module/setting/model"
	"elf-server/pkg/utils"
)

type SettingDecoder struct{}

func (*SettingDecoder) List(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.ListReq{}, error(nil)
	req.Page, _ = utils.GetQueryDefault(r, "page", "1").Int()
	req.Limit, _ = utils.GetQueryDefault(r, "limit", "10").Int()
	return req, err
}

func (*SettingDecoder) Get(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.GetReq{}, error(nil)
	req.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*SettingDecoder) Create(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.CreateReq{
		Row: &model.Setting{},
	}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req.Row)
	return req, err
}

func (*SettingDecoder) Update(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.UpdateReq{
		Row: &model.Setting{},
	}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req.Row)
	if err != nil {
		return req, err
	}
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*SettingDecoder) Delete(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.DeleteReq{
		Row: &model.Setting{},
	}, error(nil)
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*SettingDecoder) ListAll(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}
