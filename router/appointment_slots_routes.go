package router

import (
	appointmentslotsapi "github.com/praveennagaraj97/online-consultation/api/appointment_slots"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) appointmentSlotsRoutes() {

	api := appointmentslotsapi.AppointmentSlotsAPI{}
	api.Initialize(r.app, r.repos.GetAppointmentSlotsRepository(), r.repos.GetDoctorAppointmentSlotSetRepository())

	doctorRoutes := r.engine.Group("/api/v1/doctor/appointment_slots")
	routes := r.engine.Group("/api/v1/appointment_slots")

	doctorRoutes.Use(r.middlewares.IsAuthorized(constants.DOCTOR_AUTH_TOKEN), r.middlewares.UserRole([]constants.UserType{constants.DOCTOR}))

	doctorRoutes.POST("", api.AddNewSlots())
	doctorRoutes.GET("/:id", api.GetAppointmentSlotById())
	doctorRoutes.GET("", api.GetAppointmentSlotsByDocIdAndDate())

	routes.GET("/:doctor_id", api.GetAppointmentSlotsByDocIdAndDate())
	routes.GET("/slot/:id", api.GetAppointmentSlotDetailsById())
}
