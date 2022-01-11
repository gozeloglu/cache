package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Cache is the main cache type.
type Cache struct {
	// len is the total cached data count.
	len int

	// cap is the maximum capacity of the cache.
	cap int

	// mu is the mutex variable to prevent race conditions.
	mu sync.Mutex

	// lst is the doubly-linked list that stores the cached data.
	lst *list.List
}

// Item is the cached data type.
type Item struct {
	// Key is the value's key.
	Key interface{}

	// Val is the value of the cached data.
	Val interface{}

	// Expiration is the amount of time to saved on memory.
	Expiration int64
}

// New creates a new cache and returns it with error type. Capacity of the cache
// needs to be more than zero.
func New(cap int) (*Cache, error) {
	if cap == 0 {
		return nil, errors.New("capacity of the cache must be more than 0")
	}
	if cap < 0 {
		return nil, errors.New("capacity cannot be negative")
	}
	lst := list.New()
	return &Cache{
		cap: cap,
		mu:  sync.Mutex{},
		lst: lst,
	}, nil
}

// Add saves data to cache if it is not saved yet. If the capacity is full,
// the least-recently used one will be removed and new data will be added.
// If you do not want to add an expired time for data, you need to pass 0.
func (c *Cache) Add(key interface{}, val interface{}, exp time.Duration) error {
	_, found := c.get(key)
	if found {
		return errors.New("key already exists")
	}
	item := Item{
		Key:        key,
		Val:        val,
		Expiration: time.Now().Add(exp).UnixNano(),
	}
	if exp == 0 {
		item.Expiration = 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Len() == c.Cap() {
		lruKey := c.getLRU()
		c.delete(lruKey.Key)
	}

	c.lst.PushFront(item)
	c.len++
	return nil
}

// Get retrieves the data from list and returns it with bool information which
// indicates whether found. If there is no such data in cache, it returns nil
// and false.
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	if c.Len() == 0 {
		return nil, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	val, found := c.get(key)
	if val == nil {
		return nil, found
	}
	e := c.lst.Remove(val)
	c.lst.PushFront(e)
	return val.Value.(Item).Val, found
}

// Remove deletes the item from the cache. Updates the length of the cache
// decrementing by one.
func (c *Cache) Remove(key interface{}) error {
	if c.Len() == 0 {
		return errors.New("empty cache")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.delete(key)
	return nil
}

// Contains checks the given key and returns the information that it exists
// on cache or not. Calling this function doesn't change the access order of
// the cache.
func (c *Cache) Contains(key interface{}) bool {
	if c.Len() == 0 {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.get(key)
	return found
}

// Clear deletes all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clear()
}

// Keys returns all keys in cache. It does not change frequency of the item
// access.
func (c *Cache) Keys() []interface{} {
	var keys []interface{}

	c.mu.Lock()
	defer c.mu.Unlock()
	for e := c.lst.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(Item).Key)
	}

	return keys
}

// Peek returns the given key without updating access frequency of the item.
func (c *Cache) Peek(key interface{}) (interface{}, bool) {
	if c.Len() == 0 {
		return nil, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	val, found := c.get(key)
	if !found {
		return nil, found
	}
	return val.Value.(Item).Val, found
}

// RemoveOldest removes the least recently used one. Returns removed key, value,
// and bool value that indicates whether remove operation is done successfully.
func (c *Cache) RemoveOldest() (k interface{}, v interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	k, v, ok = c.removeOldest()
	return
}

// Resize changes the size of the capacity. If new capacity is lower than
// existing capacity, the oldest items will be removed. It returns the number
// of the removed oldest elements from the cache. If it is zero, means that
// no data removed from the cache.
func (c *Cache) Resize(size int) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	diff := c.resize(size)
	return diff
}

// Len returns length of the cache.
func (c *Cache) Len() int {
	return c.len
}

// Cap returns capacity of the cache.
func (c *Cache) Cap() int {
	return c.cap
}

// Replace changes the value of the given key, if the key exists. If the key
// does not exist, it returns error. Calling Replace function does not change
// the cache order.
func (c *Cache) Replace(key interface{}, val interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, found := c.get(key)
	if !found {
		return errors.New("key does not exist")
	}
	e.Value = Item{
		Key:        key,
		Val:        val,
		Expiration: e.Value.(Item).Expiration,
	}
	return nil
}

// ClearExpiredData deletes the all expired data in cache.
func (c *Cache) ClearExpiredData() {
	c.mu.Lock()
	defer c.mu.Unlock()
	l := c.Len()
	if l == 0 {
		return
	}

	now := time.Now().UnixNano()
	c.clearExpiredData(now)
}

// UpdateVal updates the value of the given key. If there is no such a data, error
// will be returned. Cache data order is updated after updating the value. It
// returns updated item.
func (c *Cache) UpdateVal(key interface{}, val interface{}) (Item, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.update(key, val, -1)
}

// UpdateExpirationDate updates the expiration date of the given key. If there
// is no such a data, error will be returned. Cache data order is updated after
// updating the expiration time. It returns updated item.
func (c *Cache) UpdateExpirationDate(key interface{}, exp time.Duration) (Item, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	newExpTime := time.Now().Add(exp).Unix()
	return c.update(key, nil, newExpTime)
}

// Expired returns true if the item expired.
func (i Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return i.Expiration < time.Now().UnixNano()
}

// get traverses the list from head to tail and looks at the given key at each
// step. It can be considered data retrieve function for cache.
func (c *Cache) get(key interface{}) (*list.Element, bool) {
	for e := c.lst.Front(); e != nil; e = e.Next() {
		if e.Value.(Item).Key == key {
			return e, true
		}
	}
	return nil, false
}

// delete removes the cached data from the list.
func (c *Cache) delete(key interface{}) {
	v, found := c.get(key)
	if !found {
		return
	}
	c.lst.Remove(v)
	c.len--
}

// getLRU returns least recently used item from list.
func (c *Cache) getLRU() Item {
	return c.lst.Back().Value.(Item)
}

// clear removes all elements from the list.
func (c *Cache) clear() {
	var next *list.Element
	for e := c.lst.Front(); e != nil; e = next {
		next = e.Next()
		c.lst.Remove(e)
		c.len--
	}
}

// removeOldest removes the oldest data from the cache.
func (c *Cache) removeOldest() (key interface{}, val interface{}, ok bool) {
	if c.Len() == 0 {
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
	if size < c.Len() {
		diff = c.Len() - size
	}

	for i := 0; i < diff; i++ {
		c.removeOldest()
	}
	c.cap = size

	return diff
}

// clearExpiredData removes the all expired data in cache.
func (c *Cache) clearExpiredData(now int64) {
	var next *list.Element
	for e := c.lst.Front(); e != nil; e = next {
		next = e.Next()
		if exp := e.Value.(Item).Expiration; exp != 0 && exp < now {
			c.lst.Remove(e)
			c.len--
		}
	}
}

// update changes the val and/or expiration date.
func (c *Cache) update(key interface{}, val interface{}, exp int64) (Item, error) {
	var next *list.Element
	for e := c.lst.Front(); e != nil; e = next {
		next = e.Next()
		if k := e.Value.(Item).Key; k == key {
			if val == nil {
				val = e.Value.(Item).Val
			}
			if exp == -1 {
				exp = e.Value.(Item).Expiration
			}

			c.lst.Remove(e)
			newItem := Item{
				Key:        key,
				Val:        val,
				Expiration: exp,
			}
			c.lst.PushFront(newItem)
			return newItem, nil
		}
	}
	return Item{}, errors.New("there is no such key")
}
