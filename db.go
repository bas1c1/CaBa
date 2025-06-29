package main

import (
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func has_key(s string) bool {
    return strings.HasSuffix(s, ".key")
}

func trim_key(s string) string {
    return s[:len(s)-4]
}

type dbslice struct {
	key   string
	value string
}

var _zeroslice = dbslice{}

type db struct {
	name string
}

func create_db(name string) *db {
	if has_key(name) {
		caba_err("wrong database name")
		return nil
	}
	err := os.Mkdir(name, 0755)
	_check(err)
	return &db{name}
}

func (d db) remove(key string) {
	cache_.save_cache()

	if err := os.Remove(d.name + "/" + hashgen(key)); err != nil {
		caba_err(err)
		cache_.load_cache()
	} else {
		caba_log("DELETED " + key)

		cache_.save_cache()
	}
}

func (d db) set(ks []dbslice) {
	for _, kvp := range ks {
		fname := hashgen(kvp.key)

		file, err := os.OpenFile(d.name+"/"+fname+".key", os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			caba_err(err)
		}

		file.Write([]byte(encrypt(kvp.value)))
		defer file.Close()

		caba_log("WRITED " + "\"" + kvp.key + "\";\"" + kvp.value + "\"\n")

		c := cache_.search_ds(kvp.key)
		if c != _zeroslice && config_.caching {
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

func (d db) list() []dbslice {
	var values []dbslice

	dir, err := os.ReadDir(d.name)
	_check(err)

	for _, entry := range dir {
		ename := entry.Name()

		if has_key(ename) {
			content, err := os.ReadFile(d.name + "/" + ename)
			_check(err)

			v := decrypt(content)

			k, err := decode_base32(trim_key(ename))
			if err != nil {
				caba_err(err)
			}

			ds := dbslice{}
			
			if config_.hash_keys && utf8.ValidString(string(k)) {
				ds.key = string(k)
				ds.value = v
			} else if !config_.hash_keys && !utf8.ValidString(string(k)) {
				ds.key = decrypt(k)
				ds.value = v
			}

			values = append(values, ds)
		} else {
			values = append(values, dbslice{ename, "THIS_IS_DATABASE"})
		}
	}

	return values
}

func (d *db) updatewds(key string) dbslice {
	k := hashgen(key)

	content, err := os.ReadFile(d.name + "/" + k + ".key")
	if err != nil {
		caba_err("not found")
		return dbslice{}
	}

	v := decrypt(content)

	ds := dbslice{key, v}

	if config_.caching {
		cache_.cache_ds(ds)
	}

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
