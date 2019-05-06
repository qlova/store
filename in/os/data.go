package os

import "os"
import "io"

type Data struct {
	err error
	
	Path string
}

func (data *Data) Exists() bool {
	if _, err := os.Stat(data.Path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (data *Data) Create() io.WriteCloser {
	if file, err := os.Create(data.Path); err == nil {
		return file
	} else {
		data.err = err
		return nil
	}
}
func (data *Data) Delete() error {
	if err := os.Remove(data.Path); err != nil {
		return err
	}
	return nil
}

func (data *Data) Reader() io.ReadCloser {
	if file, err := os.Open(data.Path); err == nil {
		return file
	} else {
		data.err = err
		return nil
	}
}

func (data *Data) Size() int64 {
	if stat, err := os.Stat(data.Path); err != nil {
		data.err = err
		return -1
	} else {
		return stat.Size()
	}
}

func (data *Data) String() string {
	return data.Path
}

func (data *Data) Error() error {
	return data.err
} 
