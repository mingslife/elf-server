package repo

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mingslife/bone"
	"github.com/russross/blackfriday/v2"
	"gorm.io/gorm"

	"elf-server/pkg/component"
	"elf-server/pkg/module/post/model"
	"elf-server/pkg/utils"
)

type PostRepo struct {
	Database *component.Database `inject:"component.database"`
}

var (
	postsFields = []string{"id", "title", "user_id", "category_id", "source_type", "route", "is_published", "published_at", "is_private", "is_comment_enabled", "is_comment_shown", "created_at", "updated_at"}
	postFields  = []string{"id", "title", "keywords", "description", "user_id", "category_id", "cover", "source_type", "route", "is_published", "is_private", "password", "is_comment_enabled", "is_comment_shown", "created_at", "updated_at"}
)

func (r *PostRepo) List(ctx context.Context, limit, page int) (s []*model.Post, cnt int64, err error) {
	err = r.Database.DB.WithContext(ctx).Model(&model.Post{}).
		Select(postsFields).
		Count(&cnt).
		Order("id DESC").Limit(limit).Offset(utils.Offset(limit, page)).
		Find(&s).Error
	return
}

func (r *PostRepo) Get(ctx context.Context, id uint) (*model.Post, error) {
	var v model.Post
	if err := r.Database.DB.WithContext(ctx).Select(postFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PostRepo) Create(ctx context.Context, v *model.Post) error {
	v.RouteHash = utils.Md5(v.Route)
	if v.IsPublished {
		now := time.Now()
		v.PublishedAt = &now
	}
	return r.Database.DB.WithContext(ctx).Save(v).Error
}

func (r *PostRepo) Update(ctx context.Context, v *model.Post) error {
	m := map[string]any{
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
		var o model.Post
		if err := r.Database.DB.WithContext(ctx).Select([]string{"id", "is_published"}).Take(&o, "id = ?", v.ID).Error; err != nil {
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
	return r.Database.DB.WithContext(ctx).Model(v).Updates(m).Error
}

func (r *PostRepo) Delete(ctx context.Context, v *model.Post) error {
	return r.Database.DB.WithContext(ctx).Delete(v).Error
}

func (r *PostRepo) ListByUserID(ctx context.Context, userID uint, limit, page int) (s []*model.Post, cnt int64, err error) {
	err = r.Database.DB.WithContext(ctx).Model(&model.Post{}).Select([]string{
		"posts.id", "posts.unique_id", "page_view", "thumb_up", "thumb_down", "title", "user_id", "nickname", "category_id", "category_name", "source_type", "posts.route", "is_published", "published_at", "posts.is_private", "is_comment_enabled", "is_comment_shown", "posts.created_at", "posts.updated_at",
	}).
		Where("user_id = ?", userID).
		Count(&cnt).
		Joins("LEFT JOIN users on users.id = posts.user_id").
		Joins("LEFT JOIN categories on categories.id = posts.category_id").
		Joins("LEFT JOIN post_statistics on post_statistics.unique_id = posts.unique_id").
		Order("id DESC").Limit(limit).Offset(utils.Offset(limit, page)).
		Preload("PostStatistics", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "unique_id", "page_view", "thumb_up", "thumb_down"})
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "username", "nickname", "avatar"})
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select([]string{"id", "category_name", "is_private"})
		}).
		Find(&s).Error
	return
}

func (r *PostRepo) GetContent(ctx context.Context, id uint, userID int) (*model.Post, error) {
	var v model.Post
	if err := r.Database.DB.WithContext(ctx).Select([]string{
		"id", "source", "source_type", "content",
	}).Take(&v, "id = ? and user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *PostRepo) UpdateContent(ctx context.Context, id uint, userID int, source string) error {
	var o model.Post
	if err := r.Database.DB.WithContext(ctx).Select([]string{"id", "source_type", "is_private"}).Take(&o, "id = ? and user_id = ?", id, userID).Error; err != nil {
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
	return r.Database.DB.WithContext(ctx).Model(&model.Post{}).Where("id = ? and user_id = ?", id, userID).Updates(m).Error
}

func (r *PostRepo) ListForPortal(ctx context.Context, username, categoryRoute string, limit, page int) (s []*model.Post, err error) {
	db := r.Database.DB.WithContext(ctx).Select([]string{
		"posts.id", "title", "posts.description", "user_id", "category_id", "posts.cover", "source_type", "posts.route", "is_published", "published_at", "posts.is_private", "posts.created_at", "posts.updated_at",
	}).
		Joins("LEFT JOIN users on users.id = posts.user_id").
		Joins("LEFT JOIN categories on categories.id = posts.category_id").
		Order("id DESC").Limit(limit).Offset(utils.Offset(limit, page)).
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
	err = db.Error
	return
}

func (r *PostRepo) CountForPortal(ctx context.Context, username, categoryRoute string) (cnt int64, err error) {
	db := r.Database.DB.WithContext(ctx).Table("posts").Joins("LEFT JOIN users on users.id = posts.user_id").
		Joins("LEFT JOIN categories on categories.id = posts.category_id").
		Where("is_published = 1")
	if username != "" {
		db = db.Where("users.username = ?", username)
	}
	if categoryRoute != "" {
		categoryRouteHash := utils.Md5(categoryRoute)
		db = db.Where("categories.route_hash = ?", categoryRouteHash)
	}
	db = db.Count(&cnt)
	err = db.Error
	return
}

func (r *PostRepo) GetForPortal(ctx context.Context, route string) (*model.Post, error) {
	routeHash := utils.Md5(route)
	var v model.Post
	if err := r.Database.DB.WithContext(ctx).Select([]string{
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
		return nil, err
	}
	return &v, nil
}

func (r *PostRepo) GetForPortalByUniqueID(ctx context.Context, uniqueID string) (*model.Post, error) {
	var v model.Post
	if err := r.Database.DB.WithContext(ctx).Select([]string{
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
		return nil, err
	}
	return &v, nil
}

func (r *PostRepo) GetByRoute(ctx context.Context, route string) (*model.Post, error) {
	routeHash := utils.Md5(route)
	var v model.Post
	if err := r.Database.DB.WithContext(ctx).Select([]string{
		"id", "title", "user_id", "category_id", "cover", "source_type", "route", "is_published", "published_at", "is_private", "is_comment_enabled", "is_comment_shown",
	}).
		Where("is_published = 1").
		Take(&v, "route_hash = ?", routeHash).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

var _ bone.Repo = (*PostRepo)(nil)
