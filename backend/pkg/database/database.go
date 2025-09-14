package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"flashcard/internal/config"
	"flashcard/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	cfg := config.AppConfig

	// 确保数据库文件目录存在
	dbDir := filepath.Dir(cfg.DBPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %v", err)
	}

	// 设置日志级别
	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Warn
	default:
		logLevel = logger.Error
	}

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	DB = db

	// 自动迁移数据库结构
	if err := migrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	// 创建索引
	if err := createIndexes(); err != nil {
		return fmt.Errorf("创建索引失败: %v", err)
	}

	log.Println("数据库初始化成功")
	return nil
}

// migrate 执行数据库迁移
func migrate() error {
	return DB.AutoMigrate(
		&models.Deck{},
		&models.Tag{},
		&models.Card{},
		&models.Review{},
	)
}

// createIndexes 创建必要的索引
func createIndexes() error {
	// 为tags表创建复合唯一索引（deck_id + name）
	if err := DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_deck_name ON tags(deck_id, name) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// 为cards表创建复合索引用于搜索
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_cards_deck_tag ON cards(deck_id, tag_id) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	// 为reviews表创建next_review索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_reviews_next_review ON reviews(next_review)").Error; err != nil {
		return err
	}

	// 为cards表创建全文搜索索引（可选，SQLite FTS）
	// 暂时使用普通索引，后续可优化为FTS
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_cards_question ON cards(question) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
