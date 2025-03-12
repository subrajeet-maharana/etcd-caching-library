package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	SetupLogRotation()
	cache := NewCache()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to connect to etcd:", err)
	}
	defer cli.Close()

	// // <----------------** ONE: Put and Get Method **---------------->
	// if _, err = cli.Put(context.TODO(), "myFirstKey", "my First Value"); err != nil {
	// 	log.Fatal(err)
	// }

	// resp, _ := cli.Get(context.TODO(), "myFirstKey")
	// for _, kv := range resp.Kvs {
	// 	fmt.Printf("Key: %s with value: %s\n", kv.Key, kv.Value)
	// 	log.Printf("Retrieved Key: %s | Value: %s", kv.Key, kv.Value)
	// }

	// // <----------------** TWO: Watching changing of values **---------------->
	// watchChan := cli.Watch(context.Background(), "/foo")
	// fmt.Println("Watch is on for the key: /foo")
	// log.Println("Watching key: /foo")

	// for watchResp := range watchChan {
	// 	for _, ev := range watchResp.Events {
	// 		log.Printf("Type: %s | Key: %s | Value: %s", ev.Type, ev.Kv.Key, ev.Kv.Value)
	// 		fmt.Printf("Type: %s | Key: %s | Value: %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	// 	}
	// }

	// // <----------------** THREE: Loop for log rotation checking **---------------->
	// for i := 0; i < 100000; i++ {
	// 	log.Printf("Log Entry #%d: This is a test log entry", i)
	// }

	// <----------------** FOUR: Cache **---------------->
	key := "myFirstKey"

	for i := 1; i <= 10; i++ {
		time.Sleep(2 * time.Second)
		if value, found := cache.Get(key); found {
			fmt.Printf("[CACHE] Key: %s | Value: %s\n", key, value)
			log.Printf("[CACHE] Key: %s | Value: %s\n", key, value)
		} else {
			resp, err := cli.Get(context.TODO(), key)
			if err != nil {
				log.Fatal("Failed to get value from etcd:", err)
			}

			for _, kv := range resp.Kvs {
				fmt.Printf("[ETCD] Key: %s | Value: %s\n", kv.Key, kv.Value)
				log.Printf("[ETCD] Key: %s | Value: %s\n", kv.Key, kv.Value)
				cache.Set(key, string(kv.Value))
			}
		}
	}

	// <----------------** FIVE: Benchmarking Cache Hits/Miss **---------------->
	// numReaders := 50
	// numWriters := 10
	// iterations := 100

	// BenchmarkCache(cache, numReaders, numWriters, iterations)
}
