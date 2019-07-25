package os

import (
	"errors"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/qlova/store"
)

//List is a list of folders and objects.
type List struct {
	OS
	amount int64
	index  int

	current string
	folder  *os.File
	files   []os.FileInfo

	isDir bool

	counter int
}

//Name returns the current item's name or blank if the current item is empty.
func (list *List) Name() string {
	return list.current
}

//Data returns the current item as data or nil if the current item is not data.
func (list *List) Data() store.Data {
	if !list.isDir {
		var file = File{OS: list.OS}
		file.OS.Path = path.Join(file.OS.Path, list.current)
		return file
	}
	return nil
}

//Node returns the current item as a node or nil if the current item is not a node.
func (list *List) Node() store.Node {
	if list.isDir {
		var folder = Folder{OS: list.OS}
		folder.OS.Path = path.Join(folder.OS.Path, list.current)
		return folder
	}
	return nil
}

//Reference is a reference to a child.
func (list *List) Reference() store.Reference {
	return store.Reference{
		Package:  "os",
		Internal: strconv.Itoa(list.counter),
		Fallback: path.Join(list.OS.Path, list.current),
	}
}

//SkipTo skips to the child with name 'name'.
func (list *List) SkipTo(ref store.Reference) error {
	return list.next(&ref.Internal)
}

//Next moves to the next item.
func (list *List) Next() error {
	return list.next(nil)
}

func (list *List) next(ref *string) error {
	var err error
	if list.folder == nil {
		list.folder, err = os.Open(list.OS.Path)
		if err != nil {
			return err
		}
	}

	if ref != nil {
		counter, err := strconv.Atoi(*ref)
		if err != nil {
			return err
		}

		var difference = counter - list.counter
		if difference < 0 {
			return errors.New("cannot SkipTo, reference is behind me")
		}
		if difference > 0 {
			_, err = list.folder.Readdir(difference)
			if err != nil && err != io.EOF {
				return err
			}
		}
	}

	if list.index >= len(list.files) {
		list.files, err = list.folder.Readdir(int(list.amount))
		if err != nil && err != io.EOF {
			return err
		}
	}

	if len(list.files) == 0 || list.index >= len(list.files) {
		return io.EOF
	}

	list.current = list.files[list.index].Name()
	list.index++
	list.counter++

	return nil
}
