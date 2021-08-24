package models

import (
	"fmt"
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

func Offset(limit, page int) int {
	if page != -1 && limit != -1 {
		return (page - 1) * limit
	}
	return -1
}
