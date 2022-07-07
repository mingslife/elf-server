package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gogf/gf/util/gconv"

	"elf-server/pkg/component"
	"elf-server/pkg/module/auth/model"
)

type AuthDecoder struct{}

func (*AuthDecoder) Login(_ context.Context, r *http.Request) (any, error) {
	req := &model.LoginReq{}
	err := json.NewDecoder(r.Body).Decode(req)
	return req, err
}

func (*AuthDecoder) Logout(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func (*AuthDecoder) Refresh(_ context.Context, r *http.Request) (any, error) {
	req := &model.RefreshReq{}
	mapClaims := component.ExtractClaims(r)
	if userID, ok := mapClaims["sub"]; ok {
		req.UserID = gconv.Uint(userID)
	} else {
		return nil, errors.New("incorrect token")
	}
	if userRole, ok := mapClaims["rol"]; ok {
		req.UserRole = gconv.Uint8(userRole)
	} else {
		return nil, errors.New("incorrect token")
	}
	return req, nil
}

func (*AuthDecoder) GetSettings(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func (*AuthDecoder) GetInfo(_ context.Context, r *http.Request) (any, error) {
	req := &model.GetInfoReq{}
	mapClaims := component.ExtractClaims(r)
	if userID, ok := mapClaims["sub"]; ok {
		req.UserID = gconv.Uint(userID)
	} else {
		return nil, errors.New("incorrect token")
	}
	if userRole, ok := mapClaims["rol"]; ok {
		req.UserRole = gconv.Uint8(userRole)
	} else {
		return nil, errors.New("incorrect token")
	}
	return req, nil
}
