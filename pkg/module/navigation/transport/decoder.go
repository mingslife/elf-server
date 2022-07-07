package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"elf-server/pkg/module/navigation/model"
	"elf-server/pkg/utils"
)

type NavigationDecoder struct{}

func (*NavigationDecoder) List(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.ListReq{}, error(nil)
	req.Page, _ = utils.GetQueryDefault(r, "page", "1").Int()
	req.Limit, _ = utils.GetQueryDefault(r, "limit", "10").Int()
	return req, err
}

func (*NavigationDecoder) Get(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.GetReq{}, error(nil)
	req.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*NavigationDecoder) Create(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.CreateReq{
		Row: &model.Navigation{},
	}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req.Row)
	return req, err
}

func (*NavigationDecoder) Update(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.UpdateReq{
		Row: &model.Navigation{},
	}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req.Row)
	if err != nil {
		return req, err
	}
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*NavigationDecoder) Delete(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.DeleteReq{
		Row: &model.Navigation{},
	}, error(nil)
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*NavigationDecoder) ListAll(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}
