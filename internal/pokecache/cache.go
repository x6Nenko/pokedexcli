package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntry map[string]entry
	mu    		 sync.Mutex
}

type entry struct {
	createdAt       time.Time
	val  						[]byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cacheEntry: make(map[string]entry),
	}

	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	newEntry := entry{
		createdAt: time.Now(),
		val: val,
	} 

	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheEntry[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, found := c.cacheEntry[key]
	return value.val, found
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()  // Clean up when done
	
	for {
		<-ticker.C  // Wait for tick
		
		// Lock ONLY during the cleanup work
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.cacheEntry {
				if now.Sub(item.createdAt) > interval {
					delete(c.cacheEntry, key)
				}
		}
		c.mu.Unlock()  // Unlock immediately after
	}
}