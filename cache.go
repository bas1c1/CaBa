package main

type cache struct {
	m map[string]string
}

var tmpc = cache{m: map[string]string{}}

func (c *cache) save_cache() {
	tmpc.m = map[string]string{}

	for k, v := range c.m {
		tmpc.m[k] = v
	}
	caba_log("UPDATED LAST CACHE")
}

func (c *cache) load_cache() {
	cache_.m = map[string]string{}

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
	c.m = map[string]string{}
}

func (c *cache) delete(key string) {
	delete(c.m, key)
}

func (c *cache) search_ds(key string) dbslice {
	if v, ok := c.m[key]; ok {
		return dbslice{key, v}
	}
	return _zeroslice
}
