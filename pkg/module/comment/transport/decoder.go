package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"elf-server/pkg/module/comment/model"
	"elf-server/pkg/utils"
)

type CommentDecoder struct{}

func (*CommentDecoder) List(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.ListReq{}, error(nil)
	req.Page, _ = utils.GetQueryDefault(r, "page", "1").Int()
	req.Limit, _ = utils.GetQueryDefault(r, "limit", "10").Int()
	return req, err
}

func (*CommentDecoder) Get(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.GetReq{}, error(nil)
	req.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*CommentDecoder) Create(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.CreateReq{
		Row: &model.Comment{},
	}, error(nil)
	// mapClaims := component.ExtractClaims(r)
	// if userID, ok := mapClaims["sub"]; ok {
	// 	req.UserID = gconv.Uint(userID)
	// } else {
	// 	return nil, errors.New("incorrect token")
	// }
	err = json.NewDecoder(r.Body).Decode(req.Row)
	// TODO req.Row.UserID
	return req, err
}

func (*CommentDecoder) Update(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.UpdateReq{
		Row: &model.Comment{},
	}, error(nil)
	err = json.NewDecoder(r.Body).Decode(req.Row)
	if err != nil {
		return req, err
	}
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}

func (*CommentDecoder) Delete(_ context.Context, r *http.Request) (any, error) {
	req, err := &model.DeleteReq{
		Row: &model.Comment{},
	}, error(nil)
	req.Row.ID, err = utils.GetParam(r, "id").Uint()
	return req, err
}
