package models

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

var (
	postStatisticsesFields = []string{"id", "unique_id", "page_view", "thumb_up", "thumb_down", "created_at", "updated_at"}
	postStatisticsFields   = []string{"id", "unique_id", "page_view", "thumb_up", "thumb_down", "created_at", "updated_at"}
)

func GetPostStatisticses(limit, page int) (s []*PostStatistics, count int64) {
	DB.Model(&PostStatistics{}).
		Select(postStatisticsesFields).
		Count(&count).
		Order("id ASC").Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func GetPostStatistics(id uint) *PostStatistics {
	var v PostStatistics
	if err := DB.Select(postStatisticsFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil
	}
	return &v
}

func (v *PostStatistics) Save() error {
	return DB.Save(v).Error
}

func (v *PostStatistics) Update() error {
	return DB.Model(v).Updates(map[string]interface{}{
		"page_view":  v.PageView,
		"thumb_up":   v.ThumbUp,
		"thumb_down": v.ThumbDown,
	}).Error
}

func (v *PostStatistics) Delete() error {
	return DB.Delete(v).Error
}

func GetPostStatisticsByUniqueID(uniqueID string) *PostStatistics {
	var v PostStatistics
	if err := DB.Select(postStatisticsFields).Take(&v, "unique_id = ?", uniqueID).Error; err != nil {
		return nil
	}
	return &v
}

func UpdatePostStatisticsPageView(uniqueID string) {
	if uniqueID == "" {
		return
	}
	DB.Exec("update post_statistics set page_view = page_view + 1 where unique_id = ?", uniqueID)
}

// func IncreasePostStatisticsThumbUp(uniqueID string) {}

// func IncreasePostStatisticsThumbDown(uniqueID string) {}
