package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type kvpair struct {
	key   string
	value string
}

type dbslice struct {
	key   string
	value string
}

type db struct {
	name  string
	dsize int64
}

func (d db) remove(key string) {
	for i := 0; i < len(cache_.slices); i++ {
		if cache_.slices[i].key == key {
			cache_.delete(key)
		}
	}
	d.save(d.name)
}

func (d db) set(ks []kvpair) {
	file, err := os.OpenFile(d.name, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	writer := bufio.NewWriter(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for j := 0; j < len(ks); j++ {
		nds := dbslice{ks[j].key, ks[j].value}

		for scanner.Scan() {
			ds := parseDbSlice(scanner.Text())
			if ds.key == nds.key {
				return
			}
		}
		if len(nds.key) != 0 {
			cache_.cache_ds(nds)
			writer.WriteString("\"" + nds.key + "\";\"" + nds.value + "\"\n")
		}
	}
	writer.Flush()
}

func (d db) get(key string) *dbslice {
	c := cache_.search_ds(key)
	if c != nil {
		return c
	} else {
		d.update(d.name)
		return cache_.search_ds(key)
	}
}

func (d *db) update(fname string) {
	d.name = fname
	file, err := os.Open(d.name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	cache_.clear()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ds := parseDbSlice(scanner.Text())
		cache_.cache_ds(ds)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (d *db) save(fname string) {
	d.name = fname
	file, err := os.Create(d.name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 0; i < len(cache_.slices); i++ {
		key := cache_.slices[i].key
		value := cache_.slices[i].value

		writer.WriteString("\"" + key + "\";\"" + value + "\"\n")
	}

	if len(cache_.slices) > 0 {
		if err := writer.Flush(); err != nil {
			log.Fatal(err)
		}
	}
}
