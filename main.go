package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// FIRST: Put and Get Method
	if _, err = cli.Put(context.TODO(), "myFirstKey", "my First Value"); err != nil {
		log.Fatal(err)
	}

	resp, _ := cli.Get(context.TODO(), "myFirstKey")
	for _, kv := range resp.Kvs {
		fmt.Printf("Key: %s with value: %s", kv.Key, kv.Value)
	}

	// SECOND: Watching changing of values
	watchChan := cli.Watch(context.Background(), "/foo")
	fmt.Printf("Watch is on for the key: /foo")
	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			fmt.Printf("Type: %s | Key: %s | Value: %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

}
