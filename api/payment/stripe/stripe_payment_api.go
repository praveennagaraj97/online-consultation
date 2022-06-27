package stripepaymentapi

import (
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/utils"
)

type StripePaymentAPI struct {
	conf *app.ApplicationConfig
}

func (a *StripePaymentAPI) Initialize(conf *app.ApplicationConfig) {
	a.conf = conf
}

func (a *StripePaymentAPI) PaymentIntentWebhook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body interface{}

		ctx.ShouldBind(&body)
		utils.PrettyPrint(body)

	}
}
