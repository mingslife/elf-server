package model

import (
	"time"
)

type Comment struct {
	ID          uint      `json:"id" gorm:"primary_key;autoIncrement"`
	ReaderID    uint      `json:"readerId"`
	PostID      uint      `json:"postId"`
	Level       uint      `json:"level" gorm:"not null;index"`
	ParentID    *uint     `json:"parentId" gorm:"index"`
	IP          string    `json:"ip" gorm:"size:40"`
	UserAgent   string    `json:"userAgent" gorm:"size:255"`
	Content     string    `json:"content" gorm:"type:text"`
	CommentedAt time.Time `json:"commentedAt" gorm:"not null"`
	IsBlocked   bool      `json:"isBlocked" gorm:"not null"`
	IsPrivate   bool      `json:"isPrivate" gorm:"not null"`
	IsAnonymous bool      `json:"isAnonymous" gorm:"not null"` // no use
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type PortalComment struct {
	Level          uint      `json:"level"`
	Nickname       string    `json:"nickname"`
	IsBlocked      bool      `json:"isBlocked"`
	Username       *string   `json:"username"`
	Content        string    `json:"content"`
	CommentedAt    time.Time `json:"commentedAt"`
	ParentLevel    *uint     `json:"parentLevel"`
	ParentNickname *string   `json:"parentNickname"`
	ParentUsername *string   `json:"parentUsername"`
}
