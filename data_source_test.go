package dbcache_test

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/go-faker/faker/v4"
	"github.com/nxy7/dbcache"
)

var _ dbcache.DataSource[User] = &FakeDataSource{}

type FakeDataSource struct {
	m           map[string]User
	mutex       sync.RWMutex
	accessCount atomic.Int64
	// flag indicating whether Get calls to this source should return errors
	shouldFail bool
}

func MakeFakeDataSource(userAmount uint32, shouldFail bool) *FakeDataSource {
	f := FakeDataSource{
		m:          make(map[string]User),
		shouldFail: shouldFail,
	}
	for i := uint32(0); i < userAmount; i++ {
		f.m[fmt.Sprintf("%v", i)] = GenerateRandomUser()
	}
	return &f
}

func (f *FakeDataSource) Get(key string) (*User, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	f.accessCount.Add(1)

	if f.shouldFail {
		return nil, fmt.Errorf("simulated source error")
	}

	u, ok := f.m[key]
	if ok {
		return &u, nil
	} else {
		return nil, nil
	}
}

func (f *FakeDataSource) AccessCount() int {
	return int(f.accessCount.Load())
}

func (f *FakeDataSource) GetAllKeys() []string {
	var keys []string
	for k := range f.m {
		keys = append(keys, k)
	}
	return keys
}

// Example data retrieved by DataSource[User] used in tests
type User struct {
	Name string
	Age  uint32
}

func GenerateRandomUser() User {
	age := rand.Uint32() % 100
	name := faker.Name()
	return User{
		Age:  age,
		Name: name,
	}
}
