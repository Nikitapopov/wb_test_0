package go_cache

import goCache "github.com/patrickmn/go-cache"

type cache struct {
	cache goCache.Cache
}

type GoCacher interface {
	Add(k string, x any) error
	Get(k string) (interface{}, bool)
	GetKeys() []string
}

func NewGoCacher() GoCacher {
	return &cache{
		cache: *goCache.New(0, 0),
	}
}

func (c *cache) Add(k string, x interface{}) error {
	c.cache.Add(k, x, 0)
	return nil
}

func (c *cache) Get(k string) (interface{}, bool) {
	return c.cache.Get(k)
}

func (c *cache) GetKeys() []string {
	cacheMap := c.cache.Items()
	keys := make([]string, len(cacheMap))

	i := 0
	for k := range cacheMap {
		keys[i] = k
		i++
	}
	return keys
}
