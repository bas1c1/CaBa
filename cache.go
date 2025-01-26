package main

import "sync"

type cache struct {
	m     map[string]string
	mutex sync.RWMutex
}

var tmpc = cache{m: map[string]string{}}

func (c *cache) save_cache() {
	tmpc.mutex.Lock()
	defer tmpc.mutex.Unlock()

	tmpc.m = map[string]string{}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for k, v := range c.m {
		tmpc.m[k] = v
	}
	caba_log("UPDATED LAST CACHE")
}

func (c *cache) load_cache() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cache_.m = map[string]string{}

	tmpc.mutex.RLock()
	defer tmpc.mutex.RUnlock()

	for k, v := range tmpc.m {
		cache_.m[k] = v
	}
	caba_log("LOADED LAST CACHE")
}

func (c *cache) cache_ds(ds dbslice) {
	for len(c.m) >= config_.cache_size {
		for k := range c.m {
			c.delete(k)
			break
		}
	}

	c.m[ds.key] = ds.value
}

func (c *cache) clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.m = map[string]string{}
}

func (c *cache) delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.m, key)
}

func (c *cache) search_ds(key string) dbslice {
	if v, ok := c.m[key]; ok {
		return dbslice{key, v}
	}
	return _zeroslice
}
