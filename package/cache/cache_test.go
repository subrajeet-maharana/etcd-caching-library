package cache_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"etcd-caching-library/package/cache"
)

func TestBTreeCache(t *testing.T) {
	t.Run("SetAndGet", func(t *testing.T) {
		cache := cache.NewBTreeCache()
		cache.Set("subrajeet", "maharana")
		value, found := cache.Get("subrajeet")
		if !found || value != "maharana" {
			t.Errorf("Expected 'maharana', got %s", value)
		}
	})

	t.Run("MissingKey", func(t *testing.T) {
		cache := cache.NewBTreeCache()
		_, found := cache.Get("missingKey")
		if found {
			t.Errorf("Expected key to be missing, but it was found")
		}
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		cache := cache.NewBTreeCache()
		var wg sync.WaitGroup
		iterations := 100

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					cache.Set("key", "value")
				}
			}(i)
		}

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					_, _ = cache.Get("key")
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("List", func(t *testing.T) {
		cache := cache.NewBTreeCache()

		for i := 0; i < 1000; i++ {
			cache.Set(fmt.Sprintf("Key-%d", i), fmt.Sprintf("Value-%d", i))
		}

		start := time.Now()
		items := cache.List("Key-100", "Key-900")
		duration := time.Since(start)

		if len(items) == 0 {
			t.Error("Expected items in list, got none")
		}

		t.Logf("BTree List Query took %v and returned %d items", duration, len(items))
	})
}

func TestLRUCache(t *testing.T) {
	t.Run("SetAndGet", func(t *testing.T) {
		cache := cache.NewLRUCache(10)
		cache.Set("subrajeet", "maharana")
		value, found := cache.Get("subrajeet")
		if !found || value != "maharana" {
			t.Errorf("Expected 'maharana', got %s", value)
		}
	})

	t.Run("MissingKey", func(t *testing.T) {
		cache := cache.NewLRUCache(10)
		_, found := cache.Get("missingKey")
		if found {
			t.Errorf("Expected key to be missing, but it was found")
		}
	})

	t.Run("Eviction", func(t *testing.T) {
		cache := cache.NewLRUCache(3)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")

		// This will evict key1
		cache.Set("key4", "value4")

		_, found := cache.Get("key1")
		if found {
			t.Errorf("Expected key1 to be evicted, but it was found")
		}

		// Verify key2, key3, key4 are still there
		_, found = cache.Get("key2")
		if !found {
			t.Errorf("Expected key2 to be present, but it was not found")
		}

		_, found = cache.Get("key3")
		if !found {
			t.Errorf("Expected key3 to be present, but it was not found")
		}

		_, found = cache.Get("key4")
		if !found {
			t.Errorf("Expected key4 to be present, but it was not found")
		}
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		cache := cache.NewLRUCache(100)
		var wg sync.WaitGroup
		iterations := 100

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					cache.Set("key", "value")
				}
			}(i)
		}

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					_, _ = cache.Get("key")
				}
			}(i)
		}

		wg.Wait()
	})
}
