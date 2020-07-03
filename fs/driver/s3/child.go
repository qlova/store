package s3

import (
	"github.com/qlova/store/fs"
)

var _ fs.Child = child{}

type child struct {
	State
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
	c.State.Key = c.index.Path
	return data{c.State}
}

func (c child) Node() fs.Node {
	if !c.isdir {
		return nil
	}
	c.State.Key = c.index.Path
	return node{c.State}
}
