package s3

import (
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qlova/store"
)

//List is a list of folders and objects.
type List struct {
	S3
	amount *int64

	objects *s3.ListObjectsOutput
	index   int

	current *s3.Object
}

//Name returns the current item's name or blank if the current item is empty.
func (list *List) Name() string {
	return strings.TrimSuffix(strings.TrimPrefix(*list.current.Key, list.Key), "\n")
}

//Data returns the current item as data or nil if the current item is not data.
func (list *List) Data() store.Data {
	var name = list.Name()
	if !strings.HasSuffix(name, "\n") {
		var object = Object{list.S3, list.current}
		object.Key = path.Join(object.Key, name)
		return object
	}
	return nil
}

//Node returns the current item as a node or nil if the current item is not a node.
func (list *List) Node() store.Node {
	if strings.HasSuffix(*list.current.Key, "\n") {
		var folder = Folder{list.S3}
		folder.Key = path.Join(folder.Key, list.Name())
		return folder
	}
	return nil
}

//Reference is a reference to a child.
func (list *List) Reference() store.Reference {
	return store.Reference{
		Package:  "s3",
		Internal: *list.current.Key,
		Fallback: *list.current.Key,
	}
}

//SkipTo skips to the child with name 'name'.
func (list *List) SkipTo(ref store.Reference) error {
	return list.next(&ref.Internal)
}

//Next moves to the next item.
func (list *List) Next() error {
	return list.next(nil)
}

func (list *List) next(ref *string) error {

	//First fetch.
	if list.objects == nil {
		var objects, err = list.ListObjects(&s3.ListObjectsInput{
			Bucket:    aws.String(list.Bucket),
			MaxKeys:   list.amount,
			Prefix:    aws.String(list.Key),
			Delimiter: aws.String("/"),
			Marker:    ref,
		})

		if err != nil {
			return err
		}

		list.objects = objects
	}

	if len(list.objects.Contents) == 0 {
		return io.EOF
	}

	//Further fetches.
	if list.index >= len(list.objects.Contents) {
		var marker = list.current.Key
		if ref != nil {
			marker = ref
		}
		var objects, err = list.ListObjects(&s3.ListObjectsInput{
			Bucket:    aws.String(list.Bucket),
			MaxKeys:   list.amount,
			Prefix:    aws.String(list.Key),
			Delimiter: aws.String("/"),
			Marker:    marker,
		})

		if err != nil {
			return err
		}

		list.index = 0
		list.objects = objects
		if len(list.objects.Contents) == 0 {
			return io.EOF
		}
	}

	list.current = list.objects.Contents[list.index]
	list.index++

	if *list.current.Key == list.Key {
		return list.Next()
	}

	return nil
}
