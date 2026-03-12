package cache

import (
	"fmt"
)

// SessionKey 生成会话缓存键
func SessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}

// UserKey 生成用户缓存键
func UserKey(userId string) string {
	return fmt.Sprintf("user:%s", userId)
}

// ConfigKey 生成配置缓存键
func ConfigKey(key string) string {
	return fmt.Sprintf("config:%s", key)
}
