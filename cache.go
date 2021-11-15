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

	// mu is the mutex variable to prevent race conditions.
	mu sync.Mutex

	// lst is the doubly-linked list that stores the cached data.
	lst list.List
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
		c.delete(key)
		c.mu.Unlock()
	}
	c.mu.Lock()
	c.lst.PushFront(item)
	c.Len += 1
	c.mu.Unlock()
	return nil
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
