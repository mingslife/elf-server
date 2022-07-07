package service

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/reader/model"
	"elf-server/pkg/module/reader/repo"
)

type ReaderService struct {
	Repo *repo.ReaderRepo `inject:""`
}

func (s *ReaderService) List(ctx context.Context, limit, page int) ([]*model.Reader, int64, error) {
	return s.Repo.List(ctx, limit, page)
}

func (s *ReaderService) Get(ctx context.Context, id uint) (*model.Reader, error) {
	return s.Repo.Get(ctx, id)
}

func (s *ReaderService) Create(ctx context.Context, v *model.Reader) error {
	return s.Repo.Create(ctx, v)
}

func (s *ReaderService) Update(ctx context.Context, v *model.Reader) error {
	return s.Repo.Update(ctx, v)
}

func (s *ReaderService) Delete(ctx context.Context, v *model.Reader) error {
	return s.Repo.Delete(ctx, v)
}

func (s *ReaderService) GetByEmail(ctx context.Context, email string) (*model.Reader, error) {
	return s.Repo.GetByEmail(ctx, email)
}

func (s *ReaderService) GetByNickname(ctx context.Context, nickname string) (*model.Reader, error) {
	return s.Repo.GetByNickname(ctx, nickname)
}

func (s *ReaderService) GetByUniqueID(ctx context.Context, uniqueID string) (*model.Reader, error) {
	return s.Repo.GetByUniqueID(ctx, uniqueID)
}

func (s *ReaderService) GetByUserID(ctx context.Context, userID uint) (*model.Reader, error) {
	return s.Repo.GetByUserID(ctx, userID)
}

var _ bone.Service = (*ReaderService)(nil)
