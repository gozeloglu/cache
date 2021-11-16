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

	err = c.Add("foo", "bar", 0)
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
	_ = c.Add("key", "val", 0)

	val, found := c.Get("foo")
	if !found {
		log.Printf("data not exists in cache.")
	}
	if val != nil {
		fmt.Printf("key: foo\nvalue: %s\n", val)
	}

	err = c.Remove("foo")
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
	log.Printf("foo removed from cache.")
}
