package router

import consultationapi "github.com/praveennagaraj97/online-consultation/api/consultation"

func (r *Router) consultationRoutes() {

	api := consultationapi.ConsultationAPI{}
	api.Initialize(r.app, r.repos.GetConsultationRepository())

	routes := r.engine.Group("/api/v1/consultation")

	routes.GET("/", api.GetAll())

	routes.Use(r.middlewares.IsAuthorized(), r.middlewares.UserRole([]string{"user", "developer"}))
	routes.POST("/add_new_type", api.AddNewConsultationType())
	routes.PATCH("/:id", api.UpdateById())
}
