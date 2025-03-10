package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

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

func BenchmarkCache(cache *Cache, numReaders, numWriters, iterations int) {
	var wg sync.WaitGroup
	start := time.Now()

	for i := range numWriters {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := fmt.Sprintf("Key-%d", rand.Intn(10))
				value := fmt.Sprintf("Value-%d", rand.Intn(100))
				cache.Set(key, value)
				log.Printf("[Writer-%d] SET %s -> %s\n", id, key, value)
				time.Sleep(time.Microsecond * 10)
			}
		}(i)
	}

	for i := range numReaders {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				key := fmt.Sprintf("Key-%d", rand.Intn(10))
				if value, found := cache.Get(key); found {
					log.Printf("[Reader-%d] HIT %s -> %s\n", id, key, value)
				} else {
					log.Printf("[Reader-%d] MISS %s\n", id, key)
				}
				time.Sleep(time.Microsecond * 10)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Benchmark completed in %v\n", time.Since(start))
}
