package config

import (
	"log"

	"github.com/spf13/viper"
)

var globalConfig *Config

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Cache    CacheConfig    `mapstructure:"cache"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	SessionTTL  string `mapstructure:"session_ttl"`
	UserInfoTTL string `mapstructure:"user_info_ttl"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, err
	}

	globalConfig = config
	return config, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// InitConfig 初始化配置
func InitConfig() error {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return err
	}
	globalConfig = config
	log.Printf("Config loaded successfully")
	return nil
}
