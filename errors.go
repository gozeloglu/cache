package cache

import "errors"

var (
	errEmptyCache   = errors.New("cache is empty")
	errNegCapacity  = errors.New("capacity cannot be negative")
	errZeroCapacity = errors.New("cache capacity should be more than zero")
	errKeyExist     = errors.New("key already exists")
	errKeyNotExist  = errors.New("key does not exist")
	errNoKey        = errors.New("there is no such key")
)
