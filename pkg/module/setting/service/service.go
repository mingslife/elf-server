package service

import (
	"context"

	"elf-server/pkg/module/setting/model"
	"elf-server/pkg/module/setting/repo"

	"github.com/mingslife/bone"
)

type SettingService struct {
	Repo *repo.SettingRepo `inject:""`
}

func (s *SettingService) List(ctx context.Context, limit, page int) ([]*model.Setting, int64, error) {
	return s.Repo.List(ctx, limit, page)
}

func (s *SettingService) Get(ctx context.Context, id uint) (*model.Setting, error) {
	return s.Repo.Get(ctx, id)
}

func (s *SettingService) Create(ctx context.Context, v *model.Setting) error {
	return s.Repo.Create(ctx, v)
}

func (s *SettingService) Update(ctx context.Context, v *model.Setting) error {
	return s.Repo.Update(ctx, v)
}

func (s *SettingService) Delete(ctx context.Context, v *model.Setting) error {
	return s.Repo.Delete(ctx, v)
}

func (s *SettingService) ListAll(ctx context.Context) ([]*model.Setting, error) {
	return s.Repo.ListAll(ctx)
}

func (s *SettingService) GetByKey(ctx context.Context, settingKey string) (*model.Setting, error) {
	return s.Repo.GetByKey(ctx, settingKey)
}

func (s *SettingService) ListAllMap(ctx context.Context) (m map[string]string, err error) {
	if settings, err := s.Repo.ListAll(ctx); err != nil {
		return nil, err
	} else {
		m = map[string]string{}
		for _, setting := range settings {
			m[setting.SettingKey] = setting.SettingValue
		}
		return m, nil
	}
}

func (s *SettingService) ListAllPublicMap(ctx context.Context) (m map[string]string, err error) {
	if settings, err := s.Repo.ListAll(ctx); err != nil {
		return nil, err
	} else {
		m = map[string]string{}
		for _, setting := range settings {
			if !setting.IsPublic {
				continue
			}
			m[setting.SettingKey] = setting.SettingValue
		}
		return m, nil
	}
}

var _ bone.Service = (*SettingService)(nil)
