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
	if cache.Len != 2 {
		t.Errorf("len needs to be 2, but it is %v", cache.Len)
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

func TestCache_GetFrontElement(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache created.")

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	for i := 0; i < len(pairs); i++ {
		err = cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}

	val, found := cache.Get(k + k + k)
	if !found {
		t.Errorf("%s needs to be found, but didn't found.", k+k+k)
	}
	if val != v+v+v {
		t.Errorf("expected value is %s, got %s", v+v+v, val)
	}
	if cache.Len != cache.lst.Len() {
		t.Errorf("expected value is %s, got %s", v+v+v, val)
	}
	order := []string{k + k + k, k + k, k}
	i := 0
	for e := cache.lst.Front(); e != nil; e = e.Next() {
		if e.Value.(Item).Key != order[i] {
			t.Errorf("order of the keys is wrong. expected %s, got %s", order[i], e.Value.(Item).Key)
		}
		i++
	}
	t.Logf("cache order is true")
}

func TestCache_GetMiddleElement(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache created.")

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	for i := 0; i < len(pairs); i++ {
		err = cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}

	val, found := cache.Get(k + k)
	if !found {
		t.Errorf("%s needs to be found, but didn't found.", k+k)
	}
	if val != v+v {
		t.Errorf("expected value is %s, got %s", v+v, val)
	}
	if cache.Len != cache.lst.Len() {
		t.Errorf("cache length is wrong. want %v, got %v", cache.Len, cache.lst.Len())
	}
	order := []string{k + k, k + k + k, k}
	i := 0
	for e := cache.lst.Front(); e != nil; e = e.Next() {
		if e.Value.(Item).Key != order[i] {
			t.Errorf("order of the keys is wrong. expected %s, got %s", order[i], e.Value.(Item).Key)
		}
		i++
	}
	t.Logf("cache order is true")
}

func TestCache_GetBackElement(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache created.")

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	for i := 0; i < len(pairs); i++ {
		err = cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}

	val, found := cache.Get(k)
	if !found {
		t.Errorf("%s needs to be found, but didn't found.", k)
	}
	if val != v {
		t.Errorf("expected value is %s, got %s", v, val)
	}
	if cache.Len != cache.lst.Len() {
		t.Errorf("expected value is %s, got %s", v, val)
	}
	order := []string{k, k + k + k, k + k}
	i := 0
	for e := cache.lst.Front(); e != nil; e = e.Next() {
		if e.Value.(Item).Key != order[i] {
			t.Errorf("order of the keys is wrong. expected %s, got %s", order[i], e.Value.(Item).Key)
		}
		i++
	}
	t.Logf("cache order is true")
}

func TestCache_Remove(t *testing.T) {
	cache, err := New(2, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache created.")

	err = cache.Add(k, v, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%s-%s added.", k, v)
	err = cache.Remove(k)
	if err != nil {
		t.Errorf(err.Error())
	}
	if cache.Len != 0 {
		t.Errorf("cache length should be 0, but got %v", cache.Len)
	}
	t.Logf("%s-%s pair removed successfully.", k, v)
}

func TestCache_RemoveEmptyCache(t *testing.T) {
	cache, err := New(1, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache created.")

	err = cache.Remove(k)
	if err == nil {
		t.Errorf("error needs to be non-nil, but it is nil.")
	}
	t.Logf(err.Error())
}

func TestCache_AddRemoveGet(t *testing.T) {
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

	err = cache.Remove(k)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%s-%s removed.", k, v)

	val, found := cache.Get(k)
	if found {
		t.Errorf("found needs to be false, but it is true")
	}
	if val != nil {
		t.Errorf("val needs to be nil, but it is %v", val)
	}
	t.Logf("not accessed %s-%s", k, v)
}

func TestCache_Contains(t *testing.T) {
	cache, err := New(1, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	err = cache.Add(k, v, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%s-%s added.", k, v)

	found := cache.Contains(k)
	if !found {
		t.Errorf("%s needs to be found, but it is not found.", k)
	}
	t.Log(found)
	t.Logf("%s found.", k)
}

func TestCache_ContainsEmptyCache(t *testing.T) {
	cache, err := New(1, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	found := cache.Contains(k)
	if found {
		t.Errorf("%s needs to be not exists, but it is found.", k)
	}
	if cache.Len != 0 {
		t.Errorf("cache needs to be empty, but it is not. len: %v", cache.Len)
	}
	t.Log(found)
	t.Logf("%s does not exists.", k)
}

func TestCache_ContainsNonExistKey(t *testing.T) {
	cache, err := New(1, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	err = cache.Add(k, v, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%s-%s added.", k, v)

	found := cache.Contains(k + k)
	if found {
		t.Errorf("%s needs to be not exists, but it is found.", k)
	}
	t.Logf("%s does not exist.", k+k)
}

func TestCache_ContainsCacheOrder(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	for i := 0; i < len(pairs); i++ {
		err = cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}

	order := []string{k + k + k, k + k, k}
	i := 0
	for e := cache.lst.Front(); e != nil; e = e.Next() {
		if e.Value.(Item).Key != order[i] {
			t.Errorf("order of the keys is wrong. expected %s, got %s", order[i], e.Value.(Item).Key)
		}
		i++
	}
	t.Logf("cache order is true")
}

func TestCache_Clear(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	err = cache.Add(k, v, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%s-%s added.", k, v)

	cache.Clear()
	if cache.Len != 0 {
		t.Errorf("expected length is %v, got %v", 0, cache.Len)
	}
	if cache.lst.Len() != 0 {
		t.Errorf("expected length of list.Len() is %v, got %v", 0, cache.lst.Len())
	}

	t.Logf("cache cleared.")
}

func TestCache_ClearEmptyCache(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	cache.Clear()
	if cache.Len != 0 {
		t.Errorf("expected length is %v, got %v", 0, cache.Len)
	}
	if cache.lst.Len() != 0 {
		t.Errorf("expected length of list.Len() is %v, got %v", 0, cache.lst.Len())
	}

	t.Logf("empty cache cleared.")
}

func TestCache_Keys(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	for i := 0; i < len(pairs); i++ {
		err = cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}

	keys := cache.Keys()
	if len(keys) != cache.Len {
		t.Errorf("keys length and cache length are not the same.\nkeys length: %v, cache length: %v", len(keys), cache.Len)
	}
	if len(keys) != len(pairs) {
		t.Errorf("keys length and pairs length are not the same.\nkeys length: %v, cache length: %v", len(keys), len(pairs))
	}
	for i, j := len(pairs)-1, 0; i >= 0; i-- {
		if keys[i] != pairs[j][0] {
			t.Errorf("%s and %s are not the same.", keys[i], pairs[j][0])
		}
		j++
	}
}

func TestCache_KeysEmptyCache(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	keys := cache.Keys()
	if len(keys) != cache.Len {
		t.Errorf("keys length and cache length are not the same.\nkeys length: %v, cache length: %v", len(keys), cache.Len)
	}
	t.Logf("cache is empty.")
}

func TestCache_Peek(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	for i := 0; i < len(pairs); i++ {
		err = cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}

	val, found := cache.Peek(k)
	if !found {
		t.Errorf("%s needs to be found, but couldn't found", k)
	}
	if val != v {
		t.Errorf("expected value is %s, but got %s", v, val)
	}
	t.Logf("peek works successfully.")
}

func TestCache_PeekEmptyCache(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	val, found := cache.Peek(k)
	if found {
		t.Errorf("%s needs to be not found, but it is found", k)
	}
	if val != nil {
		t.Errorf("expected value is %s, but got %s", v, val)
	}
	t.Logf("peek works successfully with empty cache.")
}

func TestCache_PeekFreqCheck(t *testing.T) {
	cache, err := New(3, Config{})
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache cretead.")

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	for i := 0; i < len(pairs); i++ {
		err = cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}

	val, found := cache.Peek(k)
	if !found {
		t.Errorf("%s needs to be found, but couldn't found", k)
	}
	if val != v {
		t.Errorf("expected value is %s, but got %s", v, val)
	}

	order := []string{k + k + k, k + k, k}
	for e, i := cache.lst.Front(), 0; e != nil; e = e.Next() {
		if tmpEle := e.Value.(Item).Key; order[i] != tmpEle {
			t.Errorf("expected %s, got %s", order[i], tmpEle)
		}
		i++
	}
	t.Logf("cache order is true after peek.")
}
