package appointmentslotsapi

import (
	"github.com/praveennagaraj97/online-consultation/app"
	appointmentslotsrepo "github.com/praveennagaraj97/online-consultation/repository/appointment_slots"
)

type AppointmentSlotsAPI struct {
	repo *appointmentslotsrepo.AppointmentSlotsRepository
	conf *app.ApplicationConfig
}

func (a *AppointmentSlotsAPI) Initialize(conf *app.ApplicationConfig, repo *appointmentslotsrepo.AppointmentSlotsRepository) {
	a.conf = conf
	a.repo = repo
}
