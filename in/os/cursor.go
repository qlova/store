package os

import "io"
import "os"

type Cursor struct {
	filename string

	Path string

	directory *os.File
}

//Return a unique string representing this cursor location.
func (cursor *Cursor) Name() string {
	return cursor.filename
}

//Goto the cursor location previously returned by Name()
func (cursor *Cursor) Goto(name string) {
	cursor.filename = name
}

func (cursor *Cursor) Stores(count int) []string {
	var names []string

	for i := 0; i < count; i++ {
		file, err := cursor.directory.Readdir(1)
		if err != nil {
			return nil
		}

		if len(file) > 0 && file[0].IsDir() {
			names = append(names, file[0].Name())
		}
	}

	return names
}

func (cursor *Cursor) Data(count int) []string {
	var names []string

	for i := 0; i < count; i++ {
		file, err := cursor.directory.Readdir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil
		}

		if len(file) > 0 && !file[0].IsDir() {
			names = append(names, file[0].Name())
		}
	}

	return names
}
