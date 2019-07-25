package in

import (
	"github.com/qlova/store"
	"github.com/qlova/store/in/bolt"
	"github.com/qlova/store/in/os"
	"github.com/qlova/store/in/s3"
)

//Bolt opens the db file in your current directory.
//It will be created if it doesn't exist.
func Bolt(file string) (store.Root, error) {
	return bolt.Open(file)
}

//S3 opens the given Amazon S3 bucket as a store.
//For how to configure credentials, check https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
func S3(bucket string) (store.Root, error) {
	return s3.Open(bucket)
}

//OS opens the local folder 'folder', as a store.
func OS(folder string) (store.Root, error) {
	return os.Open(folder)
}
