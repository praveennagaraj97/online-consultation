package router

import (
	consultationapi "github.com/praveennagaraj97/online-consultation/api/consultation"
	"github.com/praveennagaraj97/online-consultation/constants"
	consultationmodel "github.com/praveennagaraj97/online-consultation/models/consultation"
)

func (r *Router) consultationRoutes() {

	api := consultationapi.ConsultationAPI{}
	api.Initialize(r.app, r.repos.GetConsultationRepository())

	routes := r.engine.Group("/api/v1/consultation")
	adminRoutes := r.engine.Group("/api/v1/admin/consultation")

	adminRoutes.POST("/add_new_type", r.middlewares.IsAuthorized(),
		r.middlewares.UserRole([]constants.UserType{constants.SUPER_ADMIN}), api.AddNewConsultationType())
	adminRoutes.PATCH("/:id", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.SUPER_ADMIN}), api.UpdateById())

	routes.GET("/", api.GetAll())
	routes.GET("/instant", api.FindByType(consultationmodel.Instant))
	routes.GET("/schedule", api.FindByType(consultationmodel.Schedule))

}
