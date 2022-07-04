package appointmentslotsapi

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	appointmentslotsdto "github.com/praveennagaraj97/online-consultation/dto/appointment_slots"
	appointmentslotmodel "github.com/praveennagaraj97/online-consultation/models/appointment_slot"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	appointmentslotsrepo "github.com/praveennagaraj97/online-consultation/repository/appointment_slots"
	doctorrepo "github.com/praveennagaraj97/online-consultation/repository/doctor"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentSlotsAPI struct {
	repo        *appointmentslotsrepo.AppointmentSlotsRepository
	slotSetRepo *doctorrepo.DoctorAppointmentSlotSetRepository
	conf        *app.ApplicationConfig
}

func (a *AppointmentSlotsAPI) Initialize(conf *app.ApplicationConfig,
	repo *appointmentslotsrepo.AppointmentSlotsRepository, slotSetRepo *doctorrepo.DoctorAppointmentSlotSetRepository) {
	a.conf = conf
	a.repo = repo
	a.slotSetRepo = slotSetRepo
}

// Takes options as query and generates slots for doctor.
func (a *AppointmentSlotsAPI) AddNewSlots() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var slotRes *doctormodel.AppointmentSlotSetEntity

		doctorId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		var payload appointmentslotsdto.AddAppointmentSlotDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		// Get default Slot
		if payload.SlotSetId == nil {
			slotRes, err = a.slotSetRepo.FindDefault(doctorId)
			if err != nil {
				api.SendErrorResponse(ctx, "Couldn't find any default slot set", http.StatusBadRequest, nil)
				return
			}
		} else {
			slotRes, err = a.slotSetRepo.FindById(doctorId, payload.SlotSetId)
			if err != nil {
				api.SendErrorResponse(ctx, "Couldn't find any slot set with given id", http.StatusBadRequest, nil)
				return
			}
		}

		docs := generateSlotDocuments(payload.Dates, doctorId, slotRes.SlotTimes)

		if err := a.repo.CreateMany(docs); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[[]interface{}]{
			Data: docs,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Appointment Slots Created successfully",
			},
		})

	}
}

func (a *AppointmentSlotsAPI) GetAppointmentSlotById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		docId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.repo.FindByIdAndDoctorID(docId, &objectId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*appointmentslotmodel.AppointmentSlotEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Appointment slot details retrieved successfully",
			},
		})

	}
}

func (a *AppointmentSlotsAPI) GetAppointmentSlotsByDocIdAndDate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("doctor_id")
		var docId *primitive.ObjectID
		var err error

		if id != "" {
			docID, err := primitive.ObjectIDFromHex(id)
			docId = &docID
			if err != nil {
				api.SendErrorResponse(ctx, "Doctor Id should be valid primitive ObjectID", http.StatusUnprocessableEntity, nil)
				return
			}
		} else {
			docId, err = api.GetUserIdFromContext(ctx)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		dateQuery := ctx.Query("date")
		if dateQuery == "" {
			api.SendErrorResponse(ctx, "Date cannot be empty", http.StatusUnprocessableEntity, nil)
			return
		}

		date, err := time.Parse("2006-01-02", dateQuery)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		primitiveDate := primitive.NewDateTimeFromTime(date)

		res, err := a.repo.FindByDoctorIdAndDate(docId, &primitiveDate)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[[]appointmentslotmodel.AppointmentSlotEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Appointments slots retrieved successfully",
			},
		})

	}
}
