package repo

import (
	"context"

	"elf-server/pkg/component"
	"elf-server/pkg/module/reader/model"
	"elf-server/pkg/utils"

	"github.com/mingslife/bone"
)

type ReaderRepo struct {
	Database *component.Database `inject:"component.database"`
}

var (
	readersFields = []string{"id", "unique_id", "nickname", "gender", "email", "phone", "user_id", "is_active", "created_at", "updated_at"}
	readerFields  = []string{"id", "unique_id", "nickname", "gender", "birthday", "email", "phone", "user_id", "is_active", "created_at", "updated_at"}
)

func (r *ReaderRepo) List(ctx context.Context, limit, page int) (s []*model.Reader, cnt int64, err error) {
	err = r.Database.DB.WithContext(ctx).Model(&model.Reader{}).
		Select(readersFields).
		Count(&cnt).
		Order("id ASC").Limit(limit).Offset(utils.Offset(limit, page)).
		Find(&s).Error
	return
}

func (r *ReaderRepo) Get(ctx context.Context, id uint) (*model.Reader, error) {
	var v model.Reader
	if err := r.Database.DB.WithContext(ctx).Select(readerFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *ReaderRepo) Create(ctx context.Context, v *model.Reader) error {
	v.UniqueID = utils.NewID()
	v.EmailHash = utils.Md5(v.Email)
	return r.Database.DB.WithContext(ctx).Save(v).Error
}

func (r *ReaderRepo) Update(ctx context.Context, v *model.Reader) error {
	return r.Database.DB.WithContext(ctx).Model(v).Updates(map[string]any{
		"nickname":   v.Nickname,
		"gender":     v.Gender,
		"birthday":   v.Birthday,
		"email":      v.Email,
		"email_hash": utils.Md5(v.Email),
		"phone":      v.Phone,
		"is_active":  v.IsActive,
	}).Error
}

func (r *ReaderRepo) Delete(ctx context.Context, v *model.Reader) error {
	return r.Database.DB.WithContext(ctx).Delete(v).Error
}

func (r *ReaderRepo) GetByEmail(ctx context.Context, email string) (*model.Reader, error) {
	var v model.Reader
	if err := r.Database.DB.WithContext(ctx).Select([]string{"id", "unique_id", "nickname", "gender", "birthday", "email", "phone", "is_active"}).
		Take(&v, "email_hash = ?", utils.Md5(email)).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *ReaderRepo) GetByNickname(ctx context.Context, nickname string) (*model.Reader, error) {
	var v model.Reader
	if err := r.Database.DB.WithContext(ctx).Select([]string{"id", "unique_id", "nickname", "gender", "birthday", "email", "phone", "is_active"}).
		Take(&v, "nickname = ?", nickname).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *ReaderRepo) GetByUniqueID(ctx context.Context, uniqueID string) (*model.Reader, error) {
	var v model.Reader
	if err := r.Database.DB.WithContext(ctx).Select([]string{"id", "unique_id", "nickname", "gender", "birthday", "email", "phone", "is_active"}).
		Take(&v, "unique_id = ?", uniqueID).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *ReaderRepo) GetByUserID(ctx context.Context, userID uint) (*model.Reader, error) {
	var v model.Reader
	if err := r.Database.DB.WithContext(ctx).Select([]string{"id", "unique_id", "nickname", "gender", "birthday", "email", "phone", "is_active"}).
		Take(&v, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

var _ bone.Repo = (*ReaderRepo)(nil)
