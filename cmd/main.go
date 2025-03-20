package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"etcd-caching-library/package/benchmark"
	"etcd-caching-library/package/cache"
	"etcd-caching-library/package/etcd"
	"etcd-caching-library/package/logging"
)

func main() {
	// Initialize logger
	logger := logging.NewLogger("events.log", 1, 3, 7, true)
	logger.Infof("Starting cache system")

	// Initialize caches
	btreeCache := cache.NewBTreeCache()
	lruCache := cache.NewLRUCache(100)

	// Initialize etcd client
	etcdClient, err := etcd.NewClient([]string{"localhost:2379", "localhost:22379", "localhost:32379"}, 5*time.Second)
	if err != nil {
		logger.Errorf("Failed to connect to etcd: %v", err)
		return
	}
	defer etcdClient.Close()

	// Populate caches from etcd
	if err := etcdClient.PopulateCache(btreeCache, ""); err != nil {
		logger.Errorf("Failed to populate BTree cache: %v", err)
	}

	if err := etcdClient.PopulateCache(lruCache, ""); err != nil {
		logger.Errorf("Failed to populate LRU cache: %v", err)
	}

	// Setup watch for key changes
	etcdClient.Watch("/foo", func(eventType, key, value string) {
		logger.Infof("Event: Type: %s | Key: %s | Value: %s", eventType, key, value)
		// Update cache when key changes
		lruCache.Set(key, value)
		btreeCache.Set(key, value)
	})

	// Setup example key-value
	if err := etcdClient.Put("myFirstKey", "my First Value"); err != nil {
		logger.Errorf("Failed to put key-value: %v", err)
	}

	// Benchmark caches
	go func() {
		time.Sleep(2 * time.Second) // Give time for initial setup
		logger.Infof("Starting BTree cache benchmark")
		benchmark.BenchmarkBTreeCache(btreeCache, 5, 2, 10)

		logger.Infof("Starting LRU cache benchmark")
		benchmark.BenchmarkLRUCache(lruCache, 5, 2, 10)
	}()

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Infof("Shutting down cache system")
}
