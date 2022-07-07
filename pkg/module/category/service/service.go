package service

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/category/model"
	"elf-server/pkg/module/category/repo"
)

type CategoryService struct {
	Repo *repo.CategoryRepo `inject:""`
}

func (s *CategoryService) List(ctx context.Context, limit, page int) ([]*model.Category, int64, error) {
	return s.Repo.List(ctx, limit, page)
}

func (s *CategoryService) Get(ctx context.Context, id uint) (*model.Category, error) {
	return s.Repo.Get(ctx, id)
}

func (s *CategoryService) Create(ctx context.Context, v *model.Category) error {
	return s.Repo.Create(ctx, v)
}

func (s *CategoryService) Update(ctx context.Context, v *model.Category) error {
	return s.Repo.Update(ctx, v)
}

func (s *CategoryService) Delete(ctx context.Context, v *model.Category) error {
	return s.Repo.Delete(ctx, v)
}

func (s *CategoryService) ListAll(ctx context.Context) ([]*model.Category, error) {
	return s.Repo.ListAll(ctx)
}

func (s *CategoryService) GetByRoute(ctx context.Context, route string) (*model.Category, error) {
	return s.Repo.GetByRoute(ctx, route)
}

var _ bone.Service = (*CategoryService)(nil)
