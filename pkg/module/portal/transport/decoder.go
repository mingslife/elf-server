package transport

import (
	"context"
	"net/http"

	"elf-server/pkg/module/portal/model"
	"elf-server/pkg/utils"
)

type PortalDecoder struct{}

func (*PortalDecoder) Index(_ context.Context, r *http.Request) (req any, err error) {
	return nil, nil
}

func (*PortalDecoder) Post(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.PostReq{}, error(nil)
	req.Route, err = utils.GetParam(r, "route").String()
	return req, err
}

func (*PortalDecoder) User(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.UserReq{}, error(nil)
	req.Username, err = utils.GetParam(r, "username").String()
	if err != nil {
		return nil, err
	}
	if req.Page, err = utils.GetParam(r, "page").Int(); err != nil {
		req.Page, err = 1, nil
	}
	return req, err
}

func (*PortalDecoder) Category(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.CategoryReq{}, error(nil)
	req.Route, err = utils.GetParam(r, "route").String()
	if err != nil {
		return nil, err
	}
	if req.Page, err = utils.GetParam(r, "page").Int(); err != nil {
		req.Page, err = 1, nil
	}
	return req, err
}

func (*PortalDecoder) Posts(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.PostsReq{}, error(nil)
	if req.Page, err = utils.GetParam(r, "page").Int(); err != nil {
		req.Page, err = 1, nil
	}
	return req, err
}

func (*PortalDecoder) Reader(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}
