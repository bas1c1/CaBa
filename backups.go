package main

import (
	"os"
	"time"
)

type backup struct {
	c  cache
	ln string
}

var backup_last backup = backup{cache{[]dbslice{}}, ""}

func update_backup() {
	backup_last.c.slices = []dbslice{}
	backup_last.c.slices = append(backup_last.c.slices, cache_.slices...)
	caba_log("UPDATED LAST BACKUP")
}

func save_backup() string {
	load_backup()
	if backup_last.ln != "" {
		os.Remove(backup_last.ln)
	}
	t := "backup " + time.Now().Format("2006-01-02")
	backup_last.ln = t
	maindb.save(t)
	caba_log("SAVED BACKUP TO " + t)
	return t
}

func load_backup_file(fname string) {
	maindb.update(fname)
	update_backup()
	maindb.save(maindb.name)
	caba_log("LOADED BACKUP FROM " + fname)
}

func load_backup() {
	cache_.slices = []dbslice{}
	cache_.slices = append(cache_.slices, backup_last.c.slices...)
	caba_log("LOADED LAST BACKUP")
}
