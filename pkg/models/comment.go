package models

import (
	"fmt"
	"time"

	"github.com/importcjj/sensitive"
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

var sensitiveFilter *sensitive.Filter

func init() {
	sensitiveFilter = sensitive.New()
	sensitiveFilter.LoadWordDict("resources/dict.txt")
}

var (
	commentsFields = []string{"id", "post_id", "level", "parent_id", "content", "commented_at", "is_blocked", "is_private", "is_anonymous", "created_at", "updated_at"}
	commentFields  = []string{"id", "post_id", "level", "parent_id", "ip", "user_agent", "content", "commented_at", "is_blocked", "is_private", "is_anonymous", "created_at", "updated_at"}
)

func GetComments(limit, page int) (s []*Comment, count int64) {
	DB.Model(&Comment{}).
		Count(&count).
		Select(commentsFields).
		Order("id ASC").Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func GetComment(id uint) *Comment {
	var v Comment
	if err := DB.Select(commentFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil
	}
	return &v
}

func (v *Comment) Save() error {
	var maxLevel uint
	v.Content = sensitiveFilter.Replace(v.Content, '*')
	DB.Model(&Comment{}).
		Select("max(level)").
		Where("post_id = ?", v.PostID).
		Scan(&maxLevel)
	v.Level = maxLevel + 1
	v.CommentedAt = time.Now()
	return DB.Save(v).Error
}

func (v *Comment) Update() error {
	return DB.Model(v).Updates(map[string]interface{}{
		"post_id":      v.PostID,
		"ip":           v.IP,
		"user_agent":   v.UserAgent,
		"content":      sensitiveFilter.Replace(v.Content, '*'),
		"commented_at": v.CommentedAt,
	}).Error
}

func (v *Comment) Delete() error {
	return DB.Delete(v).Error
}

func GetCommentsByPostID(postID int, limit, page int) (s []*Comment, count int64) {
	db := DB.Model(&Comment{}).
		Select(commentsFields)
	if postID != -1 {
		db = db.Where("post_id = ?", postID)
	}
	db = db.Count(&count).
		Order("level DESC").
		Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func CountCommentsByPostID(postID int) (count int64) {
	db := DB.Model(&Comment{})
	if postID != -1 {
		db = db.Where("post_id = ?", postID)
	}
	db = db.Count(&count)
	return
}

func GetCommentByPostIDAndLevel(postID, level uint) *Comment {
	var v Comment
	DB.Model(&Comment{}).Select(commentFields).
		Where("post_id = ? and level = ?", postID, level).
		Scan(&v)
	return &v
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

func GetCommentsByPostIDForPortal(postID uint) (s []*PortalComment, count int64) {
	DB.Raw(fmt.Sprintf(`
	select
	a.level, c.nickname, a.is_blocked, a.commented_at,
	case when a.is_blocked = 1 then '' else a.content end as content,
	b.level as parent_level, d.nickname as parent_nickname, c.user_id,
	d.user_id as parent_user_id, e.username, f.username as parent_username
	from comments a left join comments b on a.parent_id = b.id
	left join readers c on a.reader_id = c.id
	left join readers d on b.reader_id = d.id
	left join users e on c.user_id = e.id
	left join users f on d.user_id = f.id
	where a.post_id = %d
	order by level desc
	`, postID)).Find(&s)

	DB.Model(&Comment{}).
		Where("post_id = ?", postID).
		Count(&count)
	return
}

func SetCommentIsBlocked(id uint, isBlocked bool) error {
	return DB.Model(&Comment{}).Where("id = ?", id).
		Update("is_blocked", isBlocked).Error
}
