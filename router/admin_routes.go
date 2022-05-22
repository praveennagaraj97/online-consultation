package router

import (
	adminapi "github.com/praveennagaraj97/online-consultation/api/admin"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) adminRoutes() {

	api := adminapi.AdminAPI{}
	api.Initailize(r.app, r.repos.GetAdminRepository())

	routes := r.engine.Group("/api/v1/admin")

	// Add Admistrative users
	routes.POST("/add_super_admin", api.AddNewAdmin(constants.SUPER_ADMIN))
	routes.POST("/add_admin", r.middlewares.IsAuthorized(),
		r.middlewares.UserRole([]constants.UserType{constants.SUPER_ADMIN, constants.ADMIN}), api.AddNewAdmin(constants.ADMIN))

}
