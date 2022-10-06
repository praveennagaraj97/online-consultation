package awspkg

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3DeleteObjectAPI interface {
	DeleteObject(ctx context.Context,
		params *s3.DeleteObjectInput,
		optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

func (a *AWSConfiguration) configS3() {
	a.s3Client = s3.NewFromConfig(*a.defaultConfig)
}

func (a *AWSConfiguration) UploadAsset(body io.Reader, fileKeyName string, contentType *string) (*s3.PutObjectOutput, error) {
	return manager.NewUploader(a.s3Client).S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &a.options.S3_BUCKET_NAME,
		Body:        body,
		Key:         aws.String(fileKeyName),
		ContentType: contentType,
		ACL:         "public-read",
	})
}

func deleteItem(api S3DeleteObjectAPI, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return api.DeleteObject(context.TODO(), input)
}

func (a *AWSConfiguration) DeleteAsset(objectName *string) (*s3.DeleteObjectOutput, error) {
	input := &s3.DeleteObjectInput{
		Bucket: &a.options.S3_BUCKET_NAME,
		Key:    objectName,
	}
	return deleteItem(a.s3Client, input)
}
