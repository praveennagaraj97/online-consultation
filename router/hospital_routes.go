package router

import (
	hospitalapi "github.com/praveennagaraj97/online-consultation/api/hospital"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) hospitalRoutes() {

	api := hospitalapi.HospitalAPI{}
	api.Initialize(r.app, r.repos.GetHospitalRepository())

	adminRoutes := r.engine.Group("/api/v1/admin/hospital")
	routes := r.engine.Group("/api/v1/hospital")

	adminRoutes.Use(r.middlewares.IsAuthorized(constants.ADMIN_AUTH_TOKEN))

	adminRoutes.POST("", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.EDITOR, constants.SUPER_ADMIN,
	}), api.AddNewHospital())

	adminRoutes.PATCH("/:id", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.EDITOR, constants.SUPER_ADMIN,
	}), api.UpdateHospitalById())

	adminRoutes.DELETE("/:id", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.SUPER_ADMIN,
	}), api.DeleteById())

	routes.GET("/:id", api.GetHospitalById())
	routes.GET("", api.GetAllHospitals())

}
