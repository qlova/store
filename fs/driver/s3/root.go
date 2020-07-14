package s3

import (
	"qlova.store/fs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var Public bool

//State is the stored S3 state for the types in this package.
type State struct {
	*s3.S3
	*s3manager.Uploader
	*s3manager.Downloader
	*s3manager.BatchDelete

	Bucket string
	Key    fs.Path
}

//Open opens the given local directory as a fs.Root
//Creates the directory if it doesn't exist.
func Open(bucket string) (fs.Root, error) {

	//Create the session.
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return fs.Root{}, err
	}

	var state State
	state.S3 = s3.New(session)
	state.Bucket = bucket
	state.Uploader = s3manager.NewUploader(session)
	state.Downloader = s3manager.NewDownloader(session)
	state.BatchDelete = s3manager.NewBatchDelete(session)

	var node = node{state}

	//Create the bucket if it does not exist.
	if _, err := node.HeadBucket(&s3.HeadBucketInput{Bucket: aws.String(bucket)}); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NoSuchBucket":
				_, err := node.CreateBucket(&s3.CreateBucketInput{
					Bucket: aws.String(bucket),
				})
				if err != nil {
					return fs.Root{}, err
				}
			default:
				return fs.Root{}, err
			}
		} else {
			return fs.Root{}, err
		}
	}

	return fs.NewRoot(node), nil
}
