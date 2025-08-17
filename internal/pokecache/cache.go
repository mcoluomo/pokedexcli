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
