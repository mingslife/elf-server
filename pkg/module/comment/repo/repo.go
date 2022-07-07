package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/importcjj/sensitive"
	"github.com/mingslife/bone"

	"elf-server/pkg/component"
	"elf-server/pkg/module/comment/model"
	"elf-server/pkg/utils"
)

type CommentRepo struct {
	Database *component.Database `inject:"component.database"`
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

func (r *CommentRepo) List(ctx context.Context, limit, page int) (s []*model.Comment, cnt int64, err error) {
	err = r.Database.DB.WithContext(ctx).Model(&model.Comment{}).
		Count(&cnt).
		Select(commentsFields).
		Order("id ASC").Limit(limit).Offset(utils.Offset(limit, page)).
		Find(&s).Error
	return
}

func (r *CommentRepo) Get(ctx context.Context, id uint) (*model.Comment, error) {
	var v model.Comment
	if err := r.Database.DB.WithContext(ctx).Select(commentFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *CommentRepo) Create(ctx context.Context, v *model.Comment) error {
	var maxLevel uint
	v.Content = sensitiveFilter.Replace(v.Content, '*')
	r.Database.DB.WithContext(ctx).Model(&model.Comment{}).
		Select("max(level)").
		Where("post_id = ?", v.PostID).
		Scan(&maxLevel)
	v.Level = maxLevel + 1
	v.CommentedAt = time.Now()
	return r.Database.DB.WithContext(ctx).Save(v).Error
}

func (r *CommentRepo) Update(ctx context.Context, v *model.Comment) error {
	return r.Database.DB.WithContext(ctx).Model(v).Updates(map[string]any{
		"post_id":      v.PostID,
		"ip":           v.IP,
		"user_agent":   v.UserAgent,
		"content":      sensitiveFilter.Replace(v.Content, '*'),
		"commented_at": v.CommentedAt,
	}).Error
}

func (r *CommentRepo) Delete(ctx context.Context, v *model.Comment) error {
	return r.Database.DB.WithContext(ctx).Delete(v).Error
}

func (r *CommentRepo) ListByPostID(ctx context.Context, postID int, limit, page int) (s []*model.Comment, cnt int64, err error) {
	db := r.Database.DB.WithContext(ctx).Model(&model.Comment{}).
		Select(commentsFields)
	if postID != -1 {
		db = db.Where("post_id = ?", postID)
	}
	db = db.Count(&cnt).
		Order("level DESC").
		Limit(limit).Offset(utils.Offset(limit, page)).
		Find(&s)
	err = db.Error
	return
}

func (r *CommentRepo) CountByPostID(ctx context.Context, postID int) (cnt int64, err error) {
	db := r.Database.DB.WithContext(ctx).Model(&model.Comment{})
	if postID != -1 {
		db = db.Where("post_id = ?", postID)
	}
	db = db.Count(&cnt)
	err = db.Error
	return
}

func (r *CommentRepo) GetByPostIDAndLevel(ctx context.Context, postID, level uint) (*model.Comment, error) {
	var v model.Comment
	err := r.Database.DB.WithContext(ctx).Model(&model.Comment{}).Select(commentFields).
		Where("post_id = ? and level = ?", postID, level).
		Scan(&v).Error
	return &v, err
}

func (r *CommentRepo) ListByPostIDForPortal(ctx context.Context, postID uint) (s []*model.PortalComment, cnt int64, err error) {
	r.Database.DB.WithContext(ctx).Raw(fmt.Sprintf(`
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

	err = r.Database.DB.WithContext(ctx).Model(&model.Comment{}).
		Where("post_id = ?", postID).
		Count(&cnt).Error
	return
}

func (r *CommentRepo) SetIsBlocked(ctx context.Context, id uint, isBlocked bool) error {
	return r.Database.DB.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", id).
		Update("is_blocked", isBlocked).Error
}

var _ bone.Repo = (*CommentRepo)(nil)
