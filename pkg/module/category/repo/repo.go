package repo

import (
	"context"

	"github.com/mingslife/bone"

	"elf-server/pkg/component"
	"elf-server/pkg/module/category/model"
	"elf-server/pkg/utils"
)

type CategoryRepo struct {
	Database *component.Database `inject:"component.database"`
}

var (
	categoriesFields = []string{"id", "category_name", "route", "is_private", "created_at", "updated_at"}
	categoryFields   = []string{"id", "category_name", "keywords", "description", "cover", "route", "is_private", "created_at", "updated_at"}
)

func (r *CategoryRepo) List(ctx context.Context, limit, page int) (s []*model.Category, cnt int64, err error) {
	err = r.Database.DB.WithContext(ctx).Model(&model.Category{}).
		Select(categoriesFields).
		Count(&cnt).
		Order("id ASC").Limit(limit).Offset(utils.Offset(limit, page)).
		Find(&s).Error
	return
}

func (r *CategoryRepo) Get(ctx context.Context, id uint) (*model.Category, error) {
	var v model.Category
	if err := r.Database.DB.WithContext(ctx).Select(categoryFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *CategoryRepo) Create(ctx context.Context, v *model.Category) error {
	v.RouteHash = utils.Md5(v.Route)
	return r.Database.DB.WithContext(ctx).Save(v).Error
}

func (r *CategoryRepo) Update(ctx context.Context, v *model.Category) error {
	return r.Database.DB.WithContext(ctx).Model(v).Updates(map[string]any{
		"category_name": v.CategoryName,
		"keywords":      v.Keywords,
		"description":   v.Description,
		"cover":         v.Cover,
		"route":         v.Route,
		"route_hash":    utils.Md5(v.Route),
		"is_private":    v.IsPrivate,
	}).Error
}

func (r *CategoryRepo) Delete(ctx context.Context, v *model.Category) error {
	return r.Database.DB.WithContext(ctx).Delete(v).Error
}

func (r *CategoryRepo) ListAll(ctx context.Context) (s []*model.Category, err error) {
	err = r.Database.DB.WithContext(ctx).Select([]string{"id", "category_name", "is_private"}).
		Order("id ASC").
		Find(&s).Error
	return
}

func (r *CategoryRepo) GetByRoute(ctx context.Context, route string) (*model.Category, error) {
	routeHash := utils.Md5(route)
	var v model.Category
	if err := r.Database.DB.WithContext(ctx).Select(categoryFields).Take(&v, "route_hash = ?", routeHash).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

var _ bone.Repo = (*CategoryRepo)(nil)
