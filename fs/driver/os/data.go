package os

import (
	"io"
	"net/http"
	"os"

	"github.com/qlova/store/fs"
)

var _ fs.Data = data{}

//data implements fs.Data
type data struct {
	path fs.Path
}

//Path returns a path representing the objects absolute location.
func (d data) Stat() (os.FileInfo, error) {
	return os.Stat(d.path.String())
}

//WriteTo implements io.WriterTo
func (d data) WriteTo(writer io.Writer) (int64, error) {
	f, err := os.Open(d.path.String())
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return io.Copy(writer, f)
}

//WriteTo implements io.WriterTo
func (d data) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, d.path.String())
}

//Delete implements fs.Data.Delete
func (d data) Delete() error {
	return os.Remove(d.path.String())
}

//WriteTo implements io.ReaderFrom
func (d data) ReadFrom(reader io.Reader) (int64, error) {
	f, err := os.Create(d.path.String())
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return io.Copy(f, reader)
}

//Path returns a path representing the objects absolute location.
func (d data) Path() fs.Path {
	return d.path
}
