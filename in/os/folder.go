package os

import (
	"os"
	"path"

	"github.com/qlova/store"
)

var _ = store.Node(Folder{})

//Folder is a local directory.
type Folder struct {
	OS
}

//Available returns if the directory is available for modifications.
func (folder Folder) Available() bool {
	return File{folder.OS}.Available()
}

//Children returns a list of subdirectories and files.
func (folder Folder) Children(amount ...int) store.Children {
	var list List
	list.OS = folder.OS
	if len(amount) > 0 {
		var number = int64(amount[0])
		list.amount = number
	} else {
		list.amount = -1
	}
	return &list
}

//Data returns the object with the given name as a store.Data.
func (folder Folder) Data(name string) store.Data {
	var file = File{folder.OS}
	file.OS.Path = path.Join(file.OS.Path, name)
	return store.Value{Data: file}
}

//Goto navigates to the relative folder by path.
func (folder Folder) Goto(location string) store.Node {
	folder.OS.Path = path.Join(folder.OS.Path, location)
	return folder
}

//Name returns the name of this folder.
func (folder Folder) Name() string {
	return path.Base(folder.OS.Path)
}

//Parent returns the parent folder of this folder.
func (folder Folder) Parent() store.Node {
	folder.OS.Path = path.Join(folder.OS.Path, "../")
	return folder
}

//Path returns a path representing the folders absolute location.
func (folder Folder) Path() string {
	return "/" + folder.OS.Path
}

//Create creates the current S3 folder.
func (folder Folder) Create() error {
	return os.MkdirAll(folder.OS.Path, 0700)
}

//Delete is currently Unimplemented
func (folder Folder) Delete() error {
	return os.RemoveAll(folder.OS.Path)
}
