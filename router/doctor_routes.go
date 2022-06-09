package router

import (
	doctorapi "github.com/praveennagaraj97/online-consultation/api/doctor"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) doctorRoutes() {
	api := doctorapi.DoctorAPI{}
	api.Initialize(r.app, r.repos.GetDoctorRepository(), r.repos.GetOneTimePasswordRepository())

	adminRoutes := r.engine.Group("/api/v1/admin/doctor")
	authRoutes := r.engine.Group("/api/v1/doctor/auth")
	routes := r.engine.Group("/api/v1/doctor")

	adminRoutes.Use(r.middlewares.IsAuthorized())

	adminRoutes.POST("", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.SUPER_ADMIN, constants.EDITOR}), api.AddNewDoctor())

	adminRoutes.GET("/:id", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.SUPER_ADMIN, constants.EDITOR}), api.GetDoctorById(false))

	routes.GET("/:id", api.GetDoctorById(true))
	routes.GET("/activate_account/:token", api.ActivateAccount())
	routes.GET("", api.FindAllDoctors())

	// Auth Routes
	authRoutes.POST("/send_verification_code", api.SendVerificationCode())
	authRoutes.POST("/verify_code/:verification_id", api.VerifyCode())

}
