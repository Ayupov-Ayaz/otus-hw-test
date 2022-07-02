package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	mu       *sync.RWMutex
	items    map[Key]*ListItem
}

func newLruCache(capacity int, queue List) *lruCache {
	return &lruCache{
		capacity: capacity,
		queue:    queue,
		items:    make(map[Key]*ListItem, capacity),
		mu:       &sync.RWMutex{},
	}
}

func (l *lruCache) removeLastElement() {
	lastItem := l.queue.Back()
	l.queue.Remove(lastItem)

	l.mu.Lock()
	for k, v := range l.items {
		if v == lastItem {
			delete(l.items, k)
			break
		}
	}
	l.mu.Unlock()
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	var (
		item *ListItem
		ok   bool
	)

	l.mu.RLock()
	item, ok = l.items[key]
	l.mu.RUnlock()

	if !ok && l.queue.Len() >= l.capacity {
		l.removeLastElement()
	}

	if !ok {
		item = l.queue.PushFront(value)
	} else {
		item.Value = value
		l.queue.MoveToFront(item)
	}

	l.mu.Lock()
	l.items[key] = item
	l.mu.Unlock()

	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.RLock()
	v, ok := l.items[key]
	l.mu.RUnlock()

	if !ok {
		return nil, false
	}

	l.queue.MoveToFront(v)

	return v.Value, true
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.mu.Unlock()
}

type cacheItem struct { // зачем это?
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return newLruCache(capacity, NewList())
}
