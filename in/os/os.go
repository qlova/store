package os

import (
	"os"

	"github.com/qlova/store"
)

var _ = store.Tree(Root{})

//Root is a root-level folder.
type Root struct {
	Folder
}

//OS is the stored OS state for the types in this package.
type OS struct {
	Path string
}

//Open the given local folder 'name', as a store.
func Open(name string) (store.Root, error) {
	var root = Root{Folder{
		OS: OS{
			Path: name,
		},
	}}

	if _, existence := os.Stat(name); os.IsNotExist(existence) {
		if err := os.Mkdir(name, 0700); err != nil {
			return store.Root{}, err
		}
	}

	return store.Root{Object: store.Object{
		Node: root,
	}}, nil
}
