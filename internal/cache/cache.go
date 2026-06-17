package cache

import (
	"strings"
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]*entry
	ttl   time.Duration
}

func New(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]*entry),
		ttl:   ttl,
	}
	go c.cleanupLoop()
	return c
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for k, e := range c.items {
			if time.Now().After(e.expiresAt) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	c.items[key] = &entry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}

func (c *Cache) Shutdown() {
	c.mu.Lock()
	c.items = make(map[string]*entry)
	c.mu.Unlock()
}

func (c *Cache) Del(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

func (c *Cache) Clear() {
	c.mu.Lock()
	c.items = make(map[string]*entry)
	c.mu.Unlock()
}

func (c *Cache) DeleteByPrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.items {
		if strings.HasPrefix(k, prefix) {
			delete(c.items, k)
		}
	}
}
