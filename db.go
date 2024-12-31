package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type kvpair struct {
	key   string
	value string
}

type dbslice struct {
	index int
	key   string
	value string
}

type db struct {
	name  string
	dsize int64
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
		key := ks[j].key
		value := ks[j].value
		i := 1

		for scanner.Scan() {
			ds := parseDbSlice(scanner.Text())
			i++
			if ds.key == key {
				return
			}
		}

		writer.WriteString(strconv.Itoa(i) + ";" + key + ";" + value + "\n")
		writer.Flush()
	}
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
