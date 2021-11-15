package cache

import (
	"testing"
)

const (
	k = "foo"
	v = "bar"
)

func TestCache_Add(t *testing.T) {
	cache, err := New(3, Config{
		CleanInterval:             0,
		ExpirationTimeoutInterval: 0,
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	c := cache.Cap
	if c != 3 {
		t.Errorf("capacity is wrong. want %v, got %v", 3, c)
	}

	err = cache.Add(k, v, 0)
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Logf("%s-%s key-value pair added.", k, v)
	l := cache.Len
	if l != 1 {
		t.Errorf("length is wrong. want %v, got %v", 1, l)
	}
}

func TestCache_AddWithReplace(t *testing.T) {
	cache, err := New(2, Config{
		CleanInterval:             0,
		ExpirationTimeoutInterval: 0,
	})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache created.")
	c := cache.Cap
	if c != 2 {
		t.Errorf("capacity is wrong. want %v, got %v", 2, c)
	}
	t.Logf("cache capacity is true.")
	pairs := [][]string{{"key1", "val1"}, {"key2", "val2"}, {"key3", "val3"}}
	for i := 0; i < len(pairs); i++ {
		err = cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("new item added.")
		if i == 0 && cache.Len != 1 {
			t.Errorf("len must be 1, but it is %v", cache.Len)
		}
		if i == 1 && cache.Len != 2 {
			t.Errorf("len must be 2, but it is %v", cache.Len)
		}
	}
	fKey, fVal := cache.lst.Front().Value.(Item).Val.(string), cache.lst.Front().Value.(Item).Key
	sKey, sVal := cache.lst.Back().Value.(Item).Val.(string), cache.lst.Back().Value.(Item).Key
	t.Logf("%s-%s", fKey, fVal)
	t.Logf("%s-%s", sKey, sVal)
}

func TestCache_NewZeroCap(t *testing.T) {
	_, err := New(0, Config{})
	if err == nil {
		t.Errorf("expected non-nil error, but got nil error.")
	}
	t.Logf(err.Error())
}

func TestCache_NewNegativeCap(t *testing.T) {
	_, err := New(-1, Config{})
	if err == nil {
		t.Errorf("expected non-nil error, but got nil error.")
	}
	t.Logf(err.Error())
}

func TestCache_Get(t *testing.T) {
	cache, err := New(1, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache created.")

	err = cache.Add(k, v, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%s-%s added.", k, v)

	val, found := cache.Get(k)
	if !found {
		t.Errorf("expected found true, but got false")
	}
	if val != v {
		t.Errorf("expected %s, but got %s", v, val)
	}
	t.Logf("value retrieved from cache.")
	t.Logf("got: %s-%s", k, val)
}

func TestCache_GetNotFound(t *testing.T) {
	cache, err := New(1, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache created.")

	err = cache.Add(k, v, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%s-%s added", k, v)
	val, found := cache.Get("test")
	if found {
		t.Errorf("expected false, but got true")
	}
	if val != nil {
		t.Errorf("expected nil, but got %v", val)
	}
}