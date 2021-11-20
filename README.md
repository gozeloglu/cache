# cache

cache is LRU-based cache package written in vanilla Go - with no package dependency. LRU stands for **Least Recently
Used** and it is one of the famous cache replacement algorithm. It replaces newly added data with the least recently
used one.

* Written in Vanilla Go, with no dependencies.
* Safe for concurrent use.
* Supports any data type for keys and values.

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
```

### Testing

You can run the tests with the following command.

```
go test .
```

### LICENSE

[MIT](https://github.com/gozeloglu/cache/blob/main/LICENSE)