package main

import "sync"

type Cache struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]string),
	}
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, found := c.store[key]
	return val, found
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}
