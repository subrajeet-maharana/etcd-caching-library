package cache

import (
	"container/list"
	"sync"
)

type entry struct {
	key   string
	value string
}

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	evict    *list.List
	mu       sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		evict:    list.New(),
	}
}

func (l *LRUCache) Get(key string) (string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ele, found := l.cache[key]; found {
		l.evict.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return "", false
}

func (l *LRUCache) Set(key, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ele, found := l.cache[key]; found {
		l.evict.MoveToFront(ele)
		ele.Value.(*entry).value = value
		return
	}

	if len(l.cache) >= l.capacity {
		l.evictOldest()
	}

	ent := &entry{key, value}
	element := l.evict.PushFront(ent)
	l.cache[key] = element
}

func (l *LRUCache) Delete(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ele, found := l.cache[key]; found {
		l.evict.Remove(ele)
		delete(l.cache, key)
	}
}

func (l *LRUCache) evictOldest() {
	oldest := l.evict.Back()
	if oldest != nil {
		l.evict.Remove(oldest)
		delete(l.cache, oldest.Value.(*entry).key)
	}
}
