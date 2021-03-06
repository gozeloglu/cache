package main

import (
	"fmt"
	"github.com/gozeloglu/cache"
	"log"
)

func main() {
	// Create a cache
	c, err := cache.New(5)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Add key-value pairs to cache
	err = c.Add("foo", "bar", 0)
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
	_ = c.Add("key", "val", 0)
	_ = c.Add("fuzz", "buzz", 0)

	// Retrieve value via key
	val, found := c.Get("foo")
	if !found {
		log.Printf("data not exists in cache.")
	}
	if val != nil {
		fmt.Printf("key: foo\nvalue: %s\n", val)
	}

	// Get all keys from cache
	fmt.Println("Keys:")
	keys := c.Keys()
	for _, k := range keys {
		fmt.Println(k)
	}
	fmt.Printf("cache length: %v\n", c.Len())

	// Remove data from cache via key
	err = c.Remove("foo")
	if err != nil {
		log.Printf("%s\n", err.Error())
	}

	// Check the given key whether exists.
	found = c.Contains("key")
	if found {
		fmt.Println("key found in cache.")
	} else {
		fmt.Println("key does not exist in cache.")
	}

	found = c.Contains("foo")
	if found {
		fmt.Println("foo found in cache.")
	} else {
		fmt.Println("foo does not exist in cache.")
	}

	// Peek one of the key without updating access order.
	val, found = c.Peek("key")
	if found {
		fmt.Println("key found in cache. value is", val)
	} else {
		fmt.Println("key does not exist in cache.")
	}

	// Remove the oldest element from the cache
	k, v, ok := c.RemoveOldest()
	if ok {
		fmt.Printf("Oldest data (%s-%s) pair is removed.\n", k, v)
	} else {
		fmt.Println("Oldest data in cache did not remove.")
	}

	// Change the capacity of the cache
	c.Resize(10)
	fmt.Println("new cache capacity is", c.Cap())

	err = c.Replace("fuzz", "fuzz_buzz")
	if err != nil {
		fmt.Println(err.Error())
	}
	val, found = c.Peek("fuzz")
	if !found {
		fmt.Println("not found.")
	}
	fmt.Printf("new value is %s\n", val)
	// Clear cache. Remove everything from cache.
	c.Clear()
	fmt.Printf("cache cleared. length is %v\n", c.Len())
	_, found = c.Get("foo")
	if !found {
		fmt.Println("key does not exist.")
	}
}
