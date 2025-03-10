package main

import (
	"context"
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// Open log file
	logFile, err := os.OpenFile("events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	// Set log output to the file
	log.SetOutput(logFile)
	log.Println("Starting etcd client...")

	// Initialize etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to connect to etcd:", err)
	}
	defer cli.Close()
	log.Println("Connected to etcd successfully.")

	// FIRST: Put and Get Method
	if _, err = cli.Put(context.TODO(), "myFirstKey", "my First Value"); err != nil {
		log.Fatal("Failed to put key:", err)
	}
	log.Println("Successfully put key: myFirstKey")

	resp, _ := cli.Get(context.TODO(), "myFirstKey")
	for _, kv := range resp.Kvs {
		log.Printf("Fetched Key: %s | Value: %s\n", kv.Key, kv.Value)
	}

	// SECOND: Watching changes of values
	watchChan := cli.Watch(context.Background(), "/foo")
	log.Println("Watching key: /foo")

	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			log.Printf("Event Type: %s | Key: %s | Value: %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
