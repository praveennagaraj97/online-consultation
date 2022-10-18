package router

import (
	discountcouponapi "github.com/praveennagaraj97/online-consultation/api/discount_coupon"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) discountCouponRoutes() {

	api := discountcouponapi.DiscountCouponAPi{}

	api.Initialize(r.app, r.repos.GetDiscountCouponRepository())

	adminRoutes := r.engine.Group("/api/v1/admin/discount_coupon")

	// Editor and Super Admin have permission to manage coupons
	adminRoutes.Use(r.middlewares.IsAuthorized(constants.AUTH_TOKEN),
		r.middlewares.UserRole([]constants.UserType{
			constants.SUPER_ADMIN,
			constants.ADMIN,
			constants.EDITOR,
		}))

	offerRoutes := adminRoutes.Group("/offer")

	offerRoutes.POST("", api.CreateNewOffer())
	offerRoutes.PATCH(":id", api.UpdateOfferById())
	offerRoutes.GET("", api.GetAllOffers())
	offerRoutes.GET(":id", api.GetOfferById())
	offerRoutes.DELETE(":id", api.DeleteOfferById())

}
