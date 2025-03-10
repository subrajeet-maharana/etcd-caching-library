package main

import (
	"sync"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	cache := NewCache()
	cache.Set("subrajeet", "maharana")
	value, found := cache.Get("subrajeet")
	if !found || value != "maharana" {
		t.Errorf("Expected 'maharana', got %s", value)
	}
}

func TestMissingKey(t *testing.T) {
	cache := NewCache()
	_, found := cache.Get("missingKey")
	if found {
		t.Errorf("Expected key to be missing, but it was found")
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewCache()
	var wg sync.WaitGroup
	iterations := 100

	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				cache.Set("key", "value")
			}
		}(i)
	}

	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_, _ = cache.Get("key")
			}
		}(i)
	}

	wg.Wait()
}
