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

	adminRoutes.Use(r.middlewares.IsAuthorized(), r.middlewares.UserRole([]constants.UserType{
		constants.ADMIN, constants.SUPER_ADMIN, constants.EDITOR}))
	adminRoutes.POST("", api.AddNewDoctor())
	adminRoutes.GET("/:id", api.GetDoctorById(false))
	adminRoutes.GET("", api.FindAllDoctors(true))
	adminRoutes.PATCH("/:id", api.UpdateById())
	adminRoutes.GET("/slot_sets/:doctor_id", api.GetAllSlotSets())
	adminRoutes.POST("/slot_set/:doctor_id", api.AddNewSlotSet())
	adminRoutes.GET("/slot_set/:id/:doctor_id", api.GetSlotSetById())
	adminRoutes.PATCH("/slot_set/:id/:doctor_id", api.UpdateSlotSetById())
	adminRoutes.DELETE("/slot_set/:id/:doctor_id", api.DeleteSlotById())

	routes.GET("/:id", api.GetDoctorById(true))
	routes.GET("/activate_account/:token", api.ActivateAccount())
	routes.GET("", api.FindAllDoctors(false))

	// Profile Options
	routes.Use(r.middlewares.IsAuthorized(), r.middlewares.UserRole([]constants.UserType{constants.DOCTOR}))
	routes.PATCH("", api.UpdateById())
	routes.GET("/me", api.GetDoctorById(true))
	// Appointment Slot Set
	routes.GET("/slot_sets", api.GetAllSlotSets())
	routes.POST("/slot_set", api.AddNewSlotSet())
	routes.GET("/slot_set/:id", api.GetSlotSetById())
	routes.PATCH("/slot_set/:id", api.UpdateSlotSetById())
	routes.DELETE("/slot_set/:id", api.DeleteSlotById())

	// Auth Routes
	authRoutes.POST("/send_verification_code", api.SendVerificationCode())
	authRoutes.POST("/verify_code/:verification_id", api.VerifyCode())
	authRoutes.POST("/signin_with_phonenumber", api.SignInWithPhoneNumber())
	authRoutes.POST("/signin_with_emaillink", api.SignInWithEmailLink())
	authRoutes.GET("/login_with_token/:token", api.SendLoginCredentialsForEmailLink())
	authRoutes.GET("/logout", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]constants.UserType{constants.DOCTOR}), api.Logout())

	authRoutes.GET("/refresh_token", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]constants.UserType{constants.DOCTOR}), api.RefreshToken())

}
