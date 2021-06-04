package models

import (
	"time"

	"elf-server/pkg/utils"
)

type Category struct {
	ID           uint      `json:"id" gorm:"primary_key;autoIncrement"`
	CategoryName string    `json:"categoryName" gorm:"not null;size:128"`
	Keywords     string    `json:"keywords" gorm:"size:255"`
	Description  string    `json:"description" gorm:"size:255"`
	Cover        string    `json:"cover" gorm:"size:255"`
	Route        string    `json:"route" gorm:"not null;size:255"`
	RouteHash    string    `json:"-" gorm:"size:32;unique"`
	IsPrivate    bool      `json:"isPrivate" gorm:"not null;index"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

var (
	categoriesFields = []string{"id", "category_name", "route", "is_private", "created_at", "updated_at"}
	categoryFields   = []string{"id", "category_name", "keywords", "description", "cover", "route", "is_private", "created_at", "updated_at"}
)

func GetCategories(limit, page int) (s []*Category, count int64) {
	DB.Model(&Category{}).
		Select(categoriesFields).
		Count(&count).
		Order("id ASC").Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func GetCategory(id uint) *Category {
	var v Category
	if err := DB.Select(categoryFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil
	}
	return &v
}

func (v *Category) Save() error {
	v.RouteHash = utils.Md5(v.Route)
	return DB.Save(v).Error
}

func (v *Category) Update() error {
	return DB.Model(v).Updates(map[string]interface{}{
		"category_name": v.CategoryName,
		"keywords":      v.Keywords,
		"description":   v.Description,
		"cover":         v.Cover,
		"route":         v.Route,
		"route_hash":    utils.Md5(v.Route),
		"is_private":    v.IsPrivate,
	}).Error
}

func (v *Category) Delete() error {
	return DB.Delete(v).Error
}

func GetAllCategories() (s []*Category) {
	DB.Select([]string{"id", "category_name", "is_private"}).
		Order("id ASC").
		Find(&s)
	return
}

func GetCategoryByRoute(route string) *Category {
	routeHash := utils.Md5(route)
	var v Category
	if err := DB.Select(categoryFields).Take(&v, "route_hash = ?", routeHash).Error; err != nil {
		return nil
	}
	return &v
}
