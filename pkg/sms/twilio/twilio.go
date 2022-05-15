package twiliopkg

import (
	"fmt"

	"github.com/praveennagaraj97/online-consultation/interfaces"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	logger "github.com/praveennagaraj97/online-consultation/pkg/log"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var twilioClient *twilio.RestClient

var fromAddr string

func Initialize() {
	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.GetEnvVariable("TWILIO_ACCOUNT_SID"),
		Password: env.GetEnvVariable("TWILIO_AUTH_TOKEN"),
	})
	fromAddr = env.GetEnvVariable("TWILIO_FROM_NUMBER")

	logger.PrintLog("Twilio SMS Package Initialised 📨")
}

func SendMessage(payload *interfaces.SMSType) {
	params := &openapi.CreateMessageParams{}
	params.SetFrom(fromAddr)

	params.SetTo(payload.To)
	params.SetBody(payload.Message)

	resp, err := twilioClient.ApiV2010.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		err = nil
	} else {
		fmt.Println("Message Sid: " + *resp.Sid)
	}
}
