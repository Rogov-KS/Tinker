package database

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Connect(databaseURL string) error {
	var err error

	// Настройка GORM logger с детальным логированием SQL
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		},
	)

	// NamingStrategy который сохраняет точные имена колонок из gorm тегов
	namingStrategy := schema.NamingStrategy{
		TablePrefix:   "",
		SingularTable: false,
		NameReplacer:  nil,
		NoLowerCase:   true, // Не преобразуем в lowercase
	}

	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger:          gormLogger,
		NamingStrategy: namingStrategy,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Настройка connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Проверка подключения
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Database connected successfully")
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

