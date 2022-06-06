package repository

import (
	"github.com/praveennagaraj97/online-consultation/db"
	logger "github.com/praveennagaraj97/online-consultation/pkg/log"
	adminrepository "github.com/praveennagaraj97/online-consultation/repository/admin"
	consultationrepository "github.com/praveennagaraj97/online-consultation/repository/consultation"
	doctorrepo "github.com/praveennagaraj97/online-consultation/repository/doctor"
	languagerepo "github.com/praveennagaraj97/online-consultation/repository/language"
	onetimepasswordrepository "github.com/praveennagaraj97/online-consultation/repository/onetimepassword"
	specialityrepository "github.com/praveennagaraj97/online-consultation/repository/specialities"
	userrepository "github.com/praveennagaraj97/online-consultation/repository/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	userRepo                *userrepository.UserRepository
	userRelativeRepo        *userrepository.UserRelativesRepository
	userDeliveryAddressRepo *userrepository.UserDeliveryAddressRepository
	otpRepo                 *onetimepasswordrepository.OneTimePasswordRepository
	consultationRepo        *consultationrepository.ConsultationRepository
	adminRepo               *adminrepository.AdminRepository
	specialityRepo          *specialityrepository.SpecialitysRepository
	languageRepo            *languagerepo.LanguageRepository
	doctorAuthRepo          *doctorrepo.DoctorAuthRepository
}

func (r *Repository) Initialize(mgoClient *mongo.Client, dbName string) {

	// One Time Password Repo
	r.otpRepo = &onetimepasswordrepository.OneTimePasswordRepository{}
	r.otpRepo.InitializeRepository(db.OpenCollection(mgoClient, dbName, "otp"))

	// User Repo
	r.userRepo = &userrepository.UserRepository{}
	r.userRepo.InitializeRepository(db.OpenCollection(mgoClient, dbName, "user"))

	// Admin Repo
	r.adminRepo = &adminrepository.AdminRepository{}
	r.adminRepo.Initialize(db.OpenCollection(mgoClient, dbName, "admin"))

	// User relatives repo
	r.userRelativeRepo = &userrepository.UserRelativesRepository{}
	r.userRelativeRepo.InitializeRepository(db.OpenCollection(mgoClient, dbName, "user_relative"))

	// UserDelivery Address Repo
	r.userDeliveryAddressRepo = &userrepository.UserDeliveryAddressRepository{}
	r.userDeliveryAddressRepo.Initialize(db.OpenCollection(mgoClient, dbName, "user_delivery_address"))

	// Consultation Repo
	r.consultationRepo = &consultationrepository.ConsultationRepository{}
	r.consultationRepo.Initialize(db.OpenCollection(mgoClient, dbName, "consultation"))

	// Speciality Repo
	r.specialityRepo = &specialityrepository.SpecialitysRepository{}
	r.specialityRepo.Initialize(db.OpenCollection(mgoClient, dbName, "speciality"))

	// Language Repo
	r.languageRepo = &languagerepo.LanguageRepository{}
	r.languageRepo.Initialize(db.OpenCollection(mgoClient, dbName, "language"))

	// Doctor Auth Repo
	r.doctorAuthRepo = &doctorrepo.DoctorAuthRepository{}
	r.doctorAuthRepo.Initialize(db.OpenCollection(mgoClient, dbName, "doctor"))

	logger.PrintLog("Repositories initialized 📜")
}

func (r *Repository) GetUserRepository() *userrepository.UserRepository {
	return r.userRepo
}

func (r *Repository) GetUserRelativeRepository() *userrepository.UserRelativesRepository {
	return r.userRelativeRepo
}

func (r *Repository) GetUserDeliveryAddressRepository() *userrepository.UserDeliveryAddressRepository {
	return r.userDeliveryAddressRepo
}

func (r *Repository) GetOneTimePasswordRepository() *onetimepasswordrepository.OneTimePasswordRepository {
	return r.otpRepo
}

func (r *Repository) GetConsultationRepository() *consultationrepository.ConsultationRepository {
	return r.consultationRepo
}

func (r *Repository) GetAdminRepository() *adminrepository.AdminRepository {
	return r.adminRepo
}

func (r *Repository) GetSpecialityRepository() *specialityrepository.SpecialitysRepository {
	return r.specialityRepo
}

func (r *Repository) GetLanguageRepository() *languagerepo.LanguageRepository {
	return r.languageRepo
}

func (r *Repository) GetDoctorAuthRepository() *doctorrepo.DoctorAuthRepository {
	return r.doctorAuthRepo
}
