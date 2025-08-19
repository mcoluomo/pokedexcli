package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entry map[string]cacheEntry
	mutex *sync.Mutex
}

type cacheEntry struct {
	createAt time.Time
	val      []byte
}

func NewCache(interval time.Time) Cache {
	reapLoop()
	return Cache{}
}

func (c *Cache) Add(url string, data []byte) {
	c.Entry[url] = cacheEntry{val: data}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	_, ok := c.Entry[key]
	var data []byte
	if !ok {
		return data, ok
	}

	return c.Entry[key].val, ok
}

func (c *Cache) reapLoop() {
}
