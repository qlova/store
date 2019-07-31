package s3

import (
	"path"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qlova/store"
)

var _ = store.Node(Folder{})

//Folder is a S3 folder.
type Folder struct {
	S3
}

//Available returns if the folder is available for modifications.
func (folder Folder) Available() bool {
	var object = Object{folder.S3}
	object.Key += "\n"
	return object.Available()
}

//Children returns a list of subdirectories and files.
func (folder Folder) Children(amount ...int) store.Children {
	var list List
	list.S3 = folder.S3
	if len(amount) > 0 {
		var number = int64(amount[0])
		list.amount = &number
	}
	list.Key += "/"
	return &list
}

//Data returns the object with the given name as a store.Data.
func (folder Folder) Data(name string) store.Data {
	var object = Object{folder.S3}
	object.Key = path.Join(folder.Key, name)
	return object
}

//Goto navigates to the relative folder by path.
func (folder Folder) Goto(location string) store.Node {
	folder.Key = path.Join(folder.Key, location)
	return folder
}

//Name returns the name of this folder.
func (folder Folder) Name() string {
	return path.Base(folder.Key)
}

//Parent returns the parent folder of this folder.
func (folder Folder) Parent() store.Node {
	folder.Key = path.Join(folder.Key, "../")
	return folder
}

//Path returns a path representing the folders absolute location.
func (folder Folder) Path() string {
	return "/" + folder.Key
}

//Create creates the current S3 folder.
func (folder Folder) Create() error {
	var _, err = folder.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(folder.Bucket),
		Key:    aws.String(folder.Key + "\n"),
	})
	return err
}

//Delete is currently Unimplemented
func (folder Folder) Delete() error {
	var _, err = folder.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(folder.Bucket),
		Key:    aws.String(folder.Key + "\n"),
	})
	if err != nil {
		return err
	}

	return folder.BatchDelete.Delete(aws.BackgroundContext(), s3manager.NewDeleteListIterator(folder.S3, &s3.ListObjectsInput{
		Bucket: aws.String(folder.Bucket),
		Prefix: aws.String(folder.Key + "/"),
	}))
}
