package model

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
