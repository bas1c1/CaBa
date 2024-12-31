package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type dbslice struct {
	index int
	key   string
	value string
}

type db struct {
	name  string
	dsize int64
}

func (d db) get(key string) dbslice {
	c := cache_.search_ds(key)
	if c != nil {
		return *c
	} else {
		file, err := os.Open(d.name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			ds := parseDbSlice(scanner.Text())
			if ds.key == key {
				cache_.cache_ds(ds)
				return ds
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		return empt_ds
	}
}
