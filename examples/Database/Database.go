package main

import "fmt"
import "qlova.store/in/bolt"

func main() {
	var ConfigStore, err = bolt.Open("bolt.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	var Config = ConfigStore.Value("config.ini")

	if err := Config.SetString("[INI]\n\ta = 1234\n"); err != nil {
		fmt.Println(err)
	}

	if s, err := Config.String(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(s)
	}
}
