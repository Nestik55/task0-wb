package cachik

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Nestik55/task0/model"
)

type CacheItem struct {
	order  model.Order
	expiry time.Time
}

type Cache struct {
	data   map[string]CacheItem
	mu     sync.Mutex
	expiry time.Time
}

const (
	timeItem  = 120 * time.Second
	timeCache = 2 * timeItem
)

func NewCache() *Cache {
	return &Cache{
		data:   make(map[string]CacheItem),
		expiry: time.Now().Add(timeCache),
	}
}

func (ch *Cache) Set(key string, order model.Order) {
	ch.mu.Lock()
	ch.clear()

	ch.data[key] = CacheItem{
		order:  order,
		expiry: time.Now().Add(timeItem),
	}
	ch.mu.Unlock()
}

func (ch *Cache) Get(key string) (model.Order, error) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.clear()

	value, ok := ch.data[key]
	if !ok {
		return model.Order{}, errors.New("GET: Данного элемента нет в кеше")
	}

	return value.order, nil
}

func (ch *Cache) clear() {
	//fmt.Println("очистка")
	if time.Now().After(ch.expiry) {
		fmt.Println(ch.expiry, time.Now())
		//fmt.Println("началось")
		var keys []string
		for k, v := range ch.data {
			if time.Now().After(v.expiry) {
				keys = append(keys, k)
			}
		}
		for _, v := range keys {
			delete(ch.data, v)
		}

		ch.expiry = time.Now().Add(timeCache)
	}
}
