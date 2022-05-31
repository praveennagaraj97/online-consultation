package router

import (
	specialityapi "github.com/praveennagaraj97/online-consultation/api/speciality"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) specialityRoutes() {

	api := specialityapi.SpecialityAPI{}
	api.Initialize(r.app, r.repos.GetSpecialityRepository())

	adminRoutes := r.engine.Group("/api/v1/admin/speciality")
	routes := r.engine.Group("/api/v1/speciality")

	adminRoutes.Use(r.middlewares.IsAuthorized())

	adminRoutes.POST("", r.middlewares.UserRole([]constants.UserType{constants.ADMIN, constants.SUPER_ADMIN}), api.AddNewSpeciality())

	routes.GET("", api.GetAll())

}
