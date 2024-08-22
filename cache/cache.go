package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	order  interface{}
	expiry time.Time
}

type Cache struct {
	data      map[string]CacheItem
	mu        sync.Mutex
	timeStart time.Time
}

const (
	timeItem  = 1 * time.Minute
	timeCache = 2 * timeItem
)

func NewCache() *Cache {
	return &Cache{
		data:      make(map[string]CacheItem),
		timeStart: time.Now(),
	}
}

func (ch *Cache) Set(key string, order interface{}) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	defer ch.clear()

	ch.data[key] = CacheItem{
		order:  order,
		expiry: time.Now().Add(timeItem),
	}
}

func (ch *Cache) Get(key string) interface{} {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	//defer ch.clear()

	value, ok := ch.data[key]
	if !ok {
		return nil
	}

	return value.order
}

func (ch *Cache) clear() {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if ch.timeStart.Sub(time.Now()) > timeCache {
		for k, v := range ch.data {
			if v.expiry.After(time.Now()) {
				delete(ch.data, k)
			}
		}

		ch.timeStart = time.Now()
	}
}
