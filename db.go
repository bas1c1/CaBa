package main

import (
	"bufio"
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
	file, err := os.OpenFile(d.name, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	writer := bufio.NewWriter(file)
	if err != nil {
		caba_err(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	cache_.delete(key)

	for scanner.Scan() {
		s := scanner.Text()
		ds := parseDbSlice(decrypt(s))
		if ds.key != key {
			if _, err := writer.WriteString(s); err != nil {
				caba_err(err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		caba_err(err)
	} else {
		caba_log("DELETED " + key)

		cache_.save_cache()
	}
}

func (d db) set(ks []dbslice) {
	file, err := os.OpenFile(d.name, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	writer := bufio.NewWriter(file)
	if err != nil {
		caba_err(err)
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
			s := "\"" + kvp.key + "\";\"" + kvp.value + "\"\n"

			if _, err := writer.WriteString(encrypt(s)); err != nil {
				caba_err(err)
			}

			caba_log("WRITED " + s)
		}
	}

	if err := writer.Flush(); err != nil {
		caba_err(err)
	}
}

func (d db) get(key string) *dbslice {
	c := cache_.search_ds(key)
	if c != nil {
		return c
	} else {
		return d.updatewds(d.name, key)
	}
}

func (d *db) updatewds(fname string, key string) *dbslice {
	file, err := os.Open(fname)
	if err != nil {
		caba_err(err)
	}
	defer file.Close()

	cache_.clear()

	var tmp dbslice

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ds := parseDbSlice(decrypt(scanner.Text()))
		if ds.key == key {
			tmp = ds
			cache_.cache_ds(ds)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		caba_err(err)
	} else {
		cache_.save_cache()

		caba_log("UPDATED " + key + " FROM " + fname)
	}

	return &tmp
}

func (d *db) update(fname string) {
	file, err := os.Open(fname)
	if err != nil {
		caba_err(err)
	}
	defer file.Close()

	cache_.clear()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ds := parseDbSlice(decrypt(scanner.Text()))
		if _, ok := cache_.m[ds.key]; ok {
			cache_.m[ds.key] = ds.value
		}
	}

	if err := scanner.Err(); err != nil {
		caba_err(err)
	} else {
		cache_.save_cache()

		caba_log("UPDATED FROM " + fname)
	}
}

func (d *db) save(fname string) {
	file, err := os.Create(fname)
	writer := bufio.NewWriter(file)
	if err != nil {
		caba_err(err)
	}
	defer file.Close()

	file2, err2 := os.Open(d.name)
	scanner := bufio.NewScanner(file2)
	if err2 != nil {
		caba_err(err2)
	}
	defer file2.Close()

	for scanner.Scan() {
		writer.WriteString(scanner.Text())
	}

	if err := writer.Flush(); err != nil {
		caba_err(err)
	} else if err := scanner.Err(); err != nil {
		caba_err(err)
	} else {
		caba_log("SAVED TO " + fname)
	}
}
