package main

import (
	"log"

	"github.com/qlova/store/fs/driver/os"
)

func main() {
	configs, err := os.Open("config")
	if err != nil {
		log.Fatalln(err)
	}

	if err := configs.File("config.ini").SetString("[INI]\n\ta = 1234\n"); err != nil {
		log.Fatalln(err)
	}
}
