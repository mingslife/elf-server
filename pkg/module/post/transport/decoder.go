package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"elf-server/pkg/component"
	"elf-server/pkg/module/post/model"
	"elf-server/pkg/utils"

	"github.com/gogf/gf/util/gconv"
)

type PostDecoder struct{}

func (*PostDecoder) List(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.ListReq{}, error(nil)
	mapClaims := component.ExtractClaims(r)
	if userID, ok := mapClaims["sub"]; ok {
		req.UserID = gconv.Uint(userID)
	} else {
		return nil, errors.New("incorrect token")
	}
	req.Page, _ = utils.GetQueryDefault(r, "page", "1").Int()
	req.Limit, _ = utils.GetQueryDefault(r, "limit", "10").Int()
	return req, err
}

func (*PostDecoder) Get(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.GetReq{}, error(nil)
	req.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*PostDecoder) Create(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.CreateReq{
		Row: &model.Post{},
	}, error(nil)
	mapClaims := component.ExtractClaims(r)
	if userID, ok := mapClaims["sub"]; ok {
		req.UserID = gconv.Uint(userID)
	} else {
		return nil, errors.New("incorrect token")
	}
	err = json.NewDecoder(r.Body).Decode(req.Row)
	req.Row.UserID = req.UserID
	return req, err
}

func (*PostDecoder) Update(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.UpdateReq{
		Row: &model.Post{},
	}, error(nil)
	mapClaims := component.ExtractClaims(r)
	if userID, ok := mapClaims["sub"]; ok {
		req.UserID = gconv.Uint(userID)
	} else {
		return nil, errors.New("incorrect token")
	}
	err = json.NewDecoder(r.Body).Decode(req.Row)
	if err != nil {
		return req, err
	}
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	req.Row.UserID = req.UserID
	return req, err
}

func (*PostDecoder) Delete(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.DeleteReq{
		Row: &model.Post{},
	}, error(nil)
	mapClaims := component.ExtractClaims(r)
	if userID, ok := mapClaims["sub"]; ok {
		req.UserID = gconv.Uint(userID)
	} else {
		return nil, errors.New("incorrect token")
	}
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*PostDecoder) GetContent(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.GetContentReq{}, error(nil)
	mapClaims := component.ExtractClaims(r)
	if userID, ok := mapClaims["sub"]; ok {
		req.UserID = gconv.Uint(userID)
	} else {
		return nil, errors.New("incorrect token")
	}
	req.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*PostDecoder) UpdateContent(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.UpdateContentReq{}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return req, err
	}
	mapClaims := component.ExtractClaims(r)
	if userID, ok := mapClaims["sub"]; ok {
		req.UserID = gconv.Uint(userID)
	} else {
		return nil, errors.New("incorrect token")
	}
	req.ID, err = utils.GetParam(r, "id").Uint()
	return req, nil
}
