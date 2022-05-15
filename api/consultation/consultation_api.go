package consultationapi

import (
	specialityrepo "github.com/praveennagaraj97/online-consultation/repository/speciality"
)

type ConsultationAPI struct {
	spltyRepo *specialityrepo.SpecialityRepository
}

func (a *ConsultationAPI) Initialize(spltyRepo *specialityrepo.SpecialityRepository) {
	a.spltyRepo = spltyRepo
}
