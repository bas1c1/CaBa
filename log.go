package main

import (
	"fmt"
	"os"
	"time"
)

func caba_log(msg string) {
	T := time.Now().UTC().String()

	text := "LOG - " + T + " - " + msg

	if text[len(text)-1] != '\n' {
		text += "\n"
	}

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

func caba_err(err exception) {
	T := time.Now().UTC().String()

	errt := fmt.Sprintf("%v", err)

	text := "\nERR - " + T + " - " + errt + "\n"

	if text[len(errt)-1] != '\n' {
		text += "\n"
	}

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
