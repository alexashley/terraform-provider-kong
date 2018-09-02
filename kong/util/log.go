package util

import (
	"log"
	"os"
)

func Log(msg string) {
	f, err := os.OpenFile("kong-provider-debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	log.SetOutput(f)

	log.Println(msg)
}
