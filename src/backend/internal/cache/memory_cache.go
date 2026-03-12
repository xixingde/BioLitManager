package cache

import (
	"sync"
	"time"

	"biolitmanager/pkg/logger"
	"go.uber.org/zap"
)

var (
	globalCache *MemoryCache
	cacheOnce   sync.Once
)

// Item 缓存项
type Item struct {
	value      interface{}
	expiration int64
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	items map[string]Item
	mu    sync.RWMutex
}

// NewMemoryCache 创建新的内存缓存实例
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]Item),
	}
	go cache.cleanupExpired()
	return cache
}

// Set 设置缓存
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	c.items[key] = Item{
		value:      value,
		expiration: expiration,
	}

	logger.GetLogger().Debug("Cache item set",
		zap.String("key", key),
		zap.Duration("ttl", ttl),
		zap.Int64("expiration", expiration),
	)
}

// Get 获取缓存
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		logger.GetLogger().Debug("Cache item not found",
			zap.String("key", key),
		)
		return nil, false
	}

	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		logger.GetLogger().Debug("Cache item expired",
			zap.String("key", key),
			zap.Int64("expiration", item.expiration),
			zap.Int64("now", time.Now().UnixNano()),
		)
		return nil, false
	}

	logger.GetLogger().Debug("Cache item found",
		zap.String("key", key),
	)
	return item.value, true
}

// Delete 删除缓存
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// cleanupExpired 定期清理过期缓存
func (c *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for key, item := range c.items {
			if item.expiration > 0 && now > item.expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// GetInstance 获取全局缓存实例（单例模式）
func GetInstance() *MemoryCache {
	if globalCache == nil {
		cacheOnce.Do(func() {
			globalCache = NewMemoryCache()
			logger.GetLogger().Info("Global cache instance created")
		})
	}
	logger.GetLogger().Debug("Global cache instance requested")
	return globalCache
}

// InitCache 初始化全局缓存
func InitCache() *MemoryCache {
	// 调用 GetInstance 来确保使用单例模式
	instance := GetInstance()
	logger.GetLogger().Info("Global cache initialized via InitCache")
	return instance
}
