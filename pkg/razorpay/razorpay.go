package razorpaypayment

import (
	"encoding/json"

	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/razorpay/razorpay-go"
)

var client *razorpay.Client

func init() {
	client = razorpay.NewClient(env.GetEnvVariable("RAZOR_PAY_KEY_ID"), env.GetEnvVariable("RAZOR_PAY_KEY_SECRET"))
}

func CreateOrder(payload CreateRazorPayOrder) (*string, error) {
	data := map[string]interface{}{
		"amount":          payload.Amount * 100,
		"currency":        payload.Currency,
		"receipt":         payload.Receipt,
		"partial_payment": payload.PartialPayment,
		"notes": map[string]string{
			"paying_for": string(payload.Notes.PayingFor),
			"ref_id":     payload.Notes.ReferenceId,
		},
	}

	res, err := client.Order.Create(data, nil)

	if res["id"] != nil {
		OrderId := res["id"].(string)
		return &OrderId, err
	}

	return nil, err
}

func GetOrderDetails(paymentId string) (*RazorPayOrderDetails, error) {

	res, err := client.Order.Fetch(paymentId, nil, nil)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(res)

	if err != nil {
		return nil, err
	}

	var details RazorPayOrderDetails

	if err := json.Unmarshal(jsonBytes, &details); err != nil {
		return nil, err
	}

	return &details, nil

}
