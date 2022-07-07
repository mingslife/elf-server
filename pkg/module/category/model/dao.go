package model

import (
	"time"
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
