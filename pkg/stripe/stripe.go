package stripepayment

import (
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

func Initialize() {
	stripe.Key = env.GetEnvVariable("STRIPE_SECRET_KEY")
}

// Create Payment Intent
func CreatePaymentIntent(
	amount float64, currency string,
	description *string, email *string, metaData *map[string]string) (*stripe.PaymentIntent, error) {
	params := stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * 100)),
		Currency: stripe.String(currency),
		// Enable Payment methods from dashboard
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Description:  description,
		ReceiptEmail: email,
	}

	if metaData != nil {
		for key, value := range *metaData {
			params.AddMetadata(key, value)
		}
	}

	return paymentintent.New(&params)
}
