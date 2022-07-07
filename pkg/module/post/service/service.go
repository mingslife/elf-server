package service

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/post/model"
	"elf-server/pkg/module/post/repo"
)

type PostService struct {
	Repo *repo.PostRepo `inject:""`
}

func (s *PostService) List(ctx context.Context, limit, page int) ([]*model.Post, int64, error) {
	return s.Repo.List(ctx, limit, page)
}

func (s *PostService) Get(ctx context.Context, id uint) (*model.Post, error) {
	return s.Repo.Get(ctx, id)
}

func (s *PostService) Create(ctx context.Context, v *model.Post) error {
	return s.Repo.Create(ctx, v)
}

func (s *PostService) Update(ctx context.Context, v *model.Post) error {
	return s.Repo.Update(ctx, v)
}

func (s *PostService) Delete(ctx context.Context, v *model.Post) error {
	return s.Repo.Delete(ctx, v)
}

func (s *PostService) ListByUserID(ctx context.Context, userID uint, limit, page int) ([]*model.Post, int64, error) {
	return s.Repo.ListByUserID(ctx, userID, limit, page)
}

func (s *PostService) GetContent(ctx context.Context, id uint, userID int) (*model.Post, error) {
	return s.Repo.GetContent(ctx, id, userID)
}

func (s *PostService) UpdateContent(ctx context.Context, id uint, userID int, source string) error {
	return s.Repo.UpdateContent(ctx, id, userID, source)
}

func (s *PostService) ListForPortal(ctx context.Context, username, categoryRoute string, limit, page int) ([]*model.Post, error) {
	return s.Repo.ListForPortal(ctx, username, categoryRoute, limit, page)
}

func (s *PostService) CountForPortal(ctx context.Context, username, categoryRoute string) (int64, error) {
	return s.Repo.CountForPortal(ctx, username, categoryRoute)
}

func (s *PostService) GetForPortal(ctx context.Context, route string) (*model.Post, error) {
	return s.Repo.GetForPortal(ctx, route)
}

func (s *PostService) GetForPortalByUniqueID(ctx context.Context, uniqueID string) (*model.Post, error) {
	return s.Repo.GetForPortalByUniqueID(ctx, uniqueID)
}

func (s *PostService) GetByRoute(ctx context.Context, route string) (*model.Post, error) {
	return s.Repo.GetByRoute(ctx, route)
}

var _ bone.Service = (*PostService)(nil)
