/*
Package cache is a cache package written in Go with no dependency. It uses
Least-Recently Used (LRU) algorithm for replacement policy. Behind the package,
it uses doubly-linked list which is built-in data structure that comes from
container/list. Any data can be stored in cache.

Firstly, you need to create a new cache as follows.

	c, err := cache.New(5, cache.Config{})

The first parameter is capacity of the cache. For the example, the cache can keep
up to 5 data. If a new data is being tried to add when the cache is full, the
cache removes the least-recently used one and adds the new data.
*/
package cache
