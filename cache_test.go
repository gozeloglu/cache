package cache

import (
	"testing"
	"time"
)

const (
	k = "foo"
	v = "bar"
)

// createCache is a helper function to create cache for test functions. It is
// used for preventing code duplication.
func createCache(cap int, t *testing.T) *Cache {
	t.Helper()
	cache, err := New(cap)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("cache created.")
	return cache
}

// addItems adds the pairs to cache. It is a helper function to prevent code
// duplication.
func addItems(cache *Cache, pairs [][]string, t *testing.T) {
	t.Helper()
	for i := 0; i < len(pairs); i++ {
		err := cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}
}

func TestCache_Add(t *testing.T) {
	cache := createCache(3, t)
	c := cache.Cap()
	if c != 3 {
		t.Errorf("capacity is wrong. want %v, got %v", 3, c)
	}

	addItems(cache, [][]string{{k, v}}, t)

	if cache.lst.Front().Value.(Item).Expiration != 0 {
		t.Errorf("expiration must be 0, but it is %v", cache.lst.Front().Value.(Item).Expiration)
	}

	t.Logf("%s-%s key-value pair added.", k, v)
	l := cache.Len()
	if l != 1 {
		t.Errorf("length is wrong. want %v, got %v", 1, l)
	}
}

func TestCache_AddWithReplace(t *testing.T) {
	cache := createCache(2, t)
	c := cache.Cap()
	if c != 2 {
		t.Errorf("capacity is wrong. want %v, got %v", 2, c)
	}
	t.Logf("cache capacity is true.")
	pairs := [][]string{{"key1", "val1"}, {"key2", "val2"}, {"key3", "val3"}}
	for i := 0; i < len(pairs); i++ {
		err := cache.Add(pairs[i][0], pairs[i][1], 0)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("new item added.")
		if i == 0 && cache.Len() != 1 {
			t.Errorf("len must be 1, but it is %v", cache.Len())
		}
		if i == 1 && cache.Len() != 2 {
			t.Errorf("len must be 2, but it is %v", cache.Len())
		}
	}
	if cache.Len() != 2 {
		t.Errorf("len needs to be 2, but it is %v", cache.Len())
	}
	fKey, fVal := cache.lst.Front().Value.(Item).Val.(string), cache.lst.Front().Value.(Item).Key
	sKey, sVal := cache.lst.Back().Value.(Item).Val.(string), cache.lst.Back().Value.(Item).Key
	t.Logf("%s-%s", fKey, fVal)
	t.Logf("%s-%s", sKey, sVal)
}

func TestCache_NewZeroCap(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Errorf("expected non-nil error, but got nil error.")
	}
	t.Logf(err.Error())
}

func TestCache_NewNegativeCap(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Errorf("expected non-nil error, but got nil error.")
	}
	t.Logf(err.Error())
}

func TestCache_AddExceedCap(t *testing.T) {
	cache := createCache(1, t)

	err := cache.Add(k, v, 0)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("%s-%s added.", k, v)
	t.Logf("len: %v, cap: %v", cache.Len(), cache.Cap())

	addItems(cache, [][]string{{k + k, v + v}}, t)
	t.Logf("len: %v, cap: %v", cache.Len(), cache.Cap())

	v, found := cache.Peek(k + k)
	if !found {
		t.Errorf("%s needs to be found.", k+k)
	}
	t.Logf("%s in cache.", v)

	v, f := cache.Peek(k)
	if f {
		t.Errorf("%s should not be in the cache.", k)
	}
	if v != nil {
		t.Errorf("%v should be nil.", v)
	}
	t.Logf("%s not in cache.", k)
}

func TestCache_Get(t *testing.T) {
	cache := createCache(1, t)

	addItems(cache, [][]string{{k, v}}, t)

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
	cache := createCache(1, t)

	addItems(cache, [][]string{{k, v}}, t)

	val, found := cache.Get("test")
	if found {
		t.Errorf("expected false, but got true")
	}
	if val != nil {
		t.Errorf("expected nil, but got %v", val)
	}
}

func TestCache_GetFrontElement(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)

	val, found := cache.Get(k + k + k)
	if !found {
		t.Errorf("%s needs to be found, but didn't found.", k+k+k)
	}
	if val != v+v+v {
		t.Errorf("expected value is %s, got %s", v+v+v, val)
	}
	if cache.Len() != cache.lst.Len() {
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
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)

	val, found := cache.Get(k + k)
	if !found {
		t.Errorf("%s needs to be found, but didn't found.", k+k)
	}
	if val != v+v {
		t.Errorf("expected value is %s, got %s", v+v, val)
	}
	if cache.Len() != cache.lst.Len() {
		t.Errorf("cache length is wrong. want %v, got %v", cache.Len(), cache.lst.Len())
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
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)

	val, found := cache.Get(k)
	if !found {
		t.Errorf("%s needs to be found, but didn't found.", k)
	}
	if val != v {
		t.Errorf("expected value is %s, got %s", v, val)
	}
	if cache.Len() != cache.lst.Len() {
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
	cache := createCache(2, t)
	addItems(cache, [][]string{{k, v}}, t)

	err := cache.Remove(k)
	if err != nil {
		t.Errorf(err.Error())
	}
	if cache.Len() != 0 {
		t.Errorf("cache length should be 0, but got %v", cache.Len())
	}
	t.Logf("%s-%s pair removed successfully.", k, v)
}

func TestCache_RemoveEmptyCache(t *testing.T) {
	cache := createCache(1, t)

	err := cache.Remove(k)
	if err == nil {
		t.Errorf("error needs to be non-nil, but it is nil.")
	}
	t.Logf(err.Error())
}

func TestCache_AddRemoveGet(t *testing.T) {
	cache := createCache(1, t)
	addItems(cache, [][]string{{k, v}}, t)

	err := cache.Remove(k)
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
	cache := createCache(1, t)
	addItems(cache, [][]string{{k, v}}, t)

	found := cache.Contains(k)
	if !found {
		t.Errorf("%s needs to be found, but it is not found.", k)
	}
	t.Log(found)
	t.Logf("%s found.", k)
}

func TestCache_ContainsEmptyCache(t *testing.T) {
	cache := createCache(1, t)

	found := cache.Contains(k)
	if found {
		t.Errorf("%s needs to be not exists, but it is found.", k)
	}
	if cache.Len() != 0 {
		t.Errorf("cache needs to be empty, but it is not. len: %v", cache.Len())
	}
	t.Log(found)
	t.Logf("%s does not exists.", k)
}

func TestCache_ContainsNonExistKey(t *testing.T) {
	cache := createCache(1, t)
	addItems(cache, [][]string{{k, v}}, t)

	found := cache.Contains(k + k)
	if found {
		t.Errorf("%s needs to be not exists, but it is found.", k)
	}
	t.Logf("%s does not exist.", k+k)
}

func TestCache_ContainsCacheOrder(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)

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
	cache := createCache(3, t)
	addItems(cache, [][]string{{k, v}}, t)

	cache.Clear()
	if cache.Len() != 0 {
		t.Errorf("expected length is %v, got %v", 0, cache.Len())
	}
	if cache.lst.Len() != 0 {
		t.Errorf("expected length of list.Len() is %v, got %v", 0, cache.lst.Len())
	}

	t.Logf("cache cleared.")
}

func TestCache_ClearEmptyCache(t *testing.T) {
	cache := createCache(3, t)

	cache.Clear()
	if cache.Len() != 0 {
		t.Errorf("expected length is %v, got %v", 0, cache.Len())
	}
	if cache.lst.Len() != 0 {
		t.Errorf("expected length of list.Len() is %v, got %v", 0, cache.lst.Len())
	}

	t.Logf("empty cache cleared.")
}

func TestCache_ClearMoreThanOneDataCache(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)

	cache.Clear()
	if cache.Len() != 0 {
		t.Errorf("all data did not clear.")
	}
	if cache.lst.Front() != nil {
		t.Errorf("front node is not nil. %s-%s", cache.lst.Front().Value.(Item).Key, cache.lst.Front().Value.(Item).Val)
	}
	t.Logf("all data removed.")
}

func TestCache_Keys(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)

	keys := cache.Keys()
	if len(keys) != cache.Len() {
		t.Errorf("keys length and cache length are not the same.\nkeys length: %v, cache length: %v", len(keys), cache.Len())
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
	cache := createCache(3, t)

	keys := cache.Keys()
	if len(keys) != cache.Len() {
		t.Errorf("keys length and cache length are not the same.\nkeys length: %v, cache length: %v", len(keys), cache.Len())
	}
	t.Logf("cache is empty.")
}

func TestCache_Peek(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)

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
	cache := createCache(3, t)

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
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)

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

func TestCache_PeekNotExist(t *testing.T) {
	cache := createCache(3, t)

	val, found := cache.Peek(k)
	if found {
		t.Errorf("found should be false")
	}
	if val != nil {
		t.Errorf("val needs to be nil, but it is %v", val)
	}
	t.Logf("value is %v", val)
}

func TestCache_RemoveOldest(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)
	t.Logf("Data length in cache: %v", cache.Len())

	key, val, ok := cache.RemoveOldest()
	if !ok {
		t.Errorf("expected ok value is %v, but got %v", true, ok)
	}
	if key != k {
		t.Errorf("expected oldest key is %s, but got %s", k, key)
	}
	if val != v {
		t.Errorf("expected oldest value is %s, but got %s", v, val.(Item).Val)
	}
	if cache.Len() != 2 {
		t.Errorf("expected cache len is %v, but got %v", 2, cache.Len())
	}
	t.Logf("Oldest data removed.")
}

func TestCache_RemoveOldestEmptyCache(t *testing.T) {
	cache := createCache(3, t)

	key, val, ok := cache.RemoveOldest()
	if ok {
		t.Errorf("expected ok value is %v, but got %v", false, ok)
	}
	if key != "" {
		t.Errorf("expected key is empty string, but got %s", key)
	}
	if val != nil {
		t.Error("expected value is nil, but got ", v)
	}
	if cache.Len() != 0 {
		t.Errorf("expected cache len is %v, but got %v", 0, cache.Len())
	}
}

func TestCache_RemoveOldestCacheItemCheck(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)
	t.Logf("Data length in cache: %v", cache.Len())

	key, _, _ := cache.RemoveOldest()
	if key != k {
		t.Errorf("expected key is %s, but got %s", k, key)
	}
	order := []string{k + k + k, k + k}
	for e, i := cache.lst.Front(), 0; e != nil; e = e.Next() {
		if tmpEle := e.Value.(Item).Key; tmpEle != order[i] {
			t.Errorf("expected %s, got %s", order[i], tmpEle)
		}
		i++
	}
	t.Logf("cache order is true.")
}

func TestCache_Resize(t *testing.T) {
	cache := createCache(10, t)

	cache.len = 5
	diff := cache.Resize(8)
	if diff != 0 {
		t.Errorf("diff needs to be 0, but it is %v", diff)
	}
	if cache.Cap() != 8 {
		t.Errorf("capacity should be 8, but it is %v", cache.Cap())
	}
	t.Logf("capacity is %v", cache.Cap())
	t.Logf("diff is %v", diff)
}

func TestCache_ResizeEqualLenSize(t *testing.T) {
	cache := createCache(10, t)

	cache.len = 5
	diff := cache.Resize(5)
	if diff != 0 {
		t.Errorf("diff needs to be 0, but it is %v", diff)
	}
	if cache.Cap() != 5 {
		t.Errorf("capacity should be 5, but it is %v", cache.Cap())
	}
	t.Logf("capacity is %v", cache.Cap())
	t.Logf("diff is %v", diff)
}

func TestCache_ResizeEqualCapLenSize(t *testing.T) {
	cache := createCache(10, t)

	cache.len = 10
	diff := cache.Resize(10)
	if diff != 0 {
		t.Errorf("diff needs to be 0, but it is %v", diff)
	}
	if cache.Cap() != 10 {
		t.Errorf("capacity should be 10, but it is %v", cache.Cap())
	}
	t.Logf("capacity is %v", cache.Cap())
	t.Logf("diff is %v", diff)
}

func TestCache_ResizeExceedCap(t *testing.T) {
	cache := createCache(10, t)

	cache.len = 5
	diff := cache.Resize(12)
	if diff != 0 {
		t.Errorf("diff needs to be 0, but it is %v", diff)
	}
	if cache.Cap() != 12 {
		t.Errorf("capacity should be 8, but it is %v", cache.Cap())
	}
	t.Logf("capacity is %v", cache.Cap())
	t.Logf("diff is %v", diff)
}

func TestCache_ResizeDecreaseCap(t *testing.T) {
	cache := createCache(10, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
		{k + k + k + k, v + v + v + v},
		{k + k + k + k + k, v + v + v + v + v},
	}
	addItems(cache, pairs, t)
	t.Logf("Data length in cache: %v", cache.Len())

	diff := cache.Resize(3)
	if diff != 2 {
		t.Errorf("diff needs to be 2, but it is %v", diff)
	}
	if cache.Cap() != 3 {
		t.Errorf("new capacity needs to be 3, but it is %v", cache.Cap())
	}

	order := []string{k + k + k + k + k, k + k + k + k, k + k + k}
	for e, i := cache.lst.Front(), 0; e != nil; e = e.Next() {
		if tmpEle := e.Value.(Item).Key; tmpEle != order[i] {
			t.Errorf("expected %s, got %s", order[i], tmpEle)
		}
		i++
	}
	t.Logf("new cache order is true")
}

func TestCache_Len(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
	}
	addItems(cache, pairs, t)

	if cache.Len() != 2 {
		t.Errorf("cache length is wrong. expected %v, got %v", 2, cache.Len())
	}
	t.Logf("Data length in cache: %v", cache.Len())
}

func TestCache_Cap(t *testing.T) {
	cache := createCache(3, t)

	if cache.Cap() != 3 {
		t.Errorf("capacity should be 3, but it is %v", cache.Cap())
	}
	t.Logf("capacity is %v", cache.Cap())
}

func TestCache_Replace(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)
	t.Logf("Data length in cache: %v", cache.Len())

	err := cache.Replace(k, k+v)
	if err != nil {
		t.Errorf(err.Error())
	}
	val, found := cache.Peek(k)
	if !found {
		t.Errorf("%s does not exist.", k)
	}
	t.Logf("key (%s) value (%s) replaced with value (%s)", k, v, val)

	order := []string{k + k + k, k + k, k}
	for e, i := cache.lst.Front(), 0; e != nil; e = e.Next() {
		if ele := e.Value.(Item).Key; ele != order[i] {
			t.Errorf("expected %s, got %s", order[i], ele)
		}
		i++
	}
	t.Logf("order of the cache data is true.")
}

func TestCache_ReplaceNotExistKey(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	addItems(cache, pairs, t)
	t.Logf("Data length in cache: %v", cache.Len())

	err := cache.Replace(k+v, k+v)
	if err == nil {
		t.Errorf("it should return error because of not existing key.")
	}
	t.Logf("key did not change, because: %s", err.Error())
}

func TestCache_ClearExpiredDataEmptyCache(t *testing.T) {
	cache := createCache(3, t)
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())

	cache.ClearExpiredData()
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())
	t.Logf("No data removed.")
}

func TestCache_ClearExpiredData(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}
	for i := 0; i < len(pairs); i++ {
		err := cache.Add(pairs[i][0], pairs[i][1], -1*time.Hour)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())

	cache.ClearExpiredData()
	if cache.Len() != 0 {
		t.Errorf("all data needs to be deleted, but the length is %v", cache.Len())
	}
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())
	t.Logf("All data removed.")
}

func TestCache_ClearExpiredSomeData(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}

	var err error
	for i := 0; i < len(pairs); i++ {
		if i == 1 {
			err = cache.Add(pairs[i][0], pairs[i][1], 1*time.Hour)
		} else {
			err = cache.Add(pairs[i][0], pairs[i][1], -1*time.Hour)
		}
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())

	cache.ClearExpiredData()
	if cache.Len() != 1 {
		t.Errorf("cache len needs to be 1, but it is %v", cache.Len())
	}
	if cache.lst.Front().Value.(Item).Key != k+k {
		gotKey := cache.lst.Front().Value.(Item).Key
		gotVal := cache.lst.Front().Value.(Item).Val
		t.Errorf("front data needs to be (%s-%s) pair, but it is (%s-%s).", k+k, v+v, gotKey, gotVal)
	}
	t.Logf("Len: %v, Cap: %v", cache.Len(), cache.Cap())
	t.Logf("All data removed except one.")
}

func TestCache_ClearExpiredNoData(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}

	for i := 0; i < len(pairs); i++ {
		err := cache.Add(pairs[i][0], pairs[i][1], 1*time.Hour)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())

	cache.ClearExpiredData()
	if cache.Len() != 3 {
		t.Errorf("cache len needs to be 3, but it is %v", cache.Len())
	}
	if cache.lst.Front().Value.(Item).Key != k+k+k {
		gotKey := cache.lst.Front().Value.(Item).Key
		gotVal := cache.lst.Front().Value.(Item).Val
		t.Errorf("front data needs to be (%s-%s) pair, but it is (%s-%s).", k+k+k, v+v+v, gotKey, gotVal)
	}
	t.Logf("Len: %v, Cap: %v", cache.Len(), cache.Cap())
	t.Logf("All data removed except one.")
}

func TestCache_UpdateVal(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}

	for i := 0; i < len(pairs); i++ {
		err := cache.Add(pairs[i][0], pairs[i][1], 1*time.Hour)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())

	timeExp := cache.lst.Back().Value.(Item).Expiration
	newItem, err := cache.UpdateVal(k, k+v)
	if err != nil {
		t.Errorf(err.Error())
	}
	if newItem.Key != k {
		t.Errorf("expected key is %s, got %s", k, newItem.Key)
	}
	if newItem.Val != k+v {
		t.Errorf("expected value is %s, got %s", k+v, newItem.Val)
	}
	if newItem.Expiration != timeExp {
		t.Errorf("expected expiration time is %v, got %v", timeExp, newItem.Expiration)
	}
	t.Logf("data is updated successfully.")

	order := []string{k, k + k + k, k + k}
	i := 0
	for e := cache.lst.Front(); e != nil; e = e.Next() {
		tmpItem := e.Value.(Item)
		if tmpItem.Key != order[i] {
			t.Errorf("expected key %s, got %s", order[i], tmpItem.Key)
		}
		i++
	}
	t.Logf("cache data order is true.")
}

func TestCache_UpdateExpirationDate(t *testing.T) {
	cache := createCache(3, t)

	pairs := [][]string{
		{k, v},
		{k + k, v + v},
		{k + k + k, v + v + v},
	}

	for i := 0; i < len(pairs); i++ {
		err := cache.Add(pairs[i][0], pairs[i][1], 1*time.Hour)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())

	timeExp := cache.lst.Back().Value.(Item).Expiration
	newItem, err := cache.UpdateExpirationDate(k, 2*time.Hour)
	if err != nil {
		t.Errorf(err.Error())
	}
	if newItem.Key != k {
		t.Errorf("expected key is %s, got %s", k, newItem.Key)
	}
	if newItem.Val != v {
		t.Errorf("expected value is %s, got %s", v, newItem.Val)
	}
	if newItem.Expiration == timeExp {
		t.Errorf("expiration time needs to be updated %v, got %v", timeExp, newItem.Expiration)
	}
	t.Logf("data is updated successfully.")

	order := []string{k, k + k + k, k + k}
	i := 0
	for e := cache.lst.Front(); e != nil; e = e.Next() {
		tmpItem := e.Value.(Item)
		if tmpItem.Key != order[i] {
			t.Errorf("expected key %s, got %s", order[i], tmpItem.Key)
		}
		i++
	}
	t.Logf("cache data order is true.")
}

func TestCache_UpdateValEmptyCache(t *testing.T) {
	cache := createCache(3, t)
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())

	newItem, err := cache.UpdateVal(k, k+v)
	if err == nil {
		t.Errorf("error needs to be not nil.")
	}
	if newItem != (Item{}) {
		t.Errorf("returned item needs to be nil.")
	}
	t.Logf(err.Error())
}

func TestCache_UpdateExpirationDateEmptyCache(t *testing.T) {
	cache := createCache(3, t)
	t.Logf("Len: %v Cap: %v", cache.Len(), cache.Cap())

	newItem, err := cache.UpdateExpirationDate(k, time.Minute*5)
	if err == nil {
		t.Errorf("error needs to be not nil.")
	}
	if newItem != (Item{}) {
		t.Errorf("returned item needs to be nil.")
	}
	t.Logf(err.Error())
}

func TestItem_Expired(t *testing.T) {
	item := Item{
		Key:        k,
		Val:        v,
		Expiration: time.Now().Add(time.Minute * -1).UnixNano(),
	}

	expired := item.Expired()
	if !expired {
		t.Errorf("It needs to be expired, but it is not expired. Value is %v", expired)
	}
	t.Logf("item did not expire")
	t.Logf("expired value is %v", expired)
}

func TestItem_NotExpired(t *testing.T) {
	item := Item{
		Key:        k,
		Val:        v,
		Expiration: time.Now().Add(time.Hour * 1).UnixNano(),
	}

	expired := item.Expired()
	if expired {
		t.Errorf("It needs to not expired, but it is expired. Value is %v", expired)
	}
	t.Logf("item did not expire")
	t.Logf("expired value is %v", expired)
}

func TestItem_ExpiredNotSet(t *testing.T) {
	item := Item{
		Key:        k,
		Val:        v,
		Expiration: 0,
	}

	expired := item.Expired()
	if expired {
		t.Errorf("It needs to not expired, but it is expired. Value is %v", expired)
	}
	t.Logf("item did not expire")
	t.Logf("expired value is %v", expired)
}
