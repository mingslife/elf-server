package model

import (
	"time"

	categorymodel "elf-server/pkg/module/category/model"
	poststatisticsmodel "elf-server/pkg/module/poststatistics/model"
	usermodel "elf-server/pkg/module/user/model"
)

type Post struct {
	ID               uint                               `json:"id" gorm:"primary_key;autoIncrement"`
	UniqueID         string                             `json:"uniqueId" gorm:"size:36;unique"`
	PostStatistics   poststatisticsmodel.PostStatistics `json:"postStatistics" gorm:"foreignkey:UniqueID;references:unique_id"`
	Title            string                             `json:"title" gorm:"not null;size:128"`
	Keywords         string                             `json:"keywords" gorm:"size:255"`
	Description      string                             `json:"description" gorm:"size:255"`
	UserID           uint                               `json:"userId" gorm:"not null"`
	User             usermodel.User                     `json:"user" gorm:"foreignkey:UserID" binding:"-"`
	CategoryID       uint                               `json:"categoryId" gorm:"not null"`
	Category         categorymodel.Category             `json:"category" gorm:"foreignkey:CategoryID"`
	Cover            string                             `json:"cover" gorm:"size:255"`
	SourceType       string                             `json:"sourceType" gorm:"not null;size:32"`
	Source           string                             `json:"source" gorm:"not null;type:text"`
	Content          string                             `json:"content" gorm:"not null;type:text"`
	Route            string                             `json:"route" gorm:"not null;size:255"`
	RouteHash        string                             `json:"-" gorm:"size:32;unique"`
	IsPublished      bool                               `json:"isPublished" gorm:"not null;index"`
	PublishedAt      *time.Time                         `json:"publishedAt"`
	IsPrivate        bool                               `json:"isPrivate" gorm:"not null;index"`
	Password         string                             `json:"password"`
	IsCommentEnabled bool                               `json:"isCommentEnabled" gorm:"not null"`
	IsCommentShown   bool                               `json:"isCommentShown" gorm:"not null"`
	CreatedAt        time.Time                          `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time                          `json:"updatedAt" gorm:"autoUpdateTime"`
}
