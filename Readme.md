# Qlovastore [![Godoc](https://godoc.org/github.com/qlova/store?status.svg)](https://godoc.org/github.com/qlova/store) [![Go Report Card](https://goreportcard.com/badge/github.com/qlova/store)](https://goreportcard.com/report/github.com/qlova/store)

Qlovastore is a storage abstraction library for Go.  
Both objects and values can be stored, where an object can have children and/or values.  

**Backends:**  

* Amazon S3
* BoltDB
* Operating System

**Example:**  

```Go
package main

import "fmt"
import "github.com/qlova/store/in/os"

func main() {
    var store, err = os.Open("store")
    if err != nil {
        fmt.Println(err)
        return
    }

    var ConfigFolder = store.Goto("config")
    ConfigFolder.Create()

    var Config = ConfigFolder.Value("config.ini")

    if err := Config.SetString("[INI]\n\ta = 1234\n"); err != nil {
        fmt.Println(err)
    }
}
```

**Creating a new backend:**  
On the backend side, there are Nodes, Cursor, Data and Tree interfaces to implement.  