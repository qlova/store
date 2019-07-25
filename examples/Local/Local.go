package main

import "fmt"
import "github.com/qlova/store/in/os"

func main() {
	var ConfigStore, err = os.Open("config")
	if err != nil {
		fmt.Println(err)
		return
	}

	var Config = ConfigStore.Value("config.ini")

	if err := Config.SetString("[INI]\n\ta = 1234\n"); err != nil {
		fmt.Println(err)
	}
}
