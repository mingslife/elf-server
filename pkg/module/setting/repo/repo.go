package repo

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/component"
	"elf-server/pkg/module/setting/model"
	"elf-server/pkg/utils"
)

type SettingRepo struct {
	Database *component.Database `inject:"component.database"`
}

var (
	settingsFields = []string{"id", "setting_key", "setting_value", "setting_tag", "is_public", "created_at", "updated_at"}
	settingFields  = []string{"id", "setting_key", "setting_value", "setting_tag", "is_public", "created_at", "updated_at"}
)

func (r *SettingRepo) List(ctx context.Context, limit, page int) (s []*model.Setting, cnt int64, err error) {
	err = r.Database.DB.WithContext(ctx).Model(&model.Setting{}).
		Select(settingsFields).
		Count(&cnt).
		Order("id ASC").Limit(limit).Offset(utils.Offset(limit, page)).
		Find(&s).Error
	return
}

func (r *SettingRepo) Get(ctx context.Context, id uint) (*model.Setting, error) {
	var v model.Setting
	if err := r.Database.DB.WithContext(ctx).Select(settingFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *SettingRepo) Create(ctx context.Context, v *model.Setting) error {
	return r.Database.DB.WithContext(ctx).Save(v).Error
}

func (r *SettingRepo) Update(ctx context.Context, v *model.Setting) error {
	return r.Database.DB.WithContext(ctx).Model(v).Updates(map[string]any{
		"setting_key":   v.SettingKey,
		"setting_value": v.SettingValue,
		"setting_tag":   v.SettingTag,
		"is_public":     v.IsPublic,
	}).Error
}

func (r *SettingRepo) Delete(ctx context.Context, v *model.Setting) error {
	return r.Database.DB.WithContext(ctx).Delete(v).Error
}

func (r *SettingRepo) ListAll(ctx context.Context) (s []*model.Setting, err error) {
	err = r.Database.DB.WithContext(ctx).Select([]string{"id", "setting_key", "setting_value", "is_public"}).
		Order("id ASC").
		Find(&s).Error
	return
}

func (r *SettingRepo) GetByKey(ctx context.Context, settingKey string) (*model.Setting, error) {
	var v model.Setting
	if err := r.Database.DB.WithContext(ctx).Select(settingFields).Take(&v, "setting_key = ?", settingKey).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

var _ bone.Repo = (*SettingRepo)(nil)
