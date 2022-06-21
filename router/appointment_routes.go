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
	)

	routes := r.engine.Group("/api/v1/appointment")

	routes.Use(r.middlewares.IsAuthorized())

	routes.POST("/schedule", r.middlewares.UserRole([]constants.UserType{constants.USER}), api.BookAnScheduledAppointment())
}
