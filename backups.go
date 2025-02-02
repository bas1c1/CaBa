package main

import (
	"os"
	"strconv"
	"sync"
	"time"
)

var last_backup_name = ""

var mutex_backup sync.Mutex

func save_backup() string {
	mutex_backup.Lock()

	if last_backup_name != "" {
		os.Remove(last_backup_name)
	}
	t := "backup " + time.Now().Format("2006-01-02")
	last_backup_name = t
	maindb.save(t)
	caba_log("SAVED BACKUP TO " + t)

	mutex_backup.Unlock()
	return t
}

func save_backup_async() {
	mutex_backup.Lock()

	t := "backup " + time.Now().Format("2006-01-02") + "-"
	for i := 0; ; i++ {
		j := strconv.Itoa(i)
		if _, err := os.Stat(t + j); err != nil {
			t += j
			break
		}
	}
	last_backup_name = t
	maindb.save(t)
	caba_log("SAVED BACKUP TO " + t)

	mutex_backup.Unlock()
}

func load_backup_file(fname string) {
	mutex_backup.Lock()
	maindb.update(fname)
	maindb.save(maindb.name)
	caba_log("LOADED BACKUP FROM " + fname)
	mutex_backup.Unlock()
}
