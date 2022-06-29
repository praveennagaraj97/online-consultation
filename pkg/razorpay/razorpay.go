package razorpaypayment

import (
	"fmt"

	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/razorpay/razorpay-go"
)

var client *razorpay.Client

func init() {
	client = razorpay.NewClient(env.GetEnvVariable("RAZOR_PAY_KEY_ID"), env.GetEnvVariable("RAZOR_PAY_KEY_SECRET"))

	fmt.Println(env.GetEnvVariable("RAZOR_PAY_KEY_ID"))
	fmt.Println(env.GetEnvVariable("RAZOR_PAY_KEY_SECRET"))
}

func CreateOrder(amount float64, currency string, receiptId string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"amount":          amount * 100,
		"currency":        currency,
		"receipt":         receiptId,
		"partial_payment": false,
	}

	// orderRes, _ := client.Order.Create(data, nil)

	return client.Order.Create(data, nil)
}
