package repo

import (
	"context"
	"errors"
	"time"

	"github.com/mingslife/bone"

	"elf-server/pkg/component"
	"elf-server/pkg/module/user/model"
	"elf-server/pkg/utils"
)

type UserRepo struct {
	Database *component.Database `inject:"component.database"`
}

var (
	usersFields = []string{"id", "username", "nickname", "email", "phone", "is_active", "active_at", "gender", "role", "created_at", "updated_at"}
	userFields  = []string{"id", "username", "nickname", "email", "phone", "tags", "introduction", "is_active", "avatar", "gender", "birthday", "role", "created_at", "updated_at"}
)

func (r *UserRepo) List(ctx context.Context, limit, page int) (s []*model.User, cnt int64, err error) {
	err = r.Database.DB.WithContext(ctx).Model(&model.User{}).
		Select(usersFields).
		Count(&cnt).
		Order("id ASC").Limit(limit).Offset(utils.Offset(limit, page)).
		Find(&s).Error
	return
}

func (r *UserRepo) Get(ctx context.Context, id uint) (*model.User, error) {
	var v model.User
	if err := r.Database.DB.WithContext(ctx).Select(userFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *UserRepo) Create(ctx context.Context, v *model.User) error {
	v.Password = utils.GenerateFromPassword(v.Password)
	v.EmailHash = utils.Md5(v.Email)
	v.PhoneHash = utils.Md5(v.Phone)
	if v.IsActive {
		now := time.Now()
		v.ActiveAt = &now
	}
	return r.Database.DB.WithContext(ctx).Save(v).Error
}

func (r *UserRepo) Update(ctx context.Context, v *model.User) error {
	m := map[string]any{
		"username":     v.Username,
		"nickname":     v.Nickname,
		"email":        v.Email,
		"email_hash":   utils.Md5(v.Email),
		"phone":        v.Phone,
		"phone_hash":   utils.Md5(v.Phone),
		"tags":         v.Tags,
		"introduction": v.Introduction,
		"is_active":    v.IsActive,
		"avatar":       v.Avatar,
		"gender":       v.Gender,
		"birthday":     v.Birthday,
	}
	if v.Password != "" {
		m["password"] = utils.GenerateFromPassword(v.Password)
	}
	if v.IsActive {
		var o model.User
		if err := r.Database.DB.WithContext(ctx).Select([]string{"id", "is_active"}).Take(&o, "id = ?", v.ID).Error; err != nil {
			return err
		}
		if !o.IsActive {
			m["active_at"] = time.Now()
		}
	} else {
		m["active_at"] = nil
	}
	return r.Database.DB.WithContext(ctx).Model(v).Updates(m).Error
}

func (r *UserRepo) Delete(ctx context.Context, v *model.User) error {
	return r.Database.DB.WithContext(ctx).Delete(v).Error
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var v model.User
	if err := r.Database.DB.WithContext(ctx).Select(userFields).Take(&v, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *UserRepo) GetByAccountAndPassword(ctx context.Context, account, password string) (*model.User, error) {
	var v model.User
	accountHash := utils.Md5(account)
	if err := r.Database.DB.WithContext(ctx).Select([]string{"id", "username", "password", "is_active", "role"}).Take(&v, "(username = ? or email_hash = ? or phone_hash = ?) and is_active = 1", account, accountHash, accountHash).Error; err != nil {
		return nil, err
	}
	if !utils.MatchPassword(v.Password, password) {
		return nil, errors.New("incorrect account or password")
	}
	return &v, nil
}

func (r *UserRepo) Exists(ctx context.Context, username, email, phone string) bool {
	var v model.User
	emailHash, phoneHash := utils.Md5(email), utils.Md5(phone)
	r.Database.DB.WithContext(ctx).Model(&model.User{}).
		Select([]string{"id"}).
		Take(&v, "username = ? or email_hash = ? or phone_hash = ?", username, emailHash, phoneHash)
	return v.ID != 0
}

var _ bone.Repo = (*UserRepo)(nil)
