package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 配置结构
type Config struct {
	// 数据库配置
	DBPath string

	// 服务器配置
	Port    string
	GinMode string

	// 导入导出配置
	ImportExportDir string

	// 日志配置
	LogLevel string
}

// 全局配置实例
var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 尝试加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用环境变量")
	}

	config := &Config{
		DBPath:          getEnv("DB_PATH", "./flashcard.db"),
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		ImportExportDir: getEnv("IMPORT_EXPORT_DIR", "./data"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}

	AppConfig = config
	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量作为整数
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量作为布尔值
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
