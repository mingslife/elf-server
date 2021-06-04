package models

import (
	"time"

	"elf-server/pkg/utils"
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

var (
	usersFields = []string{"id", "username", "nickname", "email", "phone", "is_active", "active_at", "gender", "role", "created_at", "updated_at"}
	userFields  = []string{"id", "username", "nickname", "email", "phone", "tags", "introduction", "is_active", "avatar", "gender", "birthday", "role", "created_at", "updated_at"}
)

func GetUsers(limit, page int) (s []*User, count int64) {
	DB.Model(&User{}).
		Select(usersFields).
		Count(&count).
		Order("id ASC").Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func GetUser(id uint) *User {
	var v User
	if err := DB.Select(userFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil
	}
	return &v
}

func (v *User) Save() error {
	v.Password = utils.GenerateFromPassword(v.Password)
	v.EmailHash = utils.Md5(v.Email)
	v.PhoneHash = utils.Md5(v.Phone)
	if v.IsActive {
		now := time.Now()
		v.ActiveAt = &now
	}
	return DB.Save(v).Error
}

func (v *User) Update() error {
	m := map[string]interface{}{
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
		var o User
		if err := DB.Select([]string{"id", "is_active"}).Take(&o, "id = ?", v.ID).Error; err != nil {
			return err
		}
		if !o.IsActive {
			m["active_at"] = time.Now()
		}
	} else {
		m["active_at"] = nil
	}
	return DB.Model(v).Updates(m).Error
}

func (v *User) Delete() error {
	return DB.Delete(v).Error
}

func GetUserByUsername(username string) *User {
	var v User
	if err := DB.Select(userFields).Take(&v, "username = ?", username).Error; err != nil {
		return nil
	}
	return &v
}

func GetUserByAccountAndPassword(account, password string) *User {
	var v User
	accountHash := utils.Md5(account)
	if err := DB.Select([]string{"id", "username", "password", "is_active", "role"}).Take(&v, "(username = ? or email_hash = ? or phone_hash = ?) and is_active = 1", account, accountHash, accountHash).Error; err != nil {
		return nil
	}
	if !utils.MatchPassword(v.Password, password) {
		return nil
	}
	return &v
}

func ExistsUser(username, email, phone string) bool {
	var v User
	emailHash, phoneHash := utils.Md5(email), utils.Md5(phone)
	DB.Model(&User{}).
		Select([]string{"id"}).
		Take(&v, "username = ? or email_hash = ? or phone_hash = ?", username, emailHash, phoneHash)
	return v.ID != 0
}
