package os

import (
	"qlova.store/fs"
)

var _ fs.Child = child{}

type child struct {
	index fs.Index
	isdir bool
}

func (c child) Index() fs.Index {
	return c.index
}

func (c child) Data() fs.Data {
	if c.isdir {
		return nil
	}
	return data{c.index.Path}
}

func (c child) Node() fs.Node {
	if !c.isdir {
		return nil
	}
	return node{c.index.Path}
}
