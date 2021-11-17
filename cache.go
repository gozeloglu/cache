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
	e := c.lst.Remove(val)
	c.lst.PushFront(e)
	c.mu.Unlock()
	return val.Value.(Item).Val, found
}

// Remove deletes the item from the cache. Updates the length of the cache
// decrementing by one.
func (c *Cache) Remove(key string) error {
	if c.Len == 0 {
		return errors.New("empty cache")
	}

	c.mu.Lock()
	c.delete(key)
	c.mu.Unlock()
	return nil
}

// Contains checks the given key and returns the information that it exists
// on cache or not. Calling this function doesn't change the access order of
// the cache.
func (c *Cache) Contains(key string) bool {
	if c.Len == 0 {
		return false
	}
	c.mu.Lock()
	_, found := c.get(key)
	if found {
		c.mu.Unlock()
		return true
	}
	c.mu.Unlock()
	return false
}

// Clear deletes all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	c.clear()
	c.mu.Unlock()
}

// Keys returns all keys in cache. It does not change frequency of the item
// access.
func (c *Cache) Keys() []string {
	var keys []string

	c.mu.Lock()
	for e := c.lst.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(Item).Key)
	}
	c.mu.Unlock()

	return keys
}

// Peek returns the given key without updating access frequency of the item.
func (c *Cache) Peek(key string) (interface{}, bool) {
	if c.Len == 0 {
		return nil, false
	}
	c.mu.Lock()
	val, found := c.get(key)
	c.mu.Unlock()
	return val.Value.(Item).Val, found
}

// RemoveOldest removes the least recently used one. Returns removed key, value,
// and bool value that indicates whether remove operation is done successfully.
func (c *Cache) RemoveOldest() (k string, v interface{}, ok bool) {
	c.mu.Lock()
	k, v, ok = c.removeOldest()
	c.mu.Unlock()
	return
}

// Resize changes the size of the capacity. If new capacity is lower than
// existing capacity, the oldest items will be removed. It returns the number
// of the removed oldest elements from the cache. If it is zero, means that
// no data removed from the cache.
func (c *Cache) Resize(size int) int {
	c.mu.Lock()
	diff := c.resize(size)
	c.mu.Unlock()
	return diff
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
	c.Len--
}

// getLRU returns least recently used item from list.
func (c *Cache) getLRU() Item {
	return c.lst.Back().Value.(Item)
}

// clear removes all elements from the list.
func (c *Cache) clear() {
	for e := c.lst.Front(); e != nil; e = e.Next() {
		c.lst.Remove(e)
		c.Len--
	}
}

// removeOldest removes the oldest data from the cache.
func (c *Cache) removeOldest() (key string, val interface{}, ok bool) {
	if c.Len == 0 {
		return "", nil, false
	}
	oldest := c.getLRU()
	key, val = oldest.Key, oldest.Val
	c.delete(key)
	ok = true
	return
}

// resize changes the capacity of the cache. It prunes the oldest elements from
// the cache if the size is lower than length of the cache.
func (c *Cache) resize(size int) int {
	var diff int
	if size < c.Len {
		diff = c.Len - size
	}

	for i := 0; i < diff; i++ {
		c.removeOldest()
	}
	c.Cap = size

	return diff
}
