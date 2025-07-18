package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries		map[string]cacheEntry
	mu			sync.Mutex
	interval	time.Duration
}

type cacheEntry struct {
	createdAt	time.Time
	val			[]byte
}

func NewCache(t time.Duration) *Cache {
	newCache := &Cache{
		entries: make(map[string]cacheEntry),
		interval: t,
	}
	go newCache.reapLoop()
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	newEntry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	c.entries[key] = newEntry
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cache, ok := c.entries[key]; ok {
		return cache.val, true
	}
	return []byte{}, false
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)

	//start an anonymous func concurrently
	go func() {
		defer ticker.Stop()
		for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.entries {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.entries, key)
				}
			}
			c.mu.Unlock()
		}
	}
	}()
}
