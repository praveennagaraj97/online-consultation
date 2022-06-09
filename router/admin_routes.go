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
	if r.app.Environment == "development" {
		routes.POST("/add_super_admin", api.AddNewAdmin(constants.SUPER_ADMIN))
	}

	routes.POST("/login", api.Login())
	routes.POST("/forgot_password", api.ForgotPassword())
	routes.POST("/reset_password/:token", api.ResetPassword())

	routes.Use(r.middlewares.IsAuthorized())

	routes.POST("/add_admin", r.middlewares.UserRole([]constants.UserType{constants.SUPER_ADMIN}),
		api.AddNewAdmin(constants.ADMIN))
	routes.POST("/add_editor", r.middlewares.UserRole([]constants.UserType{constants.ADMIN, constants.SUPER_ADMIN}),
		api.AddNewAdmin(constants.EDITOR))
	routes.DELETE("/delete_admin", r.middlewares.UserRole([]constants.UserType{constants.SUPER_ADMIN}), api.DeleteUser())
	routes.DELETE("/change_role", r.middlewares.UserRole([]constants.UserType{constants.SUPER_ADMIN}), api.ChangeRole())
	routes.DELETE("/delete_editor", r.middlewares.UserRole([]constants.UserType{constants.ADMIN}), api.DeleteUser())

	routes.Use(r.middlewares.UserRole([]constants.UserType{constants.ADMIN, constants.EDITOR, constants.SUPER_ADMIN}))
	routes.GET("/refresh_token", api.RefreshToken())
	routes.PATCH("/update_password", api.UpdatePassword())
	routes.GET("/logout", api.Logout())
}
