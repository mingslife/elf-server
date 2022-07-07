package service

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/navigation/model"
	"elf-server/pkg/module/navigation/repo"
)

type NavigationService struct {
	Repo *repo.NavigationRepo `inject:""`
}

func (s *NavigationService) List(ctx context.Context, limit, page int) ([]*model.Navigation, int64, error) {
	return s.Repo.List(ctx, limit, page)
}

func (s *NavigationService) Get(ctx context.Context, id uint) (*model.Navigation, error) {
	return s.Repo.Get(ctx, id)
}

func (s *NavigationService) Create(ctx context.Context, v *model.Navigation) error {
	return s.Repo.Create(ctx, v)
}

func (s *NavigationService) Update(ctx context.Context, v *model.Navigation) error {
	return s.Repo.Update(ctx, v)
}

func (s *NavigationService) Delete(ctx context.Context, v *model.Navigation) error {
	return s.Repo.Delete(ctx, v)
}

func (s *NavigationService) ListAll(ctx context.Context) ([]*model.Navigation, error) {
	return s.Repo.ListAll(ctx)
}

func (s *NavigationService) ListAllActive(ctx context.Context) ([]*model.Navigation, error) {
	return s.Repo.ListAllActive(ctx)
}

var _ bone.Service = (*NavigationService)(nil)
