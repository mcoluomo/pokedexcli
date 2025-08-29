package cache

import (
	"sync"
	"time"
)

// CacheEntry represents a single cache entry
type CacheEntry struct {
	data      []byte
	timestamp time.Time
}

// Cache is a simple in-memory cache
type Cache struct {
	entries map[string]CacheEntry
	mutex   sync.RWMutex
	ttl     time.Duration
}

// NewCache creates a new cache instance
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
		mutex:   sync.RWMutex{},
		ttl:     ttl,
	}
}

// Get retrieves data from cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// Check if entry has expired
	if time.Since(entry.timestamp) > c.ttl {
		// Entry expired, remove it
		go func() {
			c.mutex.Lock()
			delete(c.entries, key)
			c.mutex.Unlock()
		}()
		return nil, false
	}

	return entry.data, true
}

// Set stores data in cache
func (c *Cache) Set(key string, data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = CacheEntry{
		data:      data,
		timestamp: time.Now(),
	}
}

// Has checks if a key exists in cache (and is not expired)
func (c *Cache) Has(key string) bool {
	_, exists := c.Get(key)
	return exists
}

// Clear removes all entries from cache
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries = make(map[string]CacheEntry)
}

// Size returns the number of entries in cache
func (c *Cache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.entries)
}

// CleanExpired removes all expired entries
func (c *Cache) CleanExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.Sub(entry.timestamp) > c.ttl {
			delete(c.entries, key)
		}
	}
}
