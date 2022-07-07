package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"elf-server/pkg/module/category/model"
	"elf-server/pkg/utils"
)

type CategoryDecoder struct{}

func (*CategoryDecoder) List(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.ListReq{}, error(nil)
	req.Page, _ = utils.GetQueryDefault(r, "page", "1").Int()
	req.Limit, _ = utils.GetQueryDefault(r, "limit", "10").Int()
	return req, err
}

func (*CategoryDecoder) Get(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.GetReq{}, error(nil)
	req.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*CategoryDecoder) Create(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.CreateReq{
		Row: &model.Category{},
	}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req.Row)
	return req, err
}

func (*CategoryDecoder) Update(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.UpdateReq{
		Row: &model.Category{},
	}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req.Row)
	if err != nil {
		return req, err
	}
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*CategoryDecoder) Delete(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.DeleteReq{
		Row: &model.Category{},
	}, error(nil)
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*CategoryDecoder) ListAll(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}
