package os

import "os"
import "path/filepath"

import qlovastore "github.com/qlova/store"

type Store struct {
	err error
	
	Path string
}

func (store Store) Data(name string) qlovastore.Data {
	return qlovastore.Data{
		Raw: &Data{
			Path: filepath.Join(store.Path, name),
		},
	}
}

//Open the given local directory as a store.
func Open(directory string) Store {
	
	err := os.MkdirAll(directory, 0700)
	
	return Store{
		Path: directory,
		err: err,
	}
}
