package s3

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/qlova/store/fs"
)

var _ fs.Data = data{}

//data implements fs.Data
type data struct {
	State
}

//Stat returns FileInfo for this data.
func (d data) Stat() (os.FileInfo, error) {
	stat, err := d.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.Key.String()),
	})
	if err != nil {
		return nil, err
	}

	if stat.LastModified == nil || stat.ContentLength == nil {
		return nil, errors.New("stat failed")
	}

	return info{d.State, *stat, d.Key.Base(), false}, nil
}

//WriteTo implements io.WriterTo
func (d data) WriteTo(writer io.Writer) (int64, error) {
	var location = &s3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.Key.String()),
	}

	if at, ok := writer.(io.WriterAt); ok {
		return d.Downloader.Download(at, location)
	}

	remote, err := d.GetObject(location)
	if err != nil {
		return 0, err
	}
	return io.Copy(writer, remote.Body)
}

//ServeHTTP implements http.Handler
func (d data) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s",
			d.Bucket, d.Config.Region, d.Key), 303)
}

//Delete implements fs.Data.Delete
func (d data) Delete() error {
	var _, err = d.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.Key.String()),
	})
	return err
}

//ReadFrom implements io.ReaderFrom
func (d data) ReadFrom(reader io.Reader) (int64, error) {
	var acl *string
	if Public {
		acl = aws.String("public-read")
	}

	_, err := d.Upload(&s3manager.UploadInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.Key.String()),
		Body:   reader,
		ACL:    acl,
	})
	if err != nil {
		return 0, err
	}

	return -1, nil
}

//Path returns a path representing the objects absolute location.
func (d data) Path() fs.Path {
	return d.Key
}
