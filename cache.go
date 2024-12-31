package main

type cache struct {
	slices []dbslice //it's array btw
}

func (c *cache) cache_ds(ds dbslice) {
	c.slices = append(c.slices, ds)
}

func (c cache) search_ds(key string) *dbslice {
	for i := 0; i < len(c.slices); i++ {
		if c.slices[i].key == key {
			return &(c.slices[i])
		}
	}
	return nil
}
