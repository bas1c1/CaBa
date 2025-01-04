package main

type config struct {
	passkey      []byte
	cache_size   int
	priority_sys bool
}

var config_ config
