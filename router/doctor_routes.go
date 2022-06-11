package router

import (
	doctorapi "github.com/praveennagaraj97/online-consultation/api/doctor"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) doctorRoutes() {
	api := doctorapi.DoctorAPI{}
	api.Initialize(r.app,
		r.repos.GetDoctorRepository(),
		r.repos.GetOneTimePasswordRepository(),
		r.repos.GetDoctorAppointmentSlotSetRepository())

	adminRoutes := r.engine.Group("/api/v1/admin/doctor")
	authRoutes := r.engine.Group("/api/v1/doctor/auth")
	routes := r.engine.Group("/api/v1/doctor")

	adminRoutes.Use(r.middlewares.IsAuthorized())

	adminRoutes.POST("", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.SUPER_ADMIN, constants.EDITOR}), api.AddNewDoctor())

	adminRoutes.GET("/:id", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.SUPER_ADMIN, constants.EDITOR}), api.GetDoctorById(false))
	adminRoutes.GET("", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.SUPER_ADMIN, constants.EDITOR}), api.FindAllDoctors(true))
	adminRoutes.PATCH("/:id", r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.SUPER_ADMIN, constants.EDITOR}), api.UpdateById())

	routes.GET("/:id", api.GetDoctorById(true))
	routes.GET("/activate_account/:token", api.ActivateAccount())
	routes.GET("", api.FindAllDoctors(false))

	// Profile Options
	routes.Use(r.middlewares.IsAuthorized(), r.middlewares.UserRole([]constants.UserType{constants.DOCTOR}))
	routes.PATCH("", api.UpdateById())
	routes.GET("/me", api.GetDoctorById(true))
	// Appointment Slot Set
	routes.POST("/slot_set", api.AddNewSlotSet())

	// Auth Routes
	authRoutes.POST("/send_verification_code", api.SendVerificationCode())
	authRoutes.POST("/verify_code/:verification_id", api.VerifyCode())
	authRoutes.POST("/signin_with_phonenumber", api.SignInWithPhoneNumber())
	authRoutes.POST("/signin_with_emaillink", api.SignInWithEmailLink())
	authRoutes.GET("/login_with_token/:token", api.SendLoginCredentialsForEmailLink())
	authRoutes.GET("/logout", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]constants.UserType{constants.DOCTOR}), api.Logout())

	authRoutes.GET("/refresh_token", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]constants.UserType{constants.DOCTOR}), api.RefreshToken())

}
