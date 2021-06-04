package models

import (
	"time"
)

type Navigation struct {
	ID        uint      `json:"id" gorm:"primary_key;autoIncrement"`
	Label     string    `json:"label" gorm:"not null;size:255"`
	URL       string    `json:"url" gorm:"not null;size:255"`
	Target    string    `json:"target" gorm:"size:8"`
	Position  uint      `json:"position" gorm:"index"`
	IsActive  bool      `json:"isActive" gorm:"not null;index"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

var (
	navigationsFields = []string{"id", "label", "url", "target", "position", "is_active", "created_at", "updated_at"}
	navigationFields  = []string{"id", "label", "url", "target", "position", "is_active", "created_at", "updated_at"}
)

func GetNavigations(limit, page int) (s []*Navigation, count int64) {
	DB.Model(&Navigation{}).
		Select(navigationsFields).
		Count(&count).
		Order("position ASC").Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func GetNavigation(id uint) *Navigation {
	var v Navigation
	if err := DB.Select(navigationFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil
	}
	return &v
}

func (v *Navigation) Save() error {
	return DB.Save(v).Error
}

func (v *Navigation) Update() error {
	return DB.Model(v).Updates(map[string]interface{}{
		"label":     v.Label,
		"url":       v.URL,
		"target":    v.Target,
		"position":  v.Position,
		"is_active": v.IsActive,
	}).Error
}

func (v *Navigation) Delete() error {
	return DB.Delete(v).Error
}

func GetAllNavigations() (s []*Navigation) {
	DB.Select([]string{"id", "label", "url", "target", "position"}).
		Order("position ASC").
		Find(&s)
	return
}

func GetAllNavigationsActive() (s []*Navigation) {
	DB.Select([]string{"id", "label", "url", "target", "position"}).
		Where("is_active = 1").
		Order("position ASC").
		Find(&s)
	return
}
