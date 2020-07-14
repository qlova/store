package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"qlova.store"
)

var Public bool

//S3 is the stored S3 state for the types in this package.
type S3 struct {
	*s3.S3
	*s3manager.Uploader
	*s3manager.Downloader
	*s3manager.BatchDelete
	Bucket string
	Key    string
}

var _ = store.Tree(Bucket{})

//Bucket is an S3 bucket.
type Bucket struct {
	Folder
}

//Open the given Amazon S3 bucket as a store.
//For how to configure credentials, check https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
func Open(bucket string) (store.Root, error) {
	var session, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return store.Root{}, err
	}

	var node = Bucket{Folder{
		S3{
			S3:          s3.New(session),
			Bucket:      bucket,
			Uploader:    s3manager.NewUploader(session),
			Downloader:  s3manager.NewDownloader(session),
			BatchDelete: s3manager.NewBatchDelete(session),
		},
	}}

	//Create the bucket if it does not exist.
	if _, err := node.HeadBucket(&s3.HeadBucketInput{Bucket: aws.String(bucket)}); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NoSuchBucket":
				node.CreateBucket(&s3.CreateBucketInput{
					Bucket: aws.String(bucket),
				})
			default:
				return store.Root{}, err
			}
		} else {
			return store.Root{}, err
		}
	}

	return store.Root{Object: store.Object{
		Node: node,
	}}, nil
}
