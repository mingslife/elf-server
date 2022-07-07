package model

import (
	"time"
)

const (
	UserRoleAdmin    = uint8(0)
	UserRoleAuthor   = uint8(1)
	UserGenderMale   = uint8(0)
	UserGenderFemale = uint8(1)
)

type User struct {
	ID           uint       `json:"id" gorm:"primary_key;autoIncrement"`
	Username     string     `json:"username" gorm:"not null;size:64;unique" binding:"alphanum"`
	Password     string     `json:"password" gorm:"not null;size:64"`
	Nickname     string     `json:"nickname" gorm:"not null;size:64"`
	Email        string     `json:"email" gorm:"not null;size:255" binding:"email"`
	EmailHash    string     `json:"-" gorm:"not null;size:32;unique"`
	Phone        string     `json:"phone" gorm:"size:50" binding:"numeric"`
	PhoneHash    string     `json:"-" gorm:"not null;size:32;unique"`
	Tags         string     `json:"tags" gorm:"size:255"`
	Introduction string     `json:"introduction" gorm:"size:255"`
	IsActive     bool       `json:"isActive" gorm:"not null"`
	ActiveAt     *time.Time `json:"activeAt"`
	Avatar       string     `json:"avatar" gorm:"size:255"`
	Gender       uint8      `json:"gender" gorm:"not null"`
	Birthday     *time.Time `json:"birthday" gorm:"type:date"`
	Role         uint8      `json:"role" gorm:"not null"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}
