package s3

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

var _ os.FileInfo = info{}

type info struct {
	State
	s3.HeadObjectOutput
	name string
	dir  bool
}

func (i info) Name() string {
	return i.name
}

func (i info) Size() int64 {
	return *i.ContentLength
}

func (i info) IsDir() bool {
	return i.dir
}

func (i info) ModTime() time.Time {
	return *i.LastModified
}

func (i info) Sys() interface{} {
	return i.State
}

func (i info) Mode() os.FileMode {
	if i.IsDir() {
		return 0777
	}
	return 0666
}
