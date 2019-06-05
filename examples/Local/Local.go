package main

import "fmt"
import "github.com/qlova/store/in/os"

func main() {
	var ConfigStore = os.Open("config")

	var Config = ConfigStore.Data("config.ini")

	if err := Config.SetString("[INI]\n\ta = 1234"); err != nil {
		fmt.Println(err)
	}
}
