package router

import (
	appointmentslotsapi "github.com/praveennagaraj97/online-consultation/api/appointment_slots"
)

func (r *Router) appointmentSlotsRoutes() {

	api := appointmentslotsapi.AppointmentSlotsAPI{}
	api.Initialize(r.app, r.repos.GetAppointmentSlotsRepository())

	// routes := r.engine.Group("/api/v1/appointment_slots")

	// doctorRoutes := r.engine.Group("/api/v1/doctor/appointment_slots")

	// doctorRoutes.Use(r.middlewares.IsAuthorized(), r.middlewares.UserRole([]constants.UserType{constants.DOCTOR}))

	// doctorRoutes.POST("/slot_set", api.AddNewSlotSet())

}
