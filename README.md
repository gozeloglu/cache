# cache [![Go Reference](https://pkg.go.dev/badge/github.com/gozeloglu/cache.svg)](https://pkg.go.dev/github.com/gozeloglu/cache) [![Go Report Card](https://goreportcard.com/badge/github.com/gozeloglu/cache)](https://goreportcard.com/report/github.com/gozeloglu/cache) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gozeloglu/cache) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/gozeloglu/cache) [![LICENSE](https://img.shields.io/badge/license-MIT-green)](https://github.com/gozeloglu/cache/blob/main/LICENSE)

cache is LRU-based cache package written in vanilla Go - with no package dependency. LRU stands for **Least Recently
Used** and it is one of the famous cache replacement algorithm. It replaces newly added data with the least recently
used one.

* Written in Vanilla Go, with no dependencies.
* Safe for concurrent use.
* Supports any data type for keys and values.
* Supports time expiration.

### Installation

```
go get github.com/gozeloglu/cache
```

### Example

Here, there is an example usage of the package.

You can import like this.

```go
import "github.com/gozeloglu/cache"
```

Then, you need to create a cache variable with `New()` function. It takes one parameter to specify cache capacity.

#### Add new data

```go
cache.Add("foo", "bar", 0) // Without expiration time
cache.Add("key", "value", time.Hour * 2) // With expiration time
```

#### Get data

```go
val, found := cache.Get("foo")
if !found {
    fmt.Println("key does not exist. val is nil.")
}
fmt.Println(val)
```

#### Get all keys

```go
keys := cache.Keys()
for _, k := range keys {
    fmt.Println(k)
}
```

#### Contains, Peek and Remove

```go
found := cache.Contains("foo")
if found {
    val, _ := cache.Peek("foo")
    cache.Remove("foo")
}
```

#### Remove Oldest

```go
k, v, ok := cache.RemoveOldest()
if ok {
    fmt.Printf("Oldest key-value pair removed: %s-%s", k, v)
}
```

#### Resize

```go
cache.Resize(20) // Capacity will be 20
```

#### Update value, update expiration date, and replace

```go
newItem, err := cache.UpdateVal("foo", "foobar") // Cache data order is also updated
if err != nil {
    fmt.Printf("New item key and value is %s-%s", newItem.Key, newItem.Val)
}
newItem, err := c.UpdateExpirationDate("foo", time.Hour * 4) // Cache data order is also updated
if err != nil {
    fmt.Printf("%v", newItem.Expiration)
}

err = c.Replace("foo", "fuzz")  // Change value of the key without updating cache access order
if err != nil {
	fmt.Printf(err.Error())
}
```

### Testing

You can run the tests with the following command.

```
go test .
```

### Code Coverage

You can get the code coverage information with the following command:

```bash
go test -cover
```

If you want to generate a graphical coverage report, you can run the following command:

```bash
go tool cover -html=coverage.out
```

A browser tab will be opened and you will be able to see the graphical report. It shows not tracked, not covered, and covered line on the source code. 

### LICENSE

[MIT](https://github.com/gozeloglu/cache/blob/main/LICENSE)