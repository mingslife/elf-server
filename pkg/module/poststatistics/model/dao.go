package model

import (
	"time"
)

type PostStatistics struct {
	ID        uint      `json:"id" gorm:"primary_key;autoIncrement"`
	UniqueID  string    `json:"uniqueId" gorm:"not null;size:36;unique"`
	PageView  int64     `json:"pageView" gorm:"not null"`
	ThumbUp   int64     `json:"thumbUp" gorm:"not null"`
	ThumbDown int64     `json:"thumbDown" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
