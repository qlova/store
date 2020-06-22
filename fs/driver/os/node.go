package os

import (
	"os"

	"github.com/qlova/store/fs"
)

var _ fs.Node = node{}

type node struct {
	path fs.Path
}

func (n node) Create() error {
	return os.MkdirAll(n.path.String(), 0777)
}

func (n node) Delete() error {
	return os.RemoveAll(n.path.String())
}

func (n node) Data(name string) fs.Data {
	return data{n.path.Join(fs.Path(name))}
}

func (n node) Goto(path fs.Path) fs.Node {
	return node{n.path.Join(path)}
}

func (n node) Slice(offset fs.Index, length int) ([]fs.Child, error) {
	f, err := os.Open(n.path.String())
	if err != nil {
		return nil, err
	}

	if offset.Int > 0 {
		_, err = f.Readdir(offset.Int)
		if err != nil {
			return nil, err
		}
	}

	infos, err := f.Readdir(length)
	if err != nil {
		return nil, err
	}

	var children = make([]fs.Child, len(infos))

	for i, info := range infos {
		name := info.Name()
		children[i] = child{
			index: fs.Index{
				Path:   n.path.Join(fs.Path(name)),
				Int:    i,
				String: name,
			},
			isdir: info.IsDir(),
		}
	}

	return children, nil
}

func (n node) Path() fs.Path {
	return n.path
}

func (n node) Stat() (os.FileInfo, error) {
	return os.Stat(n.path.String())
}
