package router

import (
	appointmentapi "github.com/praveennagaraj97/online-consultation/api/appointment"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) appointmentRoutes() {

	api := appointmentapi.AppointmentAPI{}
	api.Initialize(r.app,
		r.repos.GetAppointmentSlotsRepository(),
		r.repos.GetAppointmentRepository(),
		r.repos.GetConsultationRepository(),
		r.repos.GetUserRelativeRepository(),
		r.repos.GetAppointmentScheduleReminderRepository(),
		r.repos.GetUserRepository(),
	)

	razorpayRoutes := r.engine.Group("/api/v1/appointment/razorpay")

	routes := r.engine.Group("/api/v1/appointment")

	routes.Use(r.middlewares.IsAuthorized(constants.AUTH_TOKEN))

	routes.POST("/schedule", r.middlewares.UserRole([]constants.UserType{constants.USER}), api.BookScheduledAppointment())
	routes.DELETE("/schedule/cancel/:id", r.middlewares.UserRole([]constants.UserType{constants.USER}), api.CancelApponintmentBooking())
	// Razor Pay Payment
	razorpayRoutes.POST("/webhook/payment_intent", api.ConfirmScheduledAppointmentFromWebhook())
}
