package doctorapi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	doctordto "github.com/praveennagaraj97/online-consultation/dto/doctor"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	doctormodel "github.com/praveennagaraj97/online-consultation/models/doctor"
	awspkg "github.com/praveennagaraj97/online-consultation/pkg/aws"
	mailer "github.com/praveennagaraj97/online-consultation/pkg/email"
	"github.com/praveennagaraj97/online-consultation/pkg/env"
	"github.com/praveennagaraj97/online-consultation/pkg/tokens"
	doctorrepo "github.com/praveennagaraj97/online-consultation/repository/doctor"
	onetimepasswordrepository "github.com/praveennagaraj97/online-consultation/repository/onetimepassword"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorAPI struct {
	repo            *doctorrepo.DoctorRepository
	apptSlotSetRepo *doctorrepo.DoctorAppointmentSlotSetRepository
	otpRepo         *onetimepasswordrepository.OneTimePasswordRepository
	appConf         *app.ApplicationConfig
}

func (a *DoctorAPI) Initialize(conf *app.ApplicationConfig,
	repo *doctorrepo.DoctorRepository,
	otpRepo *onetimepasswordrepository.OneTimePasswordRepository,
	apptSlotSetRepo *doctorrepo.DoctorAppointmentSlotSetRepository) {

	a.repo = repo
	a.appConf = conf
	a.otpRepo = otpRepo
	a.apptSlotSetRepo = apptSlotSetRepo

}

func (a *DoctorAPI) AddNewDoctor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload doctordto.AddNewDoctorDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		if exist := a.repo.CheckIfDoctorExistsByEmailOrPhone(payload.Email, interfaces.PhoneType{
			Code:   payload.PhoneCode,
			Number: payload.PhoneNumber,
		}); exist {
			api.SendErrorResponse(ctx, "Doctor with given credentials already exist", http.StatusUnprocessableEntity, nil)
			return
		}

		doc := doctormodel.DoctorEntity{
			ID:                 primitive.NewObjectID(),
			Name:               payload.Name,
			Email:              payload.Email,
			Phone:              &interfaces.PhoneType{Code: payload.PhoneCode, Number: payload.PhoneNumber},
			ProfessionalTitle:  payload.ProfessionalTitle,
			Experience:         payload.Experience,
			ConsultationTypeId: payload.ConsultationType,
			HospitalId:         payload.Hospital,
			SpecialityId:       payload.Speciality,
			Education:          payload.Education,
			SpokenLanguagesIds: payload.SpokenLanguages,
			RefreshToken:       "",
			IsActive:           payload.IsActive,
		}

		multipartFile, _ := ctx.FormFile("profile_pic")
		if multipartFile != nil {
			var ch chan *awspkg.S3UploadChannelResponse = make(chan *awspkg.S3UploadChannelResponse, 1)
			defer close(ch)
			a.appConf.AwsUtils.UploadImageToS3(ctx, string(constants.DoctorProfilePic), doc.ID.Hex(), "profile_pic",
				payload.ProfilePicWidth, payload.ProfilePicHeight, ch)

			select {
			case value, ok := <-ch:
				if ok {
					if value.Err != nil {
						api.SendErrorResponse(ctx, value.Err.Error(), http.StatusInternalServerError, nil)
						return
					} else {
						doc.ProfilePic = value.Result
					}
				}
			default:
			}

		}

		err := a.repo.CreateOne(&doc)

		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		// Direct Active Do Not send activation email.
		if !payload.IsActive {

			token, err := tokens.GenerateNoExpiryTokenWithCustomType(doc.ID.Hex(), "activate-doctor", "doctor")
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}

			activateLink := fmt.Sprintf("%s/?token=%s", env.GetEnvVariable("CLIENT_DOCTOR_ACTIVATE_ACCOUNT_LINK"), token)

			if err := a.appConf.EmailClient.SendNoReplyMail([]string{doc.Email},
				"Welcome to Online Consultation", "new-doctor", "welcome",
				mailer.GetNewDoctorAddedTemplateData(doc.Name, doc.ProfessionalTitle, activateLink)); err != nil {
				api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
				return
			}
		}

		if multipartFile != nil {
			doc.ProfilePic.OriginalSrc = a.appConf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + doc.ProfilePic.OriginalImagePath
			doc.ProfilePic.BlurDataURL = a.appConf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + doc.ProfilePic.BlurImagePath
		}

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*doctormodel.DoctorEntity]{
			Data: &doc,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "New doctor account has been added successfully",
			},
		})

	}
}

func (a *DoctorAPI) GetDoctorById(filterActiveAccounts bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var objectId *primitive.ObjectID
		var err error

		id := ctx.Param("id")

		if id == "" {
			objectId, err = api.GetUserIdFromContext(ctx)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		} else {
			id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
				return
			}
			objectId = &id
		}

		res, err := a.repo.FindOne(objectId, "", nil, filterActiveAccounts)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*doctormodel.DoctorEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Doctor details retrieved successfully",
			},
		})

	}
}

func (a *DoctorAPI) ActivateAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.Param("token")

		if token == "" {
			api.SendErrorResponse(ctx, "Token is required", http.StatusUnprocessableEntity, nil)
			return
		}

		claims, err := tokens.DecodeJSONWebToken(token)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		if claims.Type != "activate-doctor" {
			api.SendErrorResponse(ctx, "Provided token is invalid", http.StatusUnprocessableEntity, nil)
			return
		}

		objectId, err := primitive.ObjectIDFromHex(claims.ID)
		if err != nil {
			api.SendErrorResponse(ctx, "Token malformed", http.StatusUnprocessableEntity, nil)
			return
		}

		user, err := a.repo.FindOne(&objectId, "", nil, false)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		if user.IsActive {
			api.SendErrorResponse(ctx, "Your account is already activated", http.StatusNotAcceptable, nil)
			return
		}

		if err := a.repo.UpdateDoctorStatus(&objectId, true); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, serialize.Response{
			StatusCode: http.StatusOK,
			Message:    "Account activated successfully",
		})
	}
}

func (a *DoctorAPI) FindAllDoctors(showInActive bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		pgOpts := api.ParsePaginationOptions(ctx, "doctors")
		sortOpts := map[string]int8{}
		fltrOpts := api.ParseFilterByOptions(ctx)
		ketSortBy := "$lt"

		// Doctor Appointments Availability on particular Date
		var slotsExistsOn *primitive.DateTime = nil
		availableOn := ctx.Query("available_on")

		if availableOn != "" {
			t, err := time.Parse("2006-01-02", availableOn)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusBadGateway, nil)
				return
			}

			slotDate := primitive.NewDateTimeFromTime(t)
			slotsExistsOn = &slotDate
		}

		// Populate Next Available | Default is true.
		populateNextKey := ctx.Query("populate_next_available")

		populateNextAvailable, err := strconv.ParseBool(populateNextKey)
		if err != nil || populateNextKey == "" {
			populateNextAvailable = true
		}

		// Sort By Availability
		sortByAvailability, _ := strconv.ParseBool(ctx.Query("next_available_slot"))
		if sortByAvailability && populateNextAvailable {
			sortOpts["next_available_slot.is_available"] = -1
			sortOpts["next_available_slot.start"] = 1
		}
		sortOpts["_id"] = -1

		res, err := a.repo.FindAll(pgOpts, fltrOpts, &sortOpts, ketSortBy,
			showInActive, slotsExistsOn, populateNextAvailable)
		if err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		resLen := len(res)

		// Cached Paginate options
		var docCount int64
		var lastId *primitive.ObjectID

		if pgOpts.PaginateId == nil {
			docCount, err = a.repo.GetDocumentsCount(fltrOpts, showInActive, slotsExistsOn)
			if err != nil {
				api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
					"reason": err.Error(),
				})
				return
			}
		}

		if resLen > 0 {
			lastId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, int64(resLen), lastId, "doctors")

		ctx.JSON(http.StatusOK, serialize.PaginatedDataResponse[[]doctormodel.DoctorEntity]{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			DataResponse: serialize.DataResponse[[]doctormodel.DoctorEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "List of doctors retrieved successfully",
				},
			},
		})

	}
}

// Update doctor By ID - Admin Via Query | Doctor via context.
func (a *DoctorAPI) UpdateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var objectId *primitive.ObjectID
		var err error

		id := ctx.Param("id")
		if id != "" {
			oId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
				return
			}
			objectId = &oId
		} else {
			objectId, err = api.GetUserIdFromContext(ctx)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		var payload doctordto.EditDoctorDTO
		var doc *doctormodel.DoctorEntity

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		if errs := payload.Validate(); errs != nil {
			api.SendErrorResponse(ctx, errs.Message, errs.StatusCode, errs.Errors)
			return
		}

		// For Doctor Profile Page update ignore fields.
		if id == "" {
			payload.Email = ""
			payload.PhoneCode = ""
			payload.PhoneNumber = ""
			payload.ConsultationType = nil
			payload.Hospital = nil
			payload.Speciality = nil
			payload.SpokenLanguages = nil
		}
		if payload.Email != "" || payload.PhoneCode != "" && payload.PhoneNumber != "" {
			doc, err = a.repo.FindById(objectId)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
				return
			}
		}

		if (payload.PhoneCode != "" && payload.PhoneNumber != "") &&
			(payload.PhoneCode != doc.Phone.Code || payload.PhoneNumber != doc.Phone.Number) {
			payload.Phone = &interfaces.PhoneType{
				Code:   payload.PhoneCode,
				Number: payload.PhoneNumber,
			}

			if exists := a.repo.CheckIfDoctorExistsByPhone(*payload.Phone); exists {
				api.SendErrorResponse(ctx, "Phone number is in use by other doctor", http.StatusUnprocessableEntity, nil)
				return
			}

		}

		if payload.Email != "" && payload.Email != doc.Email {
			if exists := a.repo.CheckIfDoctorExistsByEmail(payload.Email); exists {
				api.SendErrorResponse(ctx, "Email is in use by other doctor", http.StatusUnprocessableEntity, nil)
				return
			}
		}

		// Update Profile Pic
		file, err := ctx.FormFile("profile_pic")
		if err == nil && file != nil {
			// Get existing profile
			if doc == nil {
				doc, err = a.repo.FindById(objectId)
				if err != nil {
					api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
					return
				}
			}

			// Delete existing pic
			if doc.ProfilePic != nil {
				a.appConf.AwsUtils.DeleteAsset(&doc.ProfilePic.OriginalImagePath)
				a.appConf.AwsUtils.DeleteAsset(&doc.ProfilePic.BlurImagePath)
			}

			if payload.ProfilePicWidth == 0 || payload.ProfilePicHeight == 0 {
				payload.ProfilePicWidth = doc.ProfilePic.Width
				payload.ProfilePicHeight = doc.ProfilePic.Height
			}

			var ch chan *awspkg.S3UploadChannelResponse = make(chan *awspkg.S3UploadChannelResponse, 1)
			defer close(ch)
			a.appConf.AwsUtils.UploadImageToS3(ctx, string(constants.DoctorProfilePic), doc.ID.Hex(), "profile_pic",
				payload.ProfilePicWidth, payload.ProfilePicHeight, ch)

			select {
			case value, ok := <-ch:
				if ok {
					if value.Err != nil {
						api.SendErrorResponse(ctx, value.Err.Error(), http.StatusInternalServerError, nil)
						return
					} else {
						payload.ProfilePic = value.Result
					}
				}
			default:
			}

		}

		if err = a.repo.UpdateById(objectId, &payload); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}
