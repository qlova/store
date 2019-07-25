package os

import (
	"io"
	"os"

	"github.com/qlova/store"
)

var _ = store.Data(File{})

//File is a local file.
type File struct {
	OS
}

//Available returns if the file is available for reading.
func (file File) Available() bool {
	if _, err := os.Stat(file.OS.Path); err != nil {
		return false
	}
	return true
}

//CopyTo copies a file to the specified io.Writer
func (file File) CopyTo(writer io.Writer) error {
	var data, err = os.Open(file.OS.Path)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, data)

	return err
}

//From creates and writes a file from the specified reader.
func (file File) From(reader io.Reader) error {
	var data, err = os.Create(file.OS.Path)
	if err != nil {
		return err
	}

	_, err = io.Copy(data, reader)

	return err
}

//Path returns a path representing the objects absolute location.
func (file File) Path() string {
	return "/" + file.OS.Path
}

//Delete the object.
func (file File) Delete() error {
	return os.Remove(file.OS.Path)
}

//Size returns the size of the object.
func (file File) Size() int64 {
	stat, err := os.Stat(file.OS.Path)
	if err != nil {
		return -1
	}
	return stat.Size()
}
