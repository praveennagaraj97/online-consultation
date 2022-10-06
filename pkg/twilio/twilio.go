package twiliopkg

import (
	"github.com/praveennagaraj97/online-consultation/interfaces"
)

type TwilioUtils struct {
	// client   *twilio.RestClient
	// fromAddr string
}

// Deprecated
func (t *TwilioUtils) Initialize() {
	// t.client = twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: env.GetEnvVariable("TWILIO_ACCOUNT_SID"),
	// 	Password: env.GetEnvVariable("TWILIO_AUTH_TOKEN"),
	// })
	// t.fromAddr = env.GetEnvVariable("TWILIO_FROM_NUMBER")

	// logger.PrintLog("Twilio SMS Package Initialised 📨")
}

// Deprecated
func (t *TwilioUtils) SendMessage(payload *interfaces.SMSType) error {
	// params := &openapi.CreateMessageParams{}
	// params.SetFrom(t.fromAddr)

	// params.SetTo(fmt.Sprintf("%v %v", payload.To.Code, payload.To.Number))
	// params.SetBody(payload.Message)

	// _, err := t.client.ApiV2010.CreateMessage(params)
	// if err != nil {
	// 	return err
	// }

	return nil

}
