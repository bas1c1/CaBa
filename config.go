package main

import (
	"bufio"
	"os"
)

type config struct {
	passkey    []byte
	cache_size int
}

var config_ config

func load_cfg(fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	cache_.clear()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		parseConfig(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		caba_err(err)
	}

	return nil
}
