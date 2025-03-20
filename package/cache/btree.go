package cache

import (
	"sync"

	"github.com/google/btree"
)

type Item struct {
	Key   string
	Value string
}

func (a Item) Less(b btree.Item) bool {
	return a.Key < b.(Item).Key
}

type BTreeCache struct {
	mu   sync.RWMutex
	tree *btree.BTree
}

func NewBTreeCache() *BTreeCache {
	return &BTreeCache{
		tree: btree.New(2),
	}
}

func (c *BTreeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tree.ReplaceOrInsert(Item{Key: key, Value: value})
}

func (c *BTreeCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item := c.tree.Get(Item{Key: key})
	if item != nil {
		return item.(Item).Value, true
	}
	return "", false
}

func (c *BTreeCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tree.Delete(Item{Key: key})
}

func (c *BTreeCache) List(start, end string) []Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var result []Item
	c.tree.AscendRange(Item{Key: start}, Item{Key: end}, func(item btree.Item) bool {
		result = append(result, item.(Item))
		return true
	})
	return result
}
