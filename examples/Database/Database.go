package main

import "fmt"
import "github.com/qlova/store/in/bolt"

func main() {
	var Database = bolt.Open("bolt.db", 0700, nil)

	var Global = Database.Data("global")

	fmt.Println(Global.String())

	if err := Global.SetString("value"); err != nil {
		fmt.Println(err)
	}
}
