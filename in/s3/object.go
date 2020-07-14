package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"qlova.store"

	"github.com/aws/aws-sdk-go/service/s3"
)

var _ = store.Data(Object{})

//Object is an S3 object.
type Object struct {
	S3
	Object *s3.Object
}

//Available returns if the object is available for reading.
func (object Object) Available() bool {
	var _, err = object.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NoSuchKey":
				return false
			}
		}
		return false
	}

	return true
}

//CopyTo copies an object to the specified io.Writer
func (object Object) CopyTo(writer io.Writer) error {
	var location = &s3.GetObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	}

	if at, ok := writer.(io.WriterAt); ok {
		_, err := object.Downloader.Download(at, location)
		return err
	}
	var remote, err = object.GetObject(location)
	io.Copy(writer, remote.Body)
	return err
}

//From creates and writes an object from the specified reader.
func (object Object) From(reader io.Reader) error {
	var acl *string
	if Public {
		acl = aws.String("public-read")
	}
	_, err := object.Upload(&s3manager.UploadInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
		Body:   reader,
		ACL:    acl,
	})
	return err
}

//Path returns a path representing the objects absolute location.
func (object Object) Path() string {
	return "/" + object.Key
}

//Delete the object.
func (object Object) Delete() error {
	var _, err = object.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	})
	return err
}

//Size returns the size of the object.
func (object Object) Size() int64 {
	var head, err = object.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	})

	if err != nil {
		return -1
	}

	return *head.ContentLength
}
