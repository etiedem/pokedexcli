package main

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	Entries  map[string]cacheEntry
	interval time.Duration
	Lock     sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.Entries[key] = cacheEntry{time.Now(), val}
}
func (c *Cache) Get(key string) ([]byte, bool) {
	entry, err := c.Entries[key]
	return entry.val, err
}
func (c *Cache) reapLoop() {
	remove := []string{}
	for key, entry := range c.Entries {
		if time.Now().Add(-c.interval).Before(entry.createdAt) {
			remove = append(remove, key)
		}
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()
	for _, item := range remove {
		delete(c.Entries, item)
	}
}

func NewCache(d time.Duration) *Cache {
	ticker := time.NewTicker(d)
	cache := Cache{Entries: make(map[string]cacheEntry), interval: d}
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				cache.reapLoop()
			}
		}
	}()
	return &cache
}
