package main

import (
	"bufio"
	"io"
	"os"
)

type dbslice struct {
	key   string
	value string
}

var _zeroslice = dbslice{}

type db struct {
	name string
}

func create_db(name string) *db {
	err := os.Mkdir(name, 0755)
	_check(err)
	return &db{name}
}

func (d db) remove(key string) {
	file, err := os.Open(d.name)
	scanner := bufio.NewScanner(file)
	_check(err)
	defer file.Close()

	tempFile, err := os.CreateTemp("", "temp_*")
	_check(err)
	defer os.Remove(tempFile.Name())

	cache_.delete(key)

	writer := bufio.NewWriter(tempFile)

	for scanner.Scan() {
		s := scanner.Text()
		ds := parseDbSlice(s)
		if ds.key != key {
			if _, err := writer.WriteString(s); err != nil {
				caba_err(err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		caba_err(err)
	} else {
		file.Close()
		tempFile.Close()

		if err := os.Rename(tempFile.Name(), d.name); err != nil {
			caba_err(err)
		} else {
			caba_log("DELETED " + key)

			cache_.save_cache()
		}
	}
}

func (d db) set(ks []dbslice) {
	for _, kvp := range ks {
		fname := hashgen(kvp.key)

		file, err := os.OpenFile(d.name+"/"+fname, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			caba_err(err)
		}

		file.Write([]byte(encrypt(kvp.value)))
		defer file.Close()

		caba_log("WRITED " + "\"" + kvp.key + "\";\"" + kvp.value + "\"\n")

		c := cache_.search_ds(kvp.key)
		if c != _zeroslice {
			cache_.cache_ds(dbslice{kvp.key, kvp.value})
		}
	}
}

func (d db) get(key string) dbslice {
	c := cache_.search_ds(key)
	if c != _zeroslice {
		return c
	} else {
		return d.updatewds(key)
	}
}

func (d db) multiget(keys []string) []dbslice {
	var values []dbslice
	for _, key := range keys {
		c := cache_.search_ds(key)

		if c != _zeroslice {
			values = append(values, c)
			continue
		}

		values = append(values, d.updatewds(key))
	}
	return values
}

func (d *db) updatewds(key string) dbslice {
	k := hashgen(key)

	content, err := os.ReadFile(d.name + "/" + k)
	_check(err)

	v := decrypt(content)

	ds := dbslice{key, v}

	cache_.cache_ds(ds)

	caba_log("UPDATED " + key)

	return ds
}

func (d *db) update(fname string) {
	dir, err := os.ReadDir(fname)
	_check(err)

	for _, entry := range dir {
		ename := entry.Name()

		_, err = copy(fname+"\\"+ename, d.name+"\\"+ename)
		_check(err)
	}
}

func (d *db) save(fname string) {
	dir, err := os.ReadDir(d.name)
	_check(err)

	err = os.Mkdir(fname, 0755)
	_check(err)

	for _, entry := range dir {
		ename := entry.Name()
		_, err = copy(d.name+"\\"+ename, fname+"\\"+ename)
		_check(err)
	}
}

func copy(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	_check(err)

	if !sourceFileStat.Mode().IsRegular() {
		caba_err("Error in copy")
	}

	source, err := os.Open(src)
	_check(err)
	defer source.Close()

	destination, err := os.Create(dst)
	_check(err)
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
