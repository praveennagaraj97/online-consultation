package awspkg

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	logger "github.com/praveennagaraj97/online-consultation/pkg/log"
)

type AWSCredentials struct {
	S3_BUCKET_NAME           string
	S3_BUCKET_REGION         string
	AWS_ACCESS_KEY_ID        string
	AWS_SECRET_ACCESS        string
	S3_PUBLIC_DOMAIN         string
	S3_PUBLIC_ACCESS_BASEURL string
}

type AWSConfiguration struct {
	options                  *AWSCredentials
	s3Client                 *s3.Client
	snsClient                *sns.Client
	defaultConfig            *aws.Config
	S3_PUBLIC_ACCESS_BASEURL string
}

func (a *AWSConfiguration) Initialize() {

	// aws packages
	awsOptions := &AWSCredentials{
		S3_BUCKET_NAME:           env.GetEnvVariable("S3_BUCKET_NAME"),
		S3_BUCKET_REGION:         env.GetEnvVariable("S3_BUCKET_REGION"),
		AWS_ACCESS_KEY_ID:        env.GetEnvVariable("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS:        env.GetEnvVariable("AWS_SECRET_ACCESS"),
		S3_PUBLIC_ACCESS_BASEURL: env.GetEnvVariable("S3_ACCESS_BASEURL"),
	}

	if a.options == nil {
		a.options = awsOptions
		a.S3_PUBLIC_ACCESS_BASEURL = env.GetEnvVariable("S3_ACCESS_BASEURL")
	}

	creds := credentials.NewStaticCredentialsProvider(a.options.AWS_ACCESS_KEY_ID, a.options.AWS_SECRET_ACCESS, "")

	config, err := awsconf.LoadDefaultConfig(context.TODO(), awsconf.WithCredentialsProvider(creds),
		awsconf.WithRegion(a.options.S3_BUCKET_REGION))
	if err != nil {
		logger.ErrorLogFatal(err)
	}

	a.defaultConfig = &config

	a.configS3()
	a.configSNS()

	logger.PrintLog("AWS package initialized")

}
