package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"gorm.io/gorm"

	"elf-server/pkg/utils"
)

type Post struct {
	ID               uint           `json:"id" gorm:"primary_key;autoIncrement"`
	UniqueID         string         `json:"uniqueId" gorm:"size:36;unique"`
	PostStatistics   PostStatistics `json:"postStatistics" gorm:"foreignkey:UniqueID;references:unique_id"`
	Title            string         `json:"title" gorm:"not null;size:128"`
	Keywords         string         `json:"keywords" gorm:"size:255"`
	Description      string         `json:"description" gorm:"size:255"`
	UserID           uint           `json:"userId" gorm:"not null"`
	User             User           `json:"user" gorm:"foreignkey:UserID" binding:"-"`
	CategoryID       uint           `json:"categoryId" gorm:"not null"`
	Category         Category       `json:"category" gorm:"foreignkey:CategoryID"`
	Cover            string         `json:"cover" gorm:"size:255"`
	SourceType       string         `json:"sourceType" gorm:"not null;size:32"`
	Source           string         `json:"source" gorm:"not null;type:text"`
	Content          string         `json:"content" gorm:"not null;type:text"`
	Route            string         `json:"route" gorm:"not null;size:255"`
	RouteHash        string         `json:"-" gorm:"size:32;unique"`
	IsPublished      bool           `json:"isPublished" gorm:"not null;index"`
	PublishedAt      *time.Time     `json:"publishedAt"`
	IsPrivate        bool           `json:"isPrivate" gorm:"not null;index"`
	Password         string         `json:"password"`
	IsCommentEnabled bool           `json:"isCommentEnabled" gorm:"not null"`
	IsCommentShown   bool           `json:"isCommentShown" gorm:"not null"`
	CreatedAt        time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
}

var (
	postsFields = []string{"id", "title", "user_id", "category_id", "source_type", "route", "is_published", "published_at", "is_private", "is_comment_enabled", "is_comment_shown", "created_at", "updated_at"}
	postFields  = []string{"id", "title", "keywords", "description", "user_id", "category_id", "cover", "source_type", "route", "is_published", "is_private", "password", "is_comment_enabled", "is_comment_shown", "created_at", "updated_at"}
)

func GetPosts(limit, page int) (s []*Post) {
	DB.Select(postsFields).
		Order("id DESC").Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func CountPosts() (count int64) {
	DB.Model(&Post{}).Count(&count)
	return
}

func GetPost(id uint) *Post {
	var v Post
	if err := DB.Select(postFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil
	}
	return &v
}

func (v *Post) Save() error {
	v.RouteHash = utils.Md5(v.Route)
	v.UniqueID = utils.NewUUID()
	if v.IsPublished {
		now := time.Now()
		v.PublishedAt = &now
	}
	DB.Save(&PostStatistics{UniqueID: v.UniqueID})
	return DB.Save(v).Error
}

func (v *Post) Update() error {
	m := map[string]interface{}{
		"title":              v.Title,
		"keywords":           v.Keywords,
		"description":        v.Description,
		"user_id":            v.UserID,
		"category_id":        v.CategoryID,
		"cover":              v.Cover,
		"route":              v.Route,
		"route_hash":         utils.Md5(v.Route),
		"is_published":       v.IsPublished,
		"is_private":         v.IsPrivate,
		"is_comment_enabled": v.IsCommentEnabled,
		"is_comment_shown":   v.IsCommentShown,
	}
	if v.IsPublished {
		var o Post
		if err := DB.Select([]string{"id", "is_published"}).Take(&o, "id = ?", v.ID).Error; err != nil {
			return err
		}
		if !o.IsPublished {
			m["published_at"] = time.Now()
		}
	} else {
		m["published_at"] = nil
	}
	if v.IsPrivate {
		m["password"] = v.Password
	}
	return DB.Model(v).Updates(m).Error
}

func (v *Post) Delete() error {
	return DB.Delete(v).Error
}

func GetPostsByUserID(userID uint, limit, page int) (s []*Post, count int64) {
	DB.Model(&Post{}).Select([]string{
		"posts.id", "posts.unique_id", "page_view", "thumb_up", "thumb_down", "title", "user_id", "nickname", "category_id", "category_name", "source_type", "posts.route", "is_published", "published_at", "posts.is_private", "is_comment_enabled", "is_comment_shown", "posts.created_at", "posts.updated_at",
	}).
		Where("user_id = ?", userID).
		Count(&count).
		Joins("LEFT JOIN users on users.id = posts.user_id").
		Joins("LEFT JOIN categories on categories.id = posts.category_id").
		Joins("LEFT JOIN post_statistics on post_statistics.unique_id = posts.unique_id").
		Order("id DESC").Limit(limit).Offset(Offset(limit, page)).
		Preload("PostStatistics", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "unique_id", "page_view", "thumb_up", "thumb_down"})
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "username", "nickname", "avatar"})
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "category_name", "is_private"})
		}).
		Find(&s)
	return
}

func GetPostContent(id uint, userID int) *Post {
	var v Post
	if err := DB.Select([]string{
		"id", "source", "source_type", "content",
	}).Take(&v, "id = ? and user_id = ?", id, userID).Error; err != nil {
		return nil
	}
	return &v
}

func UpdatePostContent(id uint, userID int, source string) error {
	var o Post
	if err := DB.Select([]string{"id", "source_type", "is_private"}).Take(&o, "id = ? and user_id = ?", id, userID).Error; err != nil {
		return err
	}
	content := ""
	switch o.SourceType {
	case "plain":
		content = source
	case "markdown":
		unsafe := blackfriday.Run([]byte(source))
		policy := bluemonday.UGCPolicy()
		policy.AllowAttrs("class").
			Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).
			OnElements("code")
		content = string(policy.SanitizeBytes(unsafe))
	default:
		return errors.New("unsupport source type")
	}
	m := map[string]interface{}{
		"source":       source,
		"content":      content,
		"published_at": time.Now(),
	}
	return DB.Model(&Post{}).Where("id = ? and user_id = ?", id, userID).Updates(m).Error
}

func GetPostsForPortal(username, categoryRoute string, limit, page int) (s []*Post) {
	db := DB.Select([]string{
		"posts.id", "title", "posts.description", "user_id", "category_id", "posts.cover", "source_type", "posts.route", "is_published", "published_at", "posts.is_private", "posts.created_at", "posts.updated_at",
	}).
		Joins("LEFT JOIN users on users.id = posts.user_id").
		Joins("LEFT JOIN categories on categories.id = posts.category_id").
		Order("id DESC").Limit(limit).Offset(Offset(limit, page)).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "username", "nickname", "avatar"})
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "category_name", "route", "is_private"})
		}).
		Where("is_published = 1")
	if username != "" {
		db = db.Where("users.username = ?", username)
	}
	if categoryRoute != "" {
		categoryRouteHash := utils.Md5(categoryRoute)
		db = db.Where("categories.route_hash = ?", categoryRouteHash)
	}
	db = db.Find(&s)
	return
}

func CountPostsForPortal(username, categoryRoute string) (count int64) {
	db := DB.Table("posts").Joins("LEFT JOIN users on users.id = posts.user_id").
		Joins("LEFT JOIN categories on categories.id = posts.category_id").
		Where("is_published = 1")
	if username != "" {
		db = db.Where("users.username = ?", username)
	}
	if categoryRoute != "" {
		categoryRouteHash := utils.Md5(categoryRoute)
		db = db.Where("categories.route_hash = ?", categoryRouteHash)
	}
	db = db.Count(&count)
	return
}

func GetPostForPortal(route string) *Post {
	routeHash := utils.Md5(route)
	var v Post
	if err := DB.Select([]string{
		"id", "unique_id", "title", "posts.keywords", "posts.description", "user_id", "category_id", "cover", "source_type", "content", "route", "is_published", "published_at", "is_private", "password", "is_comment_enabled", "is_comment_shown", "posts.created_at", "posts.updated_at",
	}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "username", "nickname", "avatar"})
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "category_name", "route", "is_private"})
		}).
		Where("is_published = 1").
		Take(&v, "route_hash = ?", routeHash).Error; err != nil {
		return nil
	}
	return &v
}

func GetPostForPortalByUniqueID(uniqueID string) *Post {
	var v Post
	if err := DB.Select([]string{
		"id", "unique_id", "title", "posts.keywords", "posts.description", "user_id", "category_id", "cover", "source_type", "content", "route", "is_published", "published_at", "is_private", "password", "is_comment_enabled", "is_comment_shown", "posts.created_at", "posts.updated_at",
	}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "username", "nickname", "avatar"})
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "category_name", "route", "is_private"})
		}).
		Where("is_published = 1").
		Take(&v, "unique_id = ?", uniqueID).Error; err != nil {
		return nil
	}
	return &v
}

func GetPostByRoute(route string) *Post {
	routeHash := utils.Md5(route)
	var v Post
	if err := DB.Select([]string{
		"id", "title", "user_id", "category_id", "cover", "source_type", "route", "is_published", "published_at", "is_private", "is_comment_enabled", "is_comment_shown",
	}).
		Where("is_published = 1").
		Take(&v, "route_hash = ?", routeHash).Error; err != nil {
		return nil
	}
	return &v
}
