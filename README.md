# cache

cache is LRU-based cache package written in Go. LRU stands for **Least Recently Used** and it is one of the famous cache
replacement algorithm. It replaces newly added data with the least recently used one.

### Installation

You can install like this:

```
go get github.com/gozeloglu/cache
```

### Example

```go
package main

import (
	"fmt"
	"github.com/gozeloglu/cache"
	"log"
)

func main() {
	c, err := cache.New(5, cache.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = c.Add("foo", "bar", 0)
	if err != nil {
		log.Printf(err.Error())
	}
	_ = c.Add("key", "val", 0)

	val, found := c.Get("foo")
	if !found {
		log.Printf("data not exists in cache.")
	}
	if val != nil {
		fmt.Printf("key: foo\nvalue: %s", val)
	}
}
```