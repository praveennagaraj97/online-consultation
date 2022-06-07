package router

import (
	hospitalapi "github.com/praveennagaraj97/online-consultation/api/hospital"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) hospitalRoutes() {

	api := hospitalapi.HospitalAPI{}
	api.Initialize(r.app, r.repos.GetDoctorHospitalRepository())

	adminRoutes := r.engine.Group("/api/v1/admin/hospital")
	routes := r.engine.Group("/api/v1/hospital")

	adminRoutes.Use(r.middlewares.IsAuthorized())

	adminRoutes.POST("", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.EDITOR, constants.SUPER_ADMIN,
	}), api.AddNewHospital())

	routes.GET("/:id", api.GetHospitalById())

}
