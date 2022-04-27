package models

import (
	"time"

	"elf-server/pkg/utils"
)

const (
	ReaderGenderMale   = uint8(0)
	ReaderGenderFemale = uint8(1)
)

type Reader struct {
	ID        uint       `json:"id" gorm:"primary_key;autoIncrement"`
	UniqueID  string     `json:"uniqueId" gorm:"size:50;unique"`
	Nickname  string     `json:"nickname" gorm:"not null;size:128"`
	Password  string     `json:"-" gorm:"size:50"` // no use
	Gender    uint8      `json:"gender" gorm:"not null"`
	Birthday  *time.Time `json:"birthday" gorm:"type:date"`
	Email     string     `json:"email" gorm:"not null;size:255"`
	EmailHash string     `json:"-" gorm:"size:32;unique"`
	Phone     string     `json:"phone" gorm:"size:50"`
	UserID    *uint      `json:"userId" gorm:"unique"`
	OpenID    string     `json:"openId" gorm:"size:255"`
	IsActive  bool       `json:"isActive" gorm:"not null;index"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

var (
	readersFields = []string{"id", "unique_id", "nickname", "gender", "email", "phone", "user_id", "is_active", "created_at", "updated_at"}
	readerFields  = []string{"id", "unique_id", "nickname", "gender", "birthday", "email", "phone", "user_id", "is_active", "created_at", "updated_at"}
)

func GetReaders(limit, page int) (s []*Reader, count int64) {
	DB.Model(&Reader{}).
		Select(readersFields).
		Count(&count).
		Order("id ASC").Limit(limit).Offset(Offset(limit, page)).
		Find(&s)
	return
}

func GetReader(id uint) *Reader {
	var v Reader
	if err := DB.Select(readerFields).Take(&v, "id = ?", id).Error; err != nil {
		return nil
	}
	return &v
}

func (v *Reader) Save() error {
	v.UniqueID = utils.NewID()
	v.EmailHash = utils.Md5(v.Email)
	return DB.Save(v).Error
}

func (v *Reader) Update() error {
	return DB.Model(v).Updates(map[string]interface{}{
		"nickname":   v.Nickname,
		"gender":     v.Gender,
		"birthday":   v.Birthday,
		"email":      v.Email,
		"email_hash": utils.Md5(v.Email),
		"phone":      v.Phone,
		"is_active":  v.IsActive,
	}).Error
}

func (v *Reader) Delete() error {
	return DB.Delete(v).Error
}

func GetReaderByEmail(email string) *Reader {
	var v Reader
	if err := DB.Select([]string{"id", "unique_id", "nickname", "gender", "birthday", "email", "phone", "is_active"}).
		Take(&v, "email_hash = ?", utils.Md5(email)).Error; err != nil {
		return nil
	}
	return &v
}

func GetReaderByNickname(nickname string) *Reader {
	var v Reader
	if err := DB.Select([]string{"id", "unique_id", "nickname", "gender", "birthday", "email", "phone", "is_active"}).
		Take(&v, "nickname = ?", nickname).Error; err != nil {
		return nil
	}
	return &v
}
