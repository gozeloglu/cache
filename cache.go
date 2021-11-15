package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Cache is the main cache type.
type Cache struct {
	// Len is the total cached data count.
	Len int

	// Cap is the maximum capacity of the cache.
	Cap int

	// CleanInterval is the time duration to make cache empty.
	CleanInterval time.Duration

	// ExpirationTimeoutInterval indicates the time to delete expired items.
	ExpirationTimeoutInterval time.Duration

	// mu is the mutex variable to prevent race conditions.
	mu sync.Mutex

	// lst is the doubly-linked list that stores the cached data.
	lst *list.List
}

// Item is the cached data type.
type Item struct {
	// Key is the value's key.
	Key string

	// Val is the value of the cached data.
	Val interface{}

	// Expiration is the amount of time to saved on memory.
	Expiration time.Duration
}

// Config keeps configuration variables.
type Config struct {
	// CleanInterval is the time duration to make cache empty.
	CleanInterval time.Duration

	// ExpirationTimeoutInterval indicates the time to delete expired items.
	ExpirationTimeoutInterval time.Duration
}

// New creates a new cache and returns it with error type. Capacity of the cache
// needs to be more than zero.
func New(cap int, config Config) (*Cache, error) {
	if cap == 0 {
		return nil, errors.New("capacity of the cache must be more than 0")
	}
	if cap < 0 {
		return nil, errors.New("capacity cannot be negative")
	}
	lst := list.New()
	return &Cache{
		Cap:                       cap,
		CleanInterval:             config.CleanInterval,
		ExpirationTimeoutInterval: config.ExpirationTimeoutInterval,
		mu:                        sync.Mutex{},
		lst:                       lst,
	}, nil
}

// Add saves data to cache if it is not saved yet.
func (c *Cache) Add(key string, val interface{}, exp time.Duration) error {
	_, found := c.get(key)
	if found {
		return errors.New("key already exists")
	}
	item := Item{
		Key:        key,
		Val:        val,
		Expiration: exp,
	}
	if c.Len == c.Cap {
		c.mu.Lock()
		lruKey := c.getLRU()
		c.delete(lruKey.Key)
		c.mu.Unlock()
	}
	c.mu.Lock()
	c.lst.PushFront(item)
	c.Len++
	c.mu.Unlock()
	return nil
}

// Get retrieves the data from list and returns it with bool information which
// indicates whether found. If there is no such data in cache, it returns nil
// and false.
func (c *Cache) Get(key string) (interface{}, bool) {
	if c.Len == 0 {
		return nil, false
	}
	c.mu.Lock()
	val, found := c.get(key)
	if val == nil {
		c.mu.Unlock()
		return nil, found
	}
	c.lst.PushFront(val)
	c.mu.Unlock()
	return val.Value.(Item).Val, found
}

// get traverses the list from head to tail and looks at the given key at each
// step. It can be considered data retrieve function for cache.
func (c *Cache) get(key string) (*list.Element, bool) {
	for e := c.lst.Front(); e != nil; e = e.Next() {
		if e.Value.(Item).Key == key {
			return e, true
		}
	}
	return nil, false
}

// delete removes the cached data from the list.
func (c *Cache) delete(key string) {
	v, found := c.get(key)
	if !found {
		return
	}
	c.lst.Remove(v)
}

// getLRU returns least recently used item from list.
func (c *Cache) getLRU() Item {
	return c.lst.Back().Value.(Item)
}
