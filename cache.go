package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/btree"
)

type Item struct {
	Key   string
	Value string
}

type Cache struct {
	mu   sync.RWMutex
	tree *btree.BTree
}

func (a Item) Less(b btree.Item) bool {
	return a.Key < b.(Item).Key
}

func NewCache() *Cache {
	return &Cache{
		tree: btree.New(2),
	}
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tree.ReplaceOrInsert(Item{Key: key, Value: value})
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item := c.tree.Get(Item{Key: key})
	if item != nil {
		return item.(Item).Value, true
	}
	return "", false
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tree.Delete(Item{Key: key})
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
