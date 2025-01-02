package main

type cache struct {
	slices []dbslice //it's array btw
}

func (c *cache) cache_ds(ds dbslice) {
	c.slices = append(c.slices, ds)
}

func (c *cache) clear() {
	c.slices = []dbslice{}
}

func (c *cache) delete(key string) {
	for i, sl := range c.slices {
		if sl.key == key {
			c.slices = append(c.slices[:i], c.slices[i+1:]...)
		}
	}
}

func (c cache) search_ds(key string) *dbslice {
	for _, sl := range c.slices {
		if sl.key == key {
			return &sl
		}
	}
	return nil
}
