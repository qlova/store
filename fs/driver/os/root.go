package os

import (
	"os"

	"github.com/qlova/store/fs"
)

//Open opens the given local directory as a fs.Root
//Creates the directory if it doesn't exist.
func Open(dir string) (fs.Root, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0777); err != nil {
			return fs.Root{}, err
		}
	}

	return fs.NewRoot(node{fs.Path(dir)}), nil
}
