package go_cache

import (
	"errors"
	"fmt"

	goCache "github.com/patrickmn/go-cache"
)

type mockCache struct {
	cache goCache.Cache
}

type MockGoCacher interface {
	Add(k string, x any) error
	Get(k string) (interface{}, bool)
	GetKeys() []string
}

func NewMockGoCacher() MockGoCacher {
	return &mockCache{
		cache: *goCache.New(0, 0),
	}
}

func (c *mockCache) Add(k string, x interface{}) error {
	if k == "already_existed_id" {
		return errors.New("error")
	}
	return nil
}

func (c *mockCache) Get(k string) (interface{}, bool) {
	if k == "nonexistent_item_id" {
		return nil, false
	}
	return struct{}{}, true
}

func (c *mockCache) GetKeys() []string {
	x := []string{"id_1", "id_2", "id_3"}
	fmt.Printf("x: %v\n", x)
	return x
}
