package etcd

import (
	"context"
	"log"
	"time"

	"etcd-caching-library/package/cache"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	cli *clientv3.Client
}

func NewClient(endpoints []string, dialTimeout time.Duration) (*Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) PopulateCache(cache cache.Cache, prefix string) error {
	resp, err := c.cli.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		cache.Set(string(kv.Key), string(kv.Value))
		log.Printf("Cached Key: %s | Value: %s\n", kv.Key, kv.Value)
	}
	log.Println("Cache populated from etcd")
	return nil
}

func (c *Client) Put(key, value string) error {
	_, err := c.cli.Put(context.TODO(), key, value)
	return err
}

func (c *Client) Get(key string) (string, error) {
	resp, err := c.cli.Get(context.TODO(), key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

func (c *Client) Watch(key string, handler func(string, string, string)) {
	watchChan := c.cli.Watch(context.Background(), key, clientv3.WithPrefix())
	go func() {
		for watchResp := range watchChan {
			for _, ev := range watchResp.Events {
				handler(string(ev.Type), string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}()
}
