package service

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/poststatistics/model"
	"elf-server/pkg/module/poststatistics/repo"
)

type PostStatisticsService struct {
	Repo *repo.PostStatisticsRepo `inject:""`
}

func (s *PostStatisticsService) List(ctx context.Context, limit, page int) ([]*model.PostStatistics, int64, error) {
	return s.Repo.List(ctx, limit, page)
}

func (s *PostStatisticsService) Get(ctx context.Context, id uint) (*model.PostStatistics, error) {
	return s.Repo.Get(ctx, id)
}

func (s *PostStatisticsService) Create(ctx context.Context, v *model.PostStatistics) error {
	return s.Repo.Create(ctx, v)
}

func (s *PostStatisticsService) Update(ctx context.Context, v *model.PostStatistics) error {
	return s.Repo.Update(ctx, v)
}

func (s *PostStatisticsService) Delete(ctx context.Context, v *model.PostStatistics) error {
	return s.Repo.Delete(ctx, v)
}

func (s *PostStatisticsService) GetByUniqueID(ctx context.Context, uniqueID string) (*model.PostStatistics, error) {
	return s.Repo.GetByUniqueID(ctx, uniqueID)
}

func (s *PostStatisticsService) UpdatePageView(ctx context.Context, uniqueID string) error {
	return s.Repo.UpdatePageView(ctx, uniqueID)
}

var _ bone.Service = (*PostStatisticsService)(nil)
