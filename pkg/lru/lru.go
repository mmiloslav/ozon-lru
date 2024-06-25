// пакет работы с LRUCache
package lru

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

// errors
const (
	ErrKeyNotFound = "key not found"
	ErrKeyExpired  = "key expired"
)

// интерфейс ILRUCache
type ILRUCache interface {
	// Put запись данных в кэш
	Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	// Get получение данных из кэша по ключу
	Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error)
	// GetAll получение всего наполнения кэша в виде двух слайсов: слайса ключей и слайса значений. Пары ключ-значения из кэша располагаются на соответствующих позициях в слайсах.
	GetAll(ctx context.Context) (keys []string, values []interface{}, err error)
	// Evict ручное удаление данных по ключу
	Evict(ctx context.Context, key string) (value interface{}, err error)
	// EvictAll ручная инвалидация всего кэша
	EvictAll(ctx context.Context) error
}

// ключ-значение в кеше со временем истечения срока действия
type Pair struct {
	key       string
	value     interface{}
	expiresAt time.Time
}

// структура LRU-кеша
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	mu       sync.Mutex
}

// создание нового LRU-кеша
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// добавление значения в кеш по ключу
func (lru *LRUCache) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	expiresAt := time.Now().Add(ttl)
	if element, ok := lru.cache[key]; ok {
		lru.list.MoveToFront(element)
		element.Value = Pair{key, value, expiresAt}
		return nil
	}
	if lru.list.Len() == lru.capacity {
		back := lru.list.Back()
		if back != nil {
			lru.list.Remove(back)
			delete(lru.cache, back.Value.(Pair).key)
		}
	}
	pair := Pair{key, value, expiresAt}
	element := lru.list.PushFront(pair)
	lru.cache[key] = element
	return nil
}

// получение значения по ключу из кеша
func (lru *LRUCache) Get(ctx context.Context, key string) (interface{}, time.Time, error) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element, ok := lru.cache[key]; ok {
		if time.Now().After(element.Value.(Pair).expiresAt) {
			lru.list.Remove(element)
			delete(lru.cache, key)
			return nil, time.Time{}, errors.New(ErrKeyExpired)
		}
		lru.list.MoveToFront(element)
		return element.Value.(Pair).value, element.Value.(Pair).expiresAt, nil
	}
	return nil, time.Time{}, errors.New(ErrKeyNotFound)
}

// получение всего наполнения кэша в виде двух слайсов: слайса ключей и слайса значений.
func (lru *LRUCache) GetAll(ctx context.Context) ([]string, []interface{}, error) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	keys := make([]string, 0, lru.list.Len())
	values := make([]interface{}, 0, lru.list.Len())

	for e := lru.list.Front(); e != nil; e = e.Next() {
		if time.Now().Before(e.Value.(Pair).expiresAt) {
			keys = append(keys, e.Value.(Pair).key)
			values = append(values, e.Value.(Pair).value)
		} else {
			lru.list.Remove(e)
			delete(lru.cache, e.Value.(Pair).key)
		}
	}

	return keys, values, nil
}

// удаление элемента по ключу
func (lru *LRUCache) Evict(ctx context.Context, key string) (interface{}, error) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element, ok := lru.cache[key]; ok {
		lru.list.Remove(element)
		delete(lru.cache, key)
		return element.Value.(Pair).value, nil
	}
	return nil, errors.New(ErrKeyNotFound)
}

// очищение всего кеша
func (lru *LRUCache) EvictAll(ctx context.Context) error {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.list.Init()
	lru.cache = make(map[string]*list.Element)
	return nil
}
