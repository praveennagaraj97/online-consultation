package awspkg

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/praveennagaraj97/online-consultation/interfaces"
)

// Initialize sns client
func (a *AWSConfiguration) configSNS() {
	a.snsClient = sns.NewFromConfig(*a.defaultConfig)
}

// Send Text SMS
func (a *AWSConfiguration) SendTextSMS(payload *interfaces.SMSType) (*sns.PublishOutput, error) {

	return a.snsClient.Publish(context.TODO(), &sns.PublishInput{
		Message:     &payload.Message,
		PhoneNumber: aws.String(payload.To.Code + payload.To.Number),
	})

}
