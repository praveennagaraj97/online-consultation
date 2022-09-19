package discountcouponapi

import (
	"github.com/praveennagaraj97/online-consultation/app"
	discountcouponrepository "github.com/praveennagaraj97/online-consultation/repository/discount_coupon"
)

type DiscountCouponAPi struct {
	conf *app.ApplicationConfig
	repo *discountcouponrepository.DiscountCouponRespository
}

func (a *DiscountCouponAPi) Initialize(conf *app.ApplicationConfig, discountRepo *discountcouponrepository.DiscountCouponRespository) {
	a.conf = conf
	a.repo = discountRepo
}
