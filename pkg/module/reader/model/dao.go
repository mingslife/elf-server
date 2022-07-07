package model

import (
	"time"
)

const (
	ReaderGenderMale   = uint8(0)
	ReaderGenderFemale = uint8(1)
)

type Reader struct {
	ID        uint       `json:"id" gorm:"primary_key;autoIncrement"`
	UniqueID  string     `json:"uniqueId" gorm:"size:50;unique"`
	Nickname  string     `json:"nickname" gorm:"not null;size:128"`
	Password  string     `json:"-" gorm:"size:50"` // no use
	Gender    uint8      `json:"gender" gorm:"not null"`
	Birthday  *time.Time `json:"birthday" gorm:"type:date"`
	Email     string     `json:"email" gorm:"not null;size:255"`
	EmailHash string     `json:"-" gorm:"size:32;unique"`
	Phone     string     `json:"phone" gorm:"size:50"`
	UserID    *uint      `json:"userId" gorm:"unique"`
	OpenID    string     `json:"openId" gorm:"size:255"`
	IsActive  bool       `json:"isActive" gorm:"not null;index"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}
