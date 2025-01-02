package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type dbslice struct {
	key   string
	value string
}

type db struct {
	name  string
	dsize int64
}

func (d db) remove(key string) {
	update_backup()
	save_backup()
	for i := 0; i < len(cache_.slices); i++ {
		if cache_.slices[i].key == key {
			cache_.delete(key)
			caba_log("DELETED " + key)
		}
	}
	d.save(d.name)
	update_backup()
}

func (d db) set(ks []dbslice) {
	update_backup()
	file, err := os.OpenFile(d.name, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	writer := bufio.NewWriter(file)
	if err != nil {
		save_backup()
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for _, kvp := range ks {
		for scanner.Scan() {
			ds := parseDbSlice(decrypt(scanner.Text()))
			if ds.key == kvp.key {
				return
			}
		}
		if len(kvp.key) != 0 {
			update_backup()
			s := "\"" + kvp.key + "\";\"" + kvp.value + "\"\n"

			cache_.cache_ds(kvp)
			if _, err := writer.WriteString(encrypt(s)); err != nil {
				save_backup()
				fmt.Println(err)
				os.Exit(1)
			}

			caba_log("WRITED " + s)
		}
	}

	update_backup()

	if err := writer.Flush(); err != nil {
		save_backup()
		fmt.Println(err)
		os.Exit(1)
	}
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
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	cache_.clear()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ds := parseDbSlice(decrypt(scanner.Text()))
		cache_.cache_ds(ds)
	}

	caba_log("UPDATED FROM " + fname)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (d *db) save(fname string) {
	file, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 0; i < len(cache_.slices); i++ {
		key := cache_.slices[i].key
		value := cache_.slices[i].value

		writer.WriteString(encrypt("\"" + key + "\";\"" + value + "\"\n"))
	}

	caba_log("SAVED TO " + fname)

	if len(cache_.slices) > 0 {
		if err := writer.Flush(); err != nil {
			log.Fatal(err)
		}
	}
}
