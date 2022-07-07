package repo

import (
	"context"
	"errors"

	"github.com/mingslife/bone"

	"elf-server/pkg/component"
	"elf-server/pkg/module/poststatistics/model"
	"elf-server/pkg/utils"
)

type PostStatisticsRepo struct {
	Database *component.Database `inject:"component.database"`
}

var (
	postStatisticsesFields = []string{"id", "unique_id", "page_view", "thumb_up", "thumb_down", "created_at", "updated_at"}
	postStatisticsFields   = []string{"id", "unique_id", "page_view", "thumb_up", "thumb_down", "created_at", "updated_at"}
)

func (r *PostStatisticsRepo) List(ctx context.Context, limit, page int) (s []*model.PostStatistics, cnt int64, err error) {
	err = r.Database.DB.WithContext(ctx).Model(&model.PostStatistics{}).
		Select(postStatisticsesFields).
		Count(&cnt).
		Order("id ASC").Limit(limit).Offset(utils.Offset(limit, page)).
		Find(&s).Error
	return
}

func (r *PostStatisticsRepo) Get(ctx context.Context, id uint) (*model.PostStatistics, error) {
	var v model.PostStatistics
	if err := r.Database.DB.WithContext(ctx).Select(postStatisticsFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PostStatisticsRepo) Create(ctx context.Context, v *model.PostStatistics) error {
	v.UniqueID = utils.NewUUID()
	return r.Database.DB.WithContext(ctx).Save(v).Error
}

func (r *PostStatisticsRepo) Update(ctx context.Context, v *model.PostStatistics) error {
	return r.Database.DB.WithContext(ctx).Model(v).Updates(map[string]any{
		"page_view":  v.PageView,
		"thumb_up":   v.ThumbUp,
		"thumb_down": v.ThumbDown,
	}).Error
}

func (r *PostStatisticsRepo) Delete(ctx context.Context, v *model.PostStatistics) error {
	return r.Database.DB.WithContext(ctx).Delete(v).Error
}

func (r *PostStatisticsRepo) GetByUniqueID(ctx context.Context, uniqueID string) (*model.PostStatistics, error) {
	var v model.PostStatistics
	if err := r.Database.DB.WithContext(ctx).Select(postStatisticsFields).Take(&v, "unique_id = ?", uniqueID).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PostStatisticsRepo) UpdatePageView(ctx context.Context, uniqueID string) error {
	if uniqueID == "" {
		return errors.New("uniqueID cannot be empty")
	}
	return r.Database.DB.WithContext(ctx).Exec("update post_statistics set page_view = page_view + 1 where unique_id = ?", uniqueID).Error
}

var _ bone.Repo = (*PostStatisticsRepo)(nil)
