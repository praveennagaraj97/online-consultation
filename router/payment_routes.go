package router

import stripepaymentapi "github.com/praveennagaraj97/online-consultation/api/payment/stripe"

func (r *Router) paymentRoutes() {

	stripeAPi := stripepaymentapi.StripePaymentAPI{}
	stripeAPi.Initialize(r.app)

	stripeRoutes := r.engine.Group("/api/v1/stripe")

	// Stripe Payment
	stripeRoutes.POST("/webhook/payment_intent", stripeAPi.PaymentIntentWebhook())

}
