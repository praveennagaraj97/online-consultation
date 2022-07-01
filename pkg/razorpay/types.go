package razorpaypayment

import "github.com/praveennagaraj97/online-consultation/constants"

type PrefillData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Contact string `json:"contact"`
}

// Razor Pay Output
type RazorPayPaymentOutput struct {
	OrderId     *string     `json:"order_id"`
	Prefill     PrefillData `json:"prefill"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Amount      uint64      `json:"amount"`
	Currency    string      `json:"currency"`
}

type CreateRazorPayOrder struct {
	Amount         uint64
	Currency       string
	Receipt        string
	PartialPayment bool
	Notes          struct {
		// Reference name to handle via webhook capture
		PayingFor   constants.PaymentFor
		ReferenceId string
	}
}

type RazorPayWebhookEntity struct {
	ID      string `json:"id"`
	Entity  string `json:"entity"`
	Status  string `json:"status"`
	OrderId string `json:"order_id"`
	Notes   struct {
		PayingFor   constants.PaymentFor `json:"paying_for"`
		ReferenceId string               `json:"ref_id"`
	} `json:"notes"`
}

type RazorPayWebHook struct {
	Event    string   `json:"event"`
	Contains []string `json:"contains"`
	Payload  struct {
		Payment struct {
			Entity *RazorPayWebhookEntity `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
}

type RazorPayOrderDetails struct {
	ID         string      `json:"id"`
	Entity     string      `json:"entity"`
	Amount     int         `json:"amount"`
	AmountPaid int         `json:"amount_paid"`
	AmountDue  int         `json:"amount_due"`
	Currency   string      `json:"currency"`
	Receipt    string      `json:"receipt"`
	OfferID    interface{} `json:"offer_id"`
	Status     string      `json:"status"`
	Attempts   int         `json:"attempts"`
	Notes      struct {
		PayingFor string `json:"paying_for"`
		RefID     string `json:"ref_id"`
	} `json:"notes"`
	CreatedAt int `json:"created_at"`
}
