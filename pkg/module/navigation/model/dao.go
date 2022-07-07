package model

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
