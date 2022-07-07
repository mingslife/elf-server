package component

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mingslife/bone"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"elf-server/pkg/conf"
)

type DatabaseLogger struct {
	Logger bone.Logger `inject:"component.logger"`
}

func (l *DatabaseLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *DatabaseLogger) Info(ctx context.Context, msg string, data ...any) {
	l.Logger.WithContext(ctx).Info(fmt.Sprintf(msg, data...))
}

func (l *DatabaseLogger) Warn(ctx context.Context, msg string, data ...any) {
	l.Logger.WithContext(ctx).Warn(fmt.Sprintf(msg, data...))
}

func (l *DatabaseLogger) Error(ctx context.Context, msg string, data ...any) {
	l.Logger.WithContext(ctx).Error(fmt.Sprintf(msg, data...))
}

func (l *DatabaseLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	traceStr := "[%s] [%.3fms] [rows:%v] %s"
	source := utils.FileWithLineNum()
	elapsed := time.Since(begin)
	sql, rows := fc()
	var rowsStr string
	if rows == -1 {
		rowsStr = "-"
	} else {
		rowsStr = strconv.FormatInt(rows, 10)
	}
	if err != nil {
		l.Logger.WithContext(ctx).Error(fmt.Sprintf(traceStr, source, float64(elapsed.Nanoseconds())/1e6, rowsStr, sql), map[string]any{
			"error": err,
		})
	} else {
		l.Logger.WithContext(ctx).Info(fmt.Sprintf(traceStr, source, float64(elapsed.Nanoseconds())/1e6, rowsStr, sql))
	}
}

var _ logger.Interface = (*DatabaseLogger)(nil)

type Database struct {
	Logger *DatabaseLogger `inject:""`
	DB     *gorm.DB
}

func (*Database) Name() string {
	return "component.database"
}

func (*Database) Init() error {
	return nil
}

func (c *Database) Register() error {
	cfg := conf.GetConfig()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DbUser,
		cfg.DbPwd,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)
	var err error
	c.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: c.Logger,
	})
	return err
}

func (c *Database) Unregister() error {
	return nil
}

var _ bone.Component = (*Database)(nil)
