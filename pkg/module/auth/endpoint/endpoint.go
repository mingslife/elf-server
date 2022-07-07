package endpoint

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mingslife/bone"

	"elf-server/pkg/component"
	"elf-server/pkg/module/auth/model"
	settingservice "elf-server/pkg/module/setting/service"
	usermodel "elf-server/pkg/module/user/model"
	userservice "elf-server/pkg/module/user/service"
)

type AuthEndpoint struct {
	Jwt            *component.Jwt                 `inject:"component.jwt"`
	UserService    *userservice.UserService       `inject:""`
	SettingService *settingservice.SettingService `inject:""`
}

func (e *AuthEndpoint) Login(ctx context.Context, req any) (rsp any, err error) {
	r := req.(*model.LoginReq)

	user, _ := e.UserService.GetByAccountAndPassword(ctx, r.Account, r.Password)
	if user != nil {
		token, _, _ := e.Jwt.TokenGenerator(map[string]any{
			"sub": fmt.Sprintf("%d", user.ID),
			"rol": fmt.Sprintf("%d", user.Role),
		})
		return &model.LoginRsp{Token: token}, nil
	}

	return nil, errors.New("incorrect account or password")
}

func (*AuthEndpoint) Logout(ctx context.Context, req any) (rsp any, err error) {
	return nil, nil
}

func (e *AuthEndpoint) Refresh(ctx context.Context, req any) (rsp any, err error) {
	r := req.(*model.RefreshReq)
	token, _, _ := e.Jwt.TokenGenerator(map[string]any{
		"sub": fmt.Sprintf("%d", r.UserID),
		"rol": fmt.Sprintf("%d", r.UserRole),
	})
	return &model.RefreshRsp{Token: token}, nil
}

func (e *AuthEndpoint) Register(ctx context.Context, req any) (rsp any, err error) {
	r := req.(*model.RegisterReq)

	settings, _ := e.SettingService.ListAllMap(ctx)
	if settings["app.register"] != "true" || settings["app.inviteCode"] != r.InviteCode {
		return nil, errors.New("register failed")
	}

	// check for username, email and phone
	if e.UserService.Exists(ctx, r.Username, r.Email, r.Phone) {
		return nil, errors.New("account alreay exists")
	}

	now := time.Now()
	user := &usermodel.User{
		Username:     r.Username,
		Password:     r.Password,
		Nickname:     r.Nickname,
		Email:        r.Email,
		Phone:        r.Phone,
		Tags:         "",
		Introduction: "",
		IsActive:     true,
		ActiveAt:     &now,
		Avatar:       "/assets/avatar.svg",
		Gender:       r.Gender,
		Birthday:     r.Birthday,
		Role:         usermodel.UserRoleAuthor,
	}
	return nil, e.UserService.Create(ctx, user)
}

func (*AuthEndpoint) GetProfile(ctx context.Context, req any) (rsp any, err error) {
	return nil, nil
}

func (*AuthEndpoint) UpdateProfile(ctx context.Context, req any) (rsp any, err error) {
	return nil, nil
}

func (e *AuthEndpoint) GetSettings(ctx context.Context, req any) (rsp any, err error) {
	return e.SettingService.ListAllPublicMap(ctx)
}

func (e *AuthEndpoint) GetInfo(ctx context.Context, req any) (rsp any, err error) {
	r := req.(*model.GetInfoReq)
	return &model.GetInfoRsp{
		UserID:   r.UserID,
		UserRole: r.UserRole,
	}, nil
}

var _ bone.Endpoint = (*AuthEndpoint)(nil)
