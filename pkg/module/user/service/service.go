package service

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/user/model"
	"elf-server/pkg/module/user/repo"
)

type UserService struct {
	Repo *repo.UserRepo `inject:""`
}

func (s *UserService) List(ctx context.Context, limit, page int) ([]*model.User, int64, error) {
	return s.Repo.List(ctx, limit, page)
}

func (s *UserService) Get(ctx context.Context, id uint) (*model.User, error) {
	return s.Repo.Get(ctx, id)
}

func (s *UserService) Create(ctx context.Context, v *model.User) error {
	return s.Repo.Create(ctx, v)
}

func (s *UserService) Update(ctx context.Context, v *model.User) error {
	return s.Repo.Update(ctx, v)
}

func (s *UserService) Delete(ctx context.Context, v *model.User) error {
	return s.Repo.Delete(ctx, v)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.Repo.GetByUsername(ctx, username)
}

func (s *UserService) GetByAccountAndPassword(ctx context.Context, account, password string) (*model.User, error) {
	return s.Repo.GetByAccountAndPassword(ctx, account, password)
}

func (s *UserService) Exists(ctx context.Context, username, email, phone string) bool {
	return s.Repo.Exists(ctx, username, email, phone)
}

var _ bone.Service = (*UserService)(nil)
