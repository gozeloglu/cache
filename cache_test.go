package cache

import (
	"errors"
	"reflect"
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

// addItemsWithExp adds the pairs with expiration duration to cache. It is a helper
// function to prevent code duplication.
func addItemsWithExp(t *testing.T, cache *Cache, pairs [][]any) {
	t.Helper()
	for i := 0; i < len(pairs); i++ {
		exp := pairs[i][2].(time.Duration)
		err := cache.Add(pairs[i][0], pairs[i][1], exp)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Logf("%s-%s added.", pairs[i][0], pairs[i][1])
	}
}

// findItem finds the item corresponding to the given key in cache. It is a helper
// function to prevent code duplication.
func findItem(t *testing.T, c *Cache, key any) (Item, bool) {
	t.Helper()
	for e := c.lst.Front(); e != nil; e = e.Next() {
		if item := e.Value.(Item); item.Key == key {
			return item, true
		}
	}
	return Item{}, false
}

// cmpCacheListOrder compares the order of the items (reflecting the access frequency
// of items) in cache. It is a helper function to prevent code duplication.
func cmpCacheListOrder(t *testing.T, c *Cache, order []any) {
	t.Helper()
	if c.lst.Len() != len(order) {
		t.Errorf("expect cache list to have length, %v, want %v", c.lst.Len(), len(order))
	}

	i := 0
	e := c.lst.Front()
	for e != nil {
		o := order[i]
		k := e.Value.(Item).Key
		if !reflect.DeepEqual(k, o) {
			t.Errorf("incorrect key order, got %v, want %v at index %d", k, o, i)
		}
		e = e.Next()
		i++
	}
}

func TestCache_Add(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		wantLength        int
		wantKeysListOrder []any
	}{
		{
			name:              "successfully adds a single item to an empty cache",
			capacity:          3,
			addPairs:          [][]any{{k, v}},
			wantLength:        1,
			wantKeysListOrder: nil,
		},
		{
			name:              "successfully adds multiple items to an empty cache",
			capacity:          5,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			wantLength:        3,
			wantKeysListOrder: nil,
		},
		{
			name:              "replaces LRU item when added item count exceeds capacity",
			capacity:          2,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			wantLength:        2,
			wantKeysListOrder: []any{k + k + k, k + k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		t.Run(tt.name, func(t *testing.T) {
			for _, pair := range tt.addPairs {
				var (
					wantErr error = nil
					wantExp int64 = 0
				)
				err := c.Add(pair[0], pair[1], time.Duration(wantExp))
				if !errors.Is(err, wantErr) {
					t.Errorf("unexpected error, got error %v, want %v", err, wantErr)
					return
				}
				if exp := c.lst.Front().Value.(Item).Expiration; exp != wantExp {
					t.Errorf("unexpected expiration, got %v want %v", exp, wantExp)
				}
			}
			if c.Len() != tt.wantLength {
				t.Errorf("unexpected length, got %v, want %v", c.Len(), tt.wantLength)
			}
			if tt.wantKeysListOrder != nil {
				cmpCacheListOrder(t, c, tt.wantKeysListOrder)
			}
		})
	}
}

func TestCache_New(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		want     *Cache
		wantErr  error
	}{
		{
			name:     "returns error when provided capacity == 0",
			capacity: 0,
			want:     nil,
			wantErr:  errZeroCapacity,
		},
		{
			name:     "returns error when provided capacity < 0",
			capacity: -1,
			want:     nil,
			wantErr:  errNegCapacity,
		},
		{
			name:     "creates cache with given capacity, when capacity > 0",
			capacity: 20,
			want:     createCache(20, t),
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.capacity)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("cache.New() error = %v, want %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Get(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		getPairs          [][]any
		wantKeysListOrder []any
	}{
		{
			name:              "expect pair to be found when its added to cache",
			capacity:          1,
			addPairs:          [][]any{{k, v}},
			getPairs:          [][]any{{k, v}},
			wantKeysListOrder: nil,
		},
		{
			name:              "expect pair to not be found when its not added to cache",
			capacity:          1,
			addPairs:          [][]any{{k, v}},
			getPairs:          [][]any{{"nonexistent", nil}},
			wantKeysListOrder: nil,
		},
		{
			name:              "cache list has correct order after getting front element",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			getPairs:          [][]any{{k + k + k, v + v + v}},
			wantKeysListOrder: []any{k + k + k, k + k, k},
		},
		{
			name:              "cache list has correct order after getting middle element",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			getPairs:          [][]any{{k + k, v + v}},
			wantKeysListOrder: []any{k + k, k + k + k, k},
		},
		{
			name:              "cache list has correct order after getting back element",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			getPairs:          [][]any{{k, v}},
			wantKeysListOrder: []any{k, k + k + k, k + k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			for _, pair := range tt.getPairs {
				var (
					want          = pair[1]
					wantFound     = want != nil
					got, gotFound = c.Get(pair[0])
				)
				if gotFound != wantFound {
					t.Errorf("cache.Get() found = %v, want %v", gotFound, wantFound)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("cache.Get() = %v, want %v", got, want)
				}
			}
			if c.Len() != c.lst.Len() {
				t.Errorf("incorrect cache length, want %v, got %v", c.Len(), c.lst.Len())
			}
			if tt.wantKeysListOrder != nil {
				cmpCacheListOrder(t, c, tt.wantKeysListOrder)
			}
		})
	}
}

func TestCache_Remove(t *testing.T) {
	tests := []struct {
		name           string
		capacity       int
		addPairs       [][]any
		removeKeys     []any
		wantErrs       []error
		checkLength    bool
		getRemovedKeys bool
	}{
		{
			name:           "returns an error, when removing key from empty cache",
			capacity:       1,
			addPairs:       [][]any{},
			removeKeys:     []any{k},
			wantErrs:       []error{errEmptyCache},
			checkLength:    false,
			getRemovedKeys: false,
		},
		{
			name:           "successfully removes keys of items added to cache",
			capacity:       3,
			addPairs:       [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			removeKeys:     []any{k, k + k, k + k + k},
			wantErrs:       []error{nil, nil, nil},
			checkLength:    true,
			getRemovedKeys: false,
		},
		{
			name:           "getting elements which has been removed returns nil",
			capacity:       2,
			addPairs:       [][]any{{k, v}, {k + k, v + v}},
			removeKeys:     []any{k, k + k},
			wantErrs:       []error{nil, nil},
			checkLength:    true,
			getRemovedKeys: true,
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		originalLength := c.Len()
		t.Run(tt.name, func(t *testing.T) {
			for i, key := range tt.removeKeys {
				err := c.Remove(key)
				if !errors.Is(err, tt.wantErrs[i]) {
					t.Errorf("cache.Remove() error = %v, want %v", err, tt.wantErrs[i])
					return
				}
			}
			if tt.checkLength {
				gotLength := c.Len()
				wantLength := originalLength - len(tt.removeKeys)
				if gotLength != wantLength {
					t.Errorf("incorrect cache length, got %v, want %v", gotLength, wantLength)
				}
			}
			if tt.getRemovedKeys {
				for _, key := range tt.removeKeys {
					got, gotFound := c.Get(key)
					if gotFound != false {
						t.Errorf("cache.Get() found = %v, want %v", gotFound, false)
					}
					if got != nil {
						t.Errorf("cache.Get() = %v, want %v", got, nil)
					}
				}
			}
		})
	}
}

func TestCache_Contains(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		containKeys       []any
		wantFound         []bool
		wantKeysListOrder []any
	}{
		{
			name:              "returns false for calling .Contains() for any key for empty cache",
			capacity:          1,
			addPairs:          [][]any{},
			containKeys:       []any{k, k + k, k + k + k},
			wantFound:         []bool{false, false, false},
			wantKeysListOrder: nil,
		},
		{
			name:              "returns false for keys not added to cache",
			capacity:          1,
			addPairs:          [][]any{{k, v}},
			containKeys:       []any{"nonexistent"},
			wantFound:         []bool{false},
			wantKeysListOrder: nil,
		},
		{
			name:              "returns true for keys added to cache",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			containKeys:       []any{k, k + k, k + k + k},
			wantFound:         []bool{true, true, true},
			wantKeysListOrder: nil,
		},
		{
			name:              "preserves order for items in cache after calling .Contains()",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			containKeys:       []any{k, k + k, k + k + k},
			wantFound:         []bool{true, true, true},
			wantKeysListOrder: []any{k + k + k, k + k, k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			for i, key := range tt.containKeys {
				found := c.Contains(key)
				if found != tt.wantFound[i] {
					t.Errorf("cache.Contains() found = %v, want %v", found, tt.wantFound[i])
				}
				if tt.wantKeysListOrder != nil {
					cmpCacheListOrder(t, c, tt.wantKeysListOrder)
				}
			}
		})
	}
}

func TestCache_Clear(t *testing.T) {
	tests := []struct {
		name           string
		capacity       int
		addPairs       [][]any
		checkListFront bool
	}{
		{
			name:           "can successfully clear an empty cache",
			capacity:       1,
			addPairs:       [][]any{},
			checkListFront: false,
		},
		{
			name:           "can successfully clear a cache containing a single item",
			capacity:       1,
			addPairs:       [][]any{{k, v}},
			checkListFront: false,
		},
		{
			name:           "can successfully clear a cache containing multiple items",
			capacity:       3,
			addPairs:       [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			checkListFront: true,
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			c.Clear()
			if c.Len() != 0 {
				t.Errorf("expected length s %v, got %v", 0, c.Len())
			}
			if c.lst.Len() != 0 {
				t.Errorf("expected length of c.lst.Len() is %v, got %v", 0, c.lst.Len())
			}
			if tt.checkListFront {
				if c.lst.Front() != nil {
					f := c.lst.Front().Value.(Item)
					t.Errorf("expected c.lst.Front() to be nil, got %s-%s", f.Key, f.Val)
				}
			}
		})
	}
}

func TestCache_Keys(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		wantKeysListOrder []any
	}{
		{
			name:              "returns empty list for empty cache",
			capacity:          1,
			addPairs:          [][]any{},
			wantKeysListOrder: []any{},
		},
		{
			name:              "returns keys with order preserved for items in cache",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			wantKeysListOrder: []any{k + k + k, k + k, k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			got := c.Keys()
			if len(got) != c.Len() {
				t.Errorf("expect keys length %v and cache length %v to be the same", len(got), c.Len())
			}
			if len(got) != len(tt.addPairs) {
				t.Errorf("expect keys length %v and cache length %v to be the same", len(got), len(tt.addPairs))
			}
			cmpCacheListOrder(t, c, tt.wantKeysListOrder)
		})
	}
}

func TestCache_Peek(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		peekPairs         [][]any
		wantKeysListOrder []any
	}{
		{
			name:              "returns nil item for empty cache",
			capacity:          1,
			addPairs:          [][]any{},
			peekPairs:         [][]any{{k, nil}},
			wantKeysListOrder: nil,
		},
		{
			name:              "expect pair to not be found when its not added to cache",
			capacity:          1,
			addPairs:          [][]any{{k, v}},
			peekPairs:         [][]any{{"nonexistent", nil}},
			wantKeysListOrder: nil,
		},
		{
			name:              "expect pairs to be found when they are added to cache",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			peekPairs:         [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			wantKeysListOrder: nil,
		},
		{
			name:              "expect cache list order to be preserved after peeking pairs",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			peekPairs:         [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			wantKeysListOrder: []any{k + k + k, k + k, k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			for _, pair := range tt.peekPairs {
				var (
					want          = pair[1]
					wantFound     = want != nil
					got, gotFound = c.Peek(pair[0])
				)
				if gotFound != wantFound {
					t.Errorf("cache.Peek() found = %v, want %v", gotFound, wantFound)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("cache.Peek() = %v, want %v", got, want)
				}
				if tt.wantKeysListOrder != nil {
					cmpCacheListOrder(t, c, tt.wantKeysListOrder)
				}
			}
		})
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		want              []any
		wantLength        int
		wantKeysListOrder []any
	}{
		{
			name:              `returns ("", nil, false) for empty cache`,
			capacity:          1,
			addPairs:          [][]any{},
			want:              []any{"", nil, false},
			wantLength:        0,
			wantKeysListOrder: nil,
		},
		{
			name:              "removes expected oldest value for non-empty cache",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			want:              []any{k, v, true},
			wantLength:        2,
			wantKeysListOrder: nil,
		},
		{
			name:              "cache list has expected order after removing oldest",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			want:              []any{k, v, true},
			wantLength:        2,
			wantKeysListOrder: []any{k + k + k, k + k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotVal, gotOk := c.RemoveOldest()
			if gotKey != tt.want[0] {
				t.Errorf("removed oldest key, got %v want %v", gotKey, tt.want[0])
			}
			if gotVal != tt.want[1] {
				t.Errorf("removed oldest val, got %v want %v", gotVal, tt.want[1])
			}
			if gotOk != tt.want[2] {
				t.Errorf("ok value, got %v want %v", gotOk, tt.want[2])
			}
			if c.Len() != tt.wantLength {
				t.Errorf("cache len, got %v, want %v", c.Len(), tt.wantLength)
			}
			if tt.wantKeysListOrder != nil {
				cmpCacheListOrder(t, c, tt.wantKeysListOrder)
			}
		})
	}
}

func TestCache_Resize(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		newCapacity       int
		want              int
		wantKeysListOrder []any
	}{
		{
			name:              "resizes with no diff when newCap < cap and newCap < len",
			capacity:          7,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			newCapacity:       5,
			want:              0,
			wantKeysListOrder: nil,
		},
		{
			name:              "resizes with no diff when newCap < cap and newCap = len",
			capacity:          7,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			newCapacity:       3,
			want:              0,
			wantKeysListOrder: nil,
		},
		{
			name:              "resizes with no diff when newCap = cap and newCap = len",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			newCapacity:       3,
			want:              0,
			wantKeysListOrder: nil,
		},
		{
			name:              "resizes with no diff when newCap > cap and newCap > len",
			capacity:          1,
			addPairs:          [][]any{{k, v}},
			newCapacity:       3,
			want:              0,
			wantKeysListOrder: nil,
		},
		{
			name:              "resizes with diff and correct order when newCap < cap and newCap < len",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			newCapacity:       1,
			want:              2,
			wantKeysListOrder: []any{k + k + k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			got := c.Resize(tt.newCapacity)
			if got != tt.want {
				t.Errorf("unexpected diff, got %v, want %v", got, tt.want)
			}
			if c.Cap() != tt.newCapacity {
				t.Errorf("unexpected post resize capacity, got %v, want %v", c.Cap(), tt.newCapacity)
			}
			if tt.wantKeysListOrder != nil {
				cmpCacheListOrder(t, c, tt.wantKeysListOrder)
			}
		})
	}
}

func TestCache_Len(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		addPairs [][]any
		want     int
	}{
		{
			name:     "creates cache with expected length",
			capacity: 2,
			addPairs: [][]any{{k, v}, {k + k, v + v}},
			want:     2,
		},
		{
			name:     "cache len should never exceed cache capacity",
			capacity: 2,
			addPairs: [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			want:     2,
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Len(); got != tt.want {
				t.Errorf("unexpected cache length, got %v want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Cap(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		addPairs [][]any
		want     int
	}{
		{
			name:     "creates cache with provided capacity",
			capacity: 3,
			addPairs: [][]any{{k, v}, {k + k, v + v}},
			want:     3,
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Cap(); got != tt.want {
				t.Errorf("unexpected cache capacity, got %v want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Replace(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		replacePair       []any
		wantErr           error
		wantKeysListOrder []any
	}{
		{
			name:              "replaces key value in cache when key exists",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			replacePair:       []any{k, k + v},
			wantErr:           nil,
			wantKeysListOrder: []any{k + k + k, k + k, k},
		},
		{
			name:              "replaces key value in cache when key exists",
			capacity:          3,
			addPairs:          [][]any{{k, v}, {k + k, v + v}, {k + k + k, v + v + v}},
			replacePair:       []any{k + v, k + v},
			wantErr:           errKeyNotExist,
			wantKeysListOrder: nil,
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		for _, pair := range tt.addPairs {
			k, _ := pair[0].(string)
			v, _ := pair[1].(string)
			addItems(c, [][]string{{k, v}}, t)
		}
		t.Run(tt.name, func(t *testing.T) {
			err := c.Replace(tt.replacePair[0], tt.replacePair[1])
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("unexpected error, got %v, want %v", err, tt.wantErr)
			}
			value, found := c.Peek(tt.replacePair[0])
			if tt.wantErr == nil && !found {
				t.Error("expected key to be found but was not")
			}
			if tt.wantErr == nil && !reflect.DeepEqual(value, tt.replacePair[1]) {
				t.Errorf("unexpected key value, got %v want %v", value, tt.replacePair[1])
			}
			if tt.wantKeysListOrder != nil {
				cmpCacheListOrder(t, c, tt.wantKeysListOrder)
			}
		})
	}
}

func TestCache_ClearExpiredData(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		wantLength        int
		wantKeysListOrder []any
	}{
		{
			name:              "successfully clears empty cache",
			capacity:          3,
			addPairs:          [][]any{},
			wantLength:        0,
			wantKeysListOrder: nil,
		},
		{
			name:              "clears all expired data from cache",
			capacity:          3,
			addPairs:          [][]any{{k, v, -1 * time.Hour}, {k + k, v + v, -1 * time.Hour}, {k + k + k, v + v + v, -1 * time.Hour}},
			wantLength:        0,
			wantKeysListOrder: nil,
		},
		{
			name:              "clears only expired data from cache",
			capacity:          3,
			addPairs:          [][]any{{k, v, time.Hour}, {k + k, v + v, -1 * time.Hour}, {k + k + k, v + v + v, time.Hour}},
			wantLength:        2,
			wantKeysListOrder: []any{k + k + k, k},
		},
		{
			name:              "clears no unexpired data from cache",
			capacity:          3,
			addPairs:          [][]any{{k, v, time.Hour}, {k + k, v + v, time.Hour}, {k + k + k, v + v + v, time.Hour}},
			wantLength:        3,
			wantKeysListOrder: []any{k + k + k, k + k, k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		addItemsWithExp(t, c, tt.addPairs)
		t.Run(tt.name, func(t *testing.T) {
			c.ClearExpiredData()
			if c.Len() != tt.wantLength {
				t.Errorf("unexpected length, got %v, want %v", c.Len(), tt.wantLength)
			}
			if tt.wantKeysListOrder != nil {
				cmpCacheListOrder(t, c, tt.wantKeysListOrder)
			}
		})
	}
}

func TestCache_UpdateVal(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		updatePairs       [][]any
		wantErrs          []error
		wantKeysListOrder []any
	}{
		{
			name:              "returns error if attempting to update item in an empty cache",
			capacity:          3,
			addPairs:          [][]any{},
			updatePairs:       [][]any{{k, v}},
			wantErrs:          []error{errNoKey},
			wantKeysListOrder: nil,
		},
		{
			name:              "returns error if attempt to update non-existent key",
			capacity:          3,
			addPairs:          [][]any{{k, v, time.Hour}},
			updatePairs:       [][]any{{"nonexistant", nil}},
			wantErrs:          []error{errNoKey},
			wantKeysListOrder: nil,
		},
		{
			name:              "updates item values when item key is present in cache",
			capacity:          3,
			addPairs:          [][]any{{k, v, time.Hour}, {k + k, v + v, time.Hour}, {k + k + k, v + v + v, time.Hour}},
			updatePairs:       [][]any{{k, k + v}},
			wantErrs:          []error{nil},
			wantKeysListOrder: []any{k, k + k + k, k + k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		addItemsWithExp(t, c, tt.addPairs)
		t.Run(tt.name, func(t *testing.T) {
			for i, pair := range tt.updatePairs {
				var (
					oldItem, _   = findItem(t, c, pair[0])
					newItem, err = c.UpdateVal(pair[0], pair[1])
					wantErr      = tt.wantErrs[i]
				)
				if !errors.Is(err, wantErr) {
					t.Errorf("unexpected error, got %v, want %v", err, wantErr)
					return
				}
				if wantErr == nil && newItem.Key != oldItem.Key {
					t.Errorf("unexpected updated item key, got %v, want %v", newItem.Key, oldItem.Key)
				}
				if wantErr == nil && !reflect.DeepEqual(newItem.Val, pair[1]) {
					t.Errorf("unexpected updated item value, got %v, want %v", newItem.Val, pair[1])
				}
				if wantErr == nil && newItem.Expiration != oldItem.Expiration {
					t.Errorf("unexpected updated item expiration, got %v, want %v", newItem.Expiration, oldItem.Expiration)
				}
			}
			if tt.wantKeysListOrder != nil {
				cmpCacheListOrder(t, c, tt.wantKeysListOrder)
			}
		})
	}
}

func TestCache_UpdateExpirationDate(t *testing.T) {
	tests := []struct {
		name              string
		capacity          int
		addPairs          [][]any
		updatePairs       [][]any
		wantErrs          []error
		wantKeysListOrder []any
	}{
		{
			name:              "returns error if attempting to update item in an empty cache",
			capacity:          3,
			addPairs:          [][]any{},
			updatePairs:       [][]any{{k, time.Minute}},
			wantErrs:          []error{errNoKey},
			wantKeysListOrder: nil,
		},
		{
			name:              "returns error if attempt to update non-existent key",
			capacity:          3,
			addPairs:          [][]any{{k, v, time.Hour}},
			updatePairs:       [][]any{{"nonexistant", time.Minute}},
			wantErrs:          []error{errNoKey},
			wantKeysListOrder: nil,
		},
		{
			name:              "updates item expiration time when item key is present in cache",
			capacity:          3,
			addPairs:          [][]any{{k, v, time.Hour}, {k + k, v + v, time.Hour}, {k + k + k, v + v + v, time.Hour}},
			updatePairs:       [][]any{{k, time.Minute}},
			wantErrs:          []error{nil},
			wantKeysListOrder: []any{k, k + k + k, k + k},
		},
	}
	for _, tt := range tests {
		c := createCache(tt.capacity, t)
		addItemsWithExp(t, c, tt.addPairs)
		t.Run(tt.name, func(t *testing.T) {
			for i, pair := range tt.updatePairs {
				var (
					oldItem, _   = findItem(t, c, pair[0])
					duration     = (pair[1]).(time.Duration)
					newItem, err = c.UpdateExpirationDate(pair[0], duration)
					wantErr      = tt.wantErrs[i]
				)
				if !errors.Is(err, wantErr) {
					t.Errorf("unexpected error, got %v, want %v", err, wantErr)
					return
				}
				if wantErr == nil && newItem.Key != oldItem.Key {
					t.Errorf("unexpected updated item key, got %v, want %v", newItem.Key, oldItem.Key)
				}
				if wantErr == nil && newItem.Val != oldItem.Val {
					t.Errorf("unexpected updated item value, got %v, want %v", newItem.Val, oldItem.Val)
				}
				if wantErr == nil && newItem.Expiration == oldItem.Expiration {
					t.Errorf("expected updated item expiration time %v, got %v", oldItem.Expiration, newItem.Expiration)
				}
			}
			if tt.wantKeysListOrder != nil {
				cmpCacheListOrder(t, c, tt.wantKeysListOrder)
			}
		})
	}
}

func TestItem_NotExpired(t *testing.T) {
	tests := []struct {
		name     string
		key      any
		val      any
		duration time.Duration
		want     bool
	}{
		{
			name:     "returns true for expired item",
			key:      k,
			val:      v,
			duration: -1 * time.Hour,
			want:     true,
		},
		{
			name:     "returns false for unexpired item",
			key:      k,
			val:      v,
			duration: time.Hour,
			want:     false,
		},
		{
			name:     "returns false when expiration = 0",
			key:      k,
			val:      v,
			duration: 0,
			want:     false,
		},
	}
	for _, tt := range tests {
		var exp int64
		if tt.duration != 0 {
			exp = time.Now().Add(tt.duration).UnixNano()
		}
		item := Item{
			Key:        tt.key,
			Val:        tt.val,
			Expiration: exp,
		}
		if got := item.Expired(); got != tt.want {
			t.Errorf("unexpected expired status, got %v, want %v", got, tt.want)
		}
	}
}
