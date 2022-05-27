package awspkg

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	logger "github.com/praveennagaraj97/online-consultation/pkg/log"
)

type AWSCredentials struct {
	S3_BUCKET_NAME           string
	S3_BUCKET_REGION         string
	S3_ACCESS_KEY_ID         string
	S3_SECRET_ACCESS         string
	S3_PUBLIC_DOMAIN         string
	S3_PUBLIC_ACCESS_BASEURL string
}

type AWSConfiguration struct {
	options  *AWSCredentials
	s3Client *s3.Client
}

func (a *AWSConfiguration) Initialize() {

	// aws packages
	awsOptions := &AWSCredentials{
		S3_BUCKET_NAME:           env.GetEnvVariable("S3_BUCKET_NAME"),
		S3_BUCKET_REGION:         env.GetEnvVariable("S3_BUCKET_REGION"),
		S3_ACCESS_KEY_ID:         env.GetEnvVariable("S3_ACCESS_KEY_ID"),
		S3_SECRET_ACCESS:         env.GetEnvVariable("S3_SECRET_ACCESS"),
		S3_PUBLIC_ACCESS_BASEURL: env.GetEnvVariable("S3_ACCESS_BASEURL"),
	}

	if a.options == nil {
		a.options = awsOptions
	}

	a.configS3()

	logger.PrintLog("AWS package initialized")

}
