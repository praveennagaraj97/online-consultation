package paymentmodel

import "github.com/stripe/stripe-go/v72"

type StripePaymentModel struct {
	ID              string          `json:"id"`
	ClientSecret    string          `json:"client_secret"`
	Currency        stripe.Currency `json:"currency"`
	Amount          uint64          `json:"amount"`
	FormattedAmount uint64          `json:"formatted_amount"`
	Description     string          `json:"description"`
}

func (s *StripePaymentModel) FormatAmount() {

	switch s.Currency {
	case stripe.CurrencyINR:
		s.FormattedAmount = s.Amount / 100

	default:
		s.FormattedAmount = s.Amount / 100

	}
}
