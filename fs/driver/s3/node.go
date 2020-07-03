package s3

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/qlova/store/fs"
)

var _ fs.Node = node{}

type node struct {
	State
}

func (n node) Create() error {
	if base := path.Base(n.Key.String()); base == "_" {
		return errors.New(base + ": invalid folder name")
	}

	var parent = n.Goto("..").(node)
	if parent.Key == "" || parent.Key == "../" || parent.Key == ".." || parent.Key == "." {
		return nil
	}
	if err := parent.Create(); err != nil {
		return err
	}

	var _, err = n.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(n.Bucket),
		Key:    aws.String(n.Key.String() + "\n"),
	})
	return err
}

func (n node) Delete() error {
	var _, err = n.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(n.Bucket),
		Key:    aws.String(n.Key.String() + "\n"),
	})
	if err != nil {
		return err
	}

	return n.BatchDelete.Delete(aws.BackgroundContext(), s3manager.NewDeleteListIterator(n.S3, &s3.ListObjectsInput{
		Bucket: aws.String(n.Bucket),
		Prefix: aws.String(n.Key.String() + "/"),
	}))
}

func (n node) Data(name string) fs.Data {
	n.Key = n.Key.Join(fs.Path(name))
	return data{n.State}
}

func (n node) Goto(path fs.Path) fs.Node {
	n.Key = n.Key.Join(path)
	return n
}

func (n node) Slice(offset fs.Index, length int) ([]fs.Child, error) {
	var marker *string
	if offset.Path != "" {
		marker = aws.String(offset.Path.String())
	}

	objects, err := n.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(n.Bucket),
		MaxKeys:   aws.Int64(int64(length)),
		Prefix:    aws.String(n.Key.String()),
		Delimiter: aws.String("/"),
		Marker:    marker,
	})
	if err != nil {
		return nil, err
	}

	var children = make([]child, len(objects.Contents))

	for i, object := range objects.Contents {
		children[i] = child{
			n.State,
			fs.Index{
				Path: fs.Path(*object.Key),
			},
			strings.HasSuffix(*object.Key, "\n"),
		}
	}

	var converted = make([]fs.Child, len(children))
	for i := range children {
		converted[i] = &children[i]
	}

	return converted, nil
}

func (n node) Path() fs.Path {
	return n.Key
}

func (n node) Stat() (os.FileInfo, error) {
	stat, err := n.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(n.Bucket),
		Key:    aws.String(n.Key.String()),
	})
	if err != nil {
		return nil, err
	}

	if stat.LastModified == nil || stat.ContentLength == nil {
		return nil, errors.New("stat failed")
	}

	return info{n.State, *stat, n.Key.Base(), true}, nil
}
