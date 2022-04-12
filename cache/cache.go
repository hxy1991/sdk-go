package cache

import (
	"github.com/hxy1991/sdk-go/log"
	"sync"
	"sync/atomic"
)

type Cache struct {
	caches     sync.Map
	cacheLimit int64
	// size is used to count the number elements in the cache.
	// The atomic package is used to ensure this size is accurate when
	// using multiple goroutines.
	size int64
}

func New(cacheLimit int64) *Cache {
	return &Cache{
		cacheLimit: cacheLimit,
		caches:     sync.Map{},
	}
}

func (c *Cache) Get(cacheKey interface{}) (interface{}, bool) {
	cacheValue, found := c.caches.Load(cacheKey)
	return cacheValue, found
}

func (c *Cache) Add(cacheKey, cacheValue interface{}) {
	_, loaded := c.caches.LoadOrStore(cacheKey, cacheValue)
	if loaded {
		// 原来存在这个 Key
		c.caches.Store(cacheKey, cacheValue)
	} else {
		// 原本不存在这个 key
		size := atomic.AddInt64(&c.size, 1)
		if size > 0 && size > c.cacheLimit {
			c.deleteRandomKey(cacheKey)
		}
	}
}

func (c *Cache) deleteRandomKey(butNotKey interface{}) {
	if atomic.LoadInt64(&c.size) > 0 && atomic.LoadInt64(&c.size) > c.cacheLimit {
		c.caches.Range(func(key, value interface{}) bool {
			if butNotKey != nil && key == butNotKey {
				// 不删除刚刚添加的
				return true
			}
			if atomic.LoadInt64(&c.size) > 0 && atomic.LoadInt64(&c.size) > c.cacheLimit {
				log.Warn("exceed the cache limit [", c.cacheLimit, "] delete random key [", key, "]")
				c.Delete(key)
				return true
			}
			return false
		})
	}
}

func (c *Cache) Delete(key interface{}) {
	_, loaded := c.caches.LoadAndDelete(key)
	if loaded {
		atomic.AddInt64(&c.size, -1)
	}
}

func (c *Cache) Keys() []interface{} {
	var keys []interface{}
	c.caches.Range(func(key, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

func (c *Cache) UpdateCacheLimit(cacheLimit int64) int64 {
	oldCacheLimit := atomic.SwapInt64(&c.cacheLimit, cacheLimit)

	if atomic.LoadInt64(&c.size) > 0 && atomic.LoadInt64(&c.size) > c.cacheLimit {
		c.deleteRandomKey(nil)
	}

	return oldCacheLimit
}

func (c *Cache) CacheLimit() int64 {
	return c.cacheLimit
}
