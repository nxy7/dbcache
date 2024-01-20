// This library provides generic implementation of `Cache[T]` which can be accessed concurrently and is guaranteed to make at most
// one call to underlying data source for any given key.
package dbcache

import (
	"sync"

	"golang.org/x/sync/singleflight"
)

// Generic cache that can store any value given some string key
// it's safe to access concurrently and guarantees that data associated with a key
// will be requested at most once
type Cache[T any] struct {
	// sync map that stores data retrieved using dataRetriever
	m sync.Map

	// because sync.Map doesn't provide lazy initialization we need another
	// structure to guarantee that we'll only request user data once per user
	s singleflight.Group

	dataSource DataSource[T]
}

func MakeCache[T any](dataSource DataSource[T]) *Cache[T] {
	return &Cache[T]{
		dataSource: dataSource,
	}
}

// Get function tries to fetch data from cache or retrieves it using DataSource
// to avoid multiple requests per key in case 2 requests were made simultaneously we're
// using singleflight (key-action store) that takes care of avoiding duplicate work
func (c *Cache[T]) Get(key string) (*T, error) {
	val, ok := c.m.Load(key)
	if ok {
		as_t := val.(*T)
		return as_t, nil
	}

	v, err, _ := c.s.Do(key, func() (any, error) {
		val, err := c.dataSource.Get(key)
		if err == nil {
			c.m.Store(key, val)
		}
		return val, err
	})

	if err != nil {
		return nil, err
	}

	as_t := v.(*T)
	return as_t, nil
}
