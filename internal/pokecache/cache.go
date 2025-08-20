package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entry map[string]cacheEntry
	mutex sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{Entry: make(map[string]cacheEntry)}
	go cache.reapLoop(interval)
	return &cache
}

func (c *Cache) Add(url string, data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Entry[url] = cacheEntry{createdAt: time.Now(), val: data}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Unlock()

	_, ok := c.Entry[key]
	var data []byte

	defer c.mutex.Unlock()

	if !ok {
		return data, ok
	}

	return c.Entry[key].val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.Tick(interval)
	for range ticker {
		c.mutex.Lock()
		for key := range c.Entry {
			if time.Since(c.Entry[key].createdAt) > interval {
				delete(c.Entry, key)
			}
		}
		c.mutex.Unlock()
	}
}
