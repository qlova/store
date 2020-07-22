# Qlovastore [![Godoc](https://godoc.org/github.com/qlova/store?status.svg)](https://godoc.org/github.com/qlova/store) [![Go Report Card](https://goreportcard.com/badge/github.com/qlova/store)](https://goreportcard.com/report/github.com/qlova/store)

Qlovastore is a storage abstraction library for Go. Dealing both with filesystem and database abstraction. qlova.store/fs and qlova.store/db respectively.
The database package is still in an experimental state, the filesystem package is incomplete.

**File-system Drivers:**  

* Amazon S3 (s3)
* Operating System (os)

**Database Drivers:**  

* Postgres (postgres)

**File-system Example:**  

```Go
package main

import (
	"log"

	"qlova.store/fs/driver/os"
)

func main() {
    //Open a new fs.Root at directory called config which will be created if it doesn't exist.
	configs, err := os.Open("config")
	if err != nil {
		log.Fatalln(err)
	}

    //Set the given file to be equal to the given string.
	if err := configs.File("config.ini").SetString("[INI]\n\ta = 1234\n"); err != nil {
		log.Fatalln(err)
	}
}
```

**License**  
This work is subject to the terms of the Qlova Public
License, Version 2.0. If a copy of the QPL was not distributed with this
work, You can obtain one at https://license.qlova.org/v2

The QPL is compatible with the AGPL which is why both licenses are provided within this repository.