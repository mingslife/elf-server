package repo

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/component"
	"elf-server/pkg/module/navigation/model"
	"elf-server/pkg/utils"
)

type NavigationRepo struct {
	Database *component.Database `inject:"component.database"`
}

var (
	navigationsFields = []string{"id", "label", "url", "target", "position", "is_active", "created_at", "updated_at"}
	navigationFields  = []string{"id", "label", "url", "target", "position", "is_active", "created_at", "updated_at"}
)

func (r *NavigationRepo) List(ctx context.Context, limit, page int) (s []*model.Navigation, cnt int64, err error) {
	err = r.Database.DB.WithContext(ctx).Model(&model.Navigation{}).
		Select(navigationsFields).
		Count(&cnt).
		Order("position ASC").Limit(limit).Offset(utils.Offset(limit, page)).
		Find(&s).Error
	return
}

func (r *NavigationRepo) Get(ctx context.Context, id uint) (*model.Navigation, error) {
	var v model.Navigation
	if err := r.Database.DB.WithContext(ctx).Select(navigationFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *NavigationRepo) Create(ctx context.Context, v *model.Navigation) error {
	return r.Database.DB.WithContext(ctx).Save(v).Error
}

func (r *NavigationRepo) Update(ctx context.Context, v *model.Navigation) error {
	return r.Database.DB.WithContext(ctx).Model(v).Updates(map[string]any{
		"label":     v.Label,
		"url":       v.URL,
		"target":    v.Target,
		"position":  v.Position,
		"is_active": v.IsActive,
	}).Error
}

func (r *NavigationRepo) Delete(ctx context.Context, v *model.Navigation) error {
	return r.Database.DB.WithContext(ctx).Delete(v).Error
}

func (r *NavigationRepo) ListAll(ctx context.Context) (s []*model.Navigation, err error) {
	err = r.Database.DB.WithContext(ctx).Select([]string{"id", "label", "url", "target", "position"}).
		Order("position ASC").
		Find(&s).Error
	return
}

func (r *NavigationRepo) ListAllActive(ctx context.Context) (s []*model.Navigation, err error) {
	err = r.Database.DB.WithContext(ctx).Select([]string{"id", "label", "url", "target", "position"}).
		Where("is_active = 1").
		Order("position ASC").
		Find(&s).Error
	return
}

var _ bone.Repo = (*NavigationRepo)(nil)
