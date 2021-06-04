package models

import (
	"time"
)

type Setting struct {
	ID           uint      `json:"id" gorm:"primary_key;autoIncrement"`
	SettingKey   string    `json:"settingKey" gorm:"not null;size:64;unique"`
	SettingValue string    `json:"settingValue" gorm:"type:text"`
	SettingTag   string    `json:"settingTag" gorm:"size:255"`
	IsPublic     bool      `json:"isPublic" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

var (
	settingsFields = []string{"id", "setting_key", "setting_value", "setting_tag", "is_public", "created_at", "updated_at"}
	settingFields  = []string{"id", "setting_key", "setting_value", "setting_tag", "is_public", "created_at", "updated_at"}
)

func GetSettings(limit, page int) (s []*Setting, count int64) {
	DB.Model(&Setting{}).
		Select(settingsFields).
		Count(&count).
		Order("id ASC").Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func GetSetting(id uint) *Setting {
	var v Setting
	if err := DB.Select(settingFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil
	}
	return &v
}

func (v *Setting) Save() error {
	return DB.Save(v).Error
}

func (v *Setting) Update() error {
	return DB.Model(v).Updates(map[string]interface{}{
		"setting_key":   v.SettingKey,
		"setting_value": v.SettingValue,
		"setting_tag":   v.SettingTag,
		"is_public":     v.IsPublic,
	}).Error
}

func (v *Setting) Delete() error {
	return DB.Delete(v).Error
}

func GetAllSettings() (s []*Setting) {
	DB.Select([]string{"id", "setting_key", "setting_value", "is_public"}).
		Order("id ASC").
		Find(&s)
	return
}

func GetAllSettingsMap() (m map[string]string) {
	settings := GetAllSettings()
	m = map[string]string{}
	for _, setting := range settings {
		m[setting.SettingKey] = setting.SettingValue
	}
	return
}

func GetSettingByKey(settingKey string) *Setting {
	var v Setting
	if err := DB.Select(settingFields).Take(&v, "setting_key = ?", settingKey).Error; err != nil {
		return nil
	}
	return &v
}

func GetAllPublicSettingsMap() (m map[string]string) {
	settings := GetAllSettings()
	m = map[string]string{}
	for _, setting := range settings {
		if !setting.IsPublic {
			continue
		}
		m[setting.SettingKey] = setting.SettingValue
	}
	return
}
