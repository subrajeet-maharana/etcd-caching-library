package benchmark

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"etcd-caching-library/package/cache"
)

func BenchmarkBTreeCache(cache *cache.BTreeCache, numReaders, numWriters, iterations int) {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numWriters; i++ {
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

	for i := 0; i < numReaders; i++ {
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

func BenchmarkLRUCache(cache *cache.LRUCache, numReaders, numWriters, iterations int) {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numWriters; i++ {
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

	for i := 0; i < numReaders; i++ {
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
