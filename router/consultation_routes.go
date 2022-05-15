package router

import consultationapi "github.com/praveennagaraj97/online-consultation/api/consultation"

func (r *Router) consultationRoutes() {

	api := consultationapi.ConsultationAPI{}
	api.Initialize(r.repos.GetSpecialityRepository())

	routes := r.engine.Group("/api/v1/consultation")

	routes.GET("/specialities", api.GetSpecialities())

}
