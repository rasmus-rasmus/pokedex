package cache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.cache[key]
	if ok {
		return
	}
	entry := cacheEntry{
		time.Now(),
		val,
	}
	c.cache[key] = entry
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cache[key]
	return entry.val, ok
}

func (c Cache) reapLoop(interval time.Duration, done chan struct{}) {
	ticker := time.NewTicker(interval)
	defer func() { <-done }() // Telling caller that reapLoop returned
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			c.mu.Lock()
			for key, val := range c.cache {
				if time.Since(val.createdAt) > interval {
					delete(c.cache, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

func NewCache(interval time.Duration) (*Cache, chan struct{}) {
	c := &Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}
	done := make(chan struct{})
	go c.reapLoop(interval, done)
	return c, done
}
