package models

import (
	"time"
)

type Migration struct {
	ID        uint      `json:"id" gorm:"primary_key;autoIncrement"`
	Version   uint64    `json:"version" gorm:"not null;unique"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

var (
	migrationsFields = []string{"id", "version", "created_at", "updated_at"}
	migrationFields  = []string{"id", "version", "created_at", "updated_at"}
)

func GetMigrations(limit, page int) (s []*Migration, count int64) {
	DB.Model(&Migration{}).
		Select(migrationsFields).
		Count(&count).
		Order("id ASC").Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func GetMigration(id uint) *Migration {
	var v Migration
	if err := DB.Select(migrationFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil
	}
	return &v
}

func (v *Migration) Save() error {
	return DB.Save(v).Error
}

func (v *Migration) Update() error {
	return DB.Model(v).Updates(map[string]interface{}{
		"version": v.Version,
	}).Error
}

func (v *Migration) Delete() error {
	return DB.Delete(v).Error
}

func GetMigrationOfMaxVersion() *Migration {
	var v Migration
	if err := DB.Select(migrationFields).Order("version desc").Take(&v).Error; err != nil {
		return nil
	}
	return &v
}
