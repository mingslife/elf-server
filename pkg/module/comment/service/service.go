package service

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/module/comment/model"
	"elf-server/pkg/module/comment/repo"
)

type CommentService struct {
	Repo *repo.CommentRepo `inject:""`
}

func (s *CommentService) List(ctx context.Context, limit, page int) ([]*model.Comment, int64, error) {
	return s.Repo.List(ctx, limit, page)
}

func (s *CommentService) Get(ctx context.Context, id uint) (*model.Comment, error) {
	return s.Repo.Get(ctx, id)
}

func (s *CommentService) Create(ctx context.Context, v *model.Comment) error {
	return s.Repo.Create(ctx, v)
}

func (s *CommentService) Update(ctx context.Context, v *model.Comment) error {
	return s.Repo.Update(ctx, v)
}

func (s *CommentService) Delete(ctx context.Context, v *model.Comment) error {
	return s.Repo.Delete(ctx, v)
}

func (s *CommentService) ListByPostID(ctx context.Context, postID int, limit, page int) ([]*model.Comment, int64, error) {
	return s.Repo.ListByPostID(ctx, postID, limit, page)
}

func (s *CommentService) CountByPostID(ctx context.Context, postID int) (int64, error) {
	return s.Repo.CountByPostID(ctx, postID)
}

func (s *CommentService) GetByPostIDAndLevel(ctx context.Context, postID, level uint) (*model.Comment, error) {
	return s.Repo.GetByPostIDAndLevel(ctx, postID, level)
}

func (s *CommentService) ListByPostIDForPortal(ctx context.Context, postID uint) ([]*model.PortalComment, int64, error) {
	return s.Repo.ListByPostIDForPortal(ctx, postID)
}

func (s *CommentService) SetIsBlocked(ctx context.Context, id uint, isBlocked bool) error {
	return s.Repo.SetIsBlocked(ctx, id, isBlocked)
}

var _ bone.Service = (*CommentService)(nil)
