package main

import (
	"fmt"
	"github.com/gozeloglu/cache"
	"log"
)

func main() {
	// Create a cache
	c, err := cache.New(5, cache.Config{})
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
	fmt.Printf("cache length: %v\n", c.Len)

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

	// Clear cache. Remove everything from cache.
	c.Clear()
	fmt.Printf("cache cleared. length is %v\n", c.Len)
	val, found = c.Get("foo")
	if !found {
		fmt.Println("key does not exist.")
	}
}
