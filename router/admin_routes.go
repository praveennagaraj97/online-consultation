package router

import adminapi "github.com/praveennagaraj97/online-consultation/api/admin"

func (r *Router) adminRoutes() {

	api := adminapi.AdminAPI{}
	api.Initailize(r.app, r.repos.GetAdminRepository())

	routes := r.engine.Group("/api/v1/admin")

	// Add Admistrative users
	routes.POST("/add_super_admin")

}
