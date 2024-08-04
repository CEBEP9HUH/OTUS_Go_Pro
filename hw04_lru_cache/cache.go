package hw04lrucache

import (
	"sync"
)

type Key string

type elem struct {
	k Key
	v interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// mutable

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	if v, ok := lc.items[key]; ok {
		v.Value = elem{key, value}
		lc.queue.MoveToFront(v)
		return true
	}
	lc.items[key] = lc.queue.PushFront(elem{key, value})
	if lc.queue.Len() > lc.capacity {
		b := lc.queue.Back()
		lc.queue.Remove(b)
		delete(lc.items, b.Value.(elem).k)
	}
	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	if v, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(v)
		return v.Value.(elem).v, true
	}
	return nil, false
}

func (lc *lruCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	for k, v := range lc.items {
		lc.queue.Remove(v)
		delete(lc.items, k)
	}
}
