package main

import (
	"fmt"
	"os"
	"time"
)

func caba_log(msg string) {
	T := time.Now().UTC().String()

	text := T + " - " + msg + "\n"

	fmt.Print(text)

	f, err := os.OpenFile("log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}
