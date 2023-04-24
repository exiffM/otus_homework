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
	items    map[Key]*ListItem
	keyItems map[Key]*ListItem
	keys     List
	mutex    sync.Mutex
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()
	if _, ok := cache.items[key]; ok {
		cache.items[key].Value = value
		cache.queue.MoveToFront(cache.items[key])
		cache.keys.MoveToFront(cache.keyItems[key])
		cache.mutex.Unlock()
		return true
	}
	if cache.queue.Len() < cache.capacity {
		cache.items[key] = cache.queue.PushFront(value)
		cache.keyItems[key] = cache.keys.PushFront(key)
	} else {
		cache.queue.Remove(cache.queue.Back())
		delete(cache.items, cache.keys.Back().Value.(Key))
		delete(cache.keyItems, cache.keys.Back().Value.(Key))
		cache.keys.Remove(cache.keys.Back())
		cache.items[key] = cache.queue.PushFront(value)
		cache.keyItems[key] = cache.keys.PushFront(key)
	}
	cache.mutex.Unlock()
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	if _, ok := cache.items[key]; ok {
		cache.queue.MoveToFront(cache.items[key])
		cache.keys.MoveToFront(cache.keyItems[key])
		cache.mutex.Unlock()
		return cache.items[key].Value, true
	}
	cache.mutex.Unlock()
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	cache.capacity = 0
	cache.queue = nil
	cache.keys = nil
	cache.items = nil
	cache.keyItems = nil
	cache.mutex.Unlock()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keyItems: make(map[Key]*ListItem, capacity),
		keys:     NewList(),
	}
}
