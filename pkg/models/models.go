package models

import (
	"elf-server/pkg/utils"
	"fmt"
	"log"
	"time"

	"github.com/golang/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Model struct {
	ID        int        `json:"id" binding:"-" gorm:"size:36;primary_key"`
	CreatedAt time.Time  `json:"-"  binding:"-" gorm:"index:idx_created_at"`
	UpdatedAt time.Time  `json:"-"  binding:"-"`
	DeletedAt *time.Time `json:"-"  binding:"-" gorm:"index:idx_deleted_at"`
}

var DB *gorm.DB

func InitDB(host string, port int, user string, pwd string, database string, debug bool) {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		user, pwd, host, port, database,
	)
	logMode := logger.Info
	if debug {
		logMode = logger.Warn
	}
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		glog.Fatalf("Failed to connect database: %v", err)
		return
	}

	// defer DB.Close()

	sqlDB, err := DB.DB()
	if err != nil {
		glog.Fatalf("Failed to connect database: %v", err)
		return
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
}

func Seed() {
	now := time.Now()

	DB.AutoMigrate(&Setting{})
	DB.AutoMigrate(&Navigation{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Category{})
	DB.AutoMigrate(&Post{})
	DB.AutoMigrate(&Comment{})

	if GetSettingByKey("app.title") != nil {
		return
	}

	(&Setting{
		SettingKey:   "app.title",
		SettingValue: "ELF",
		SettingTag:   "Title",
		IsPublic:     false,
	}).Save()
	(&Setting{
		SettingKey:   "app.keywords",
		SettingValue: "blogging,go,elf",
		SettingTag:   "Keywords",
		IsPublic:     false,
	}).Save()
	(&Setting{
		SettingKey:   "app.description",
		SettingValue: "A blogging system written in Go",
		SettingTag:   "Description",
		IsPublic:     false,
	}).Save()
	(&Setting{
		SettingKey:   "app.language",
		SettingValue: "en",
		SettingTag:   "Language",
		IsPublic:     true,
	}).Save()
	(&Setting{
		SettingKey:   "app.footer",
		SettingValue: "&copy; 2021 ELFIN",
		SettingTag:   "Footer",
		IsPublic:     false,
	}).Save()
	(&Setting{
		SettingKey:   "app.limit",
		SettingValue: "10",
		SettingTag:   "Limit",
		IsPublic:     false,
	}).Save()
	(&Setting{
		SettingKey:   "app.comment",
		SettingValue: "true",
		SettingTag:   "Comment",
		IsPublic:     false,
	}).Save()
	(&Setting{
		SettingKey:   "app.script",
		SettingValue: "console.log('Powered by ELF')",
		SettingTag:   "Script",
		IsPublic:     false,
	}).Save()
	(&Setting{
		SettingKey:   "app.register",
		SettingValue: "true",
		SettingTag:   "Register",
		IsPublic:     true,
	}).Save()
	(&Setting{
		SettingKey:   "app.inviteCode",
		SettingValue: "<script>console.log('Powered by ELF')</script>",
		SettingTag:   utils.RandString(16),
		IsPublic:     false,
	}).Save()

	adminPassword := utils.RandString(8)
	(&User{
		Username:     "admin",
		Password:     adminPassword,
		Nickname:     "Admin",
		Email:        "admin@elf",
		Phone:        "0000",
		Tags:         "admin",
		Introduction: "Adminstrator account",
		IsActive:     true,
		ActiveAt:     &now,
		Avatar:       "/assets/avatar.svg",
		Gender:       UserGenderMale,
		Birthday:     &now,
		Role:         UserRoleAdmin,
	}).Save()
	log.Printf("admin password: %s\n", adminPassword)
}

func Offset(limit, page int) int {
	if page != -1 && limit != -1 {
		return (page - 1) * limit
	}
	return -1
}
