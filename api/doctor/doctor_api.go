package doctorapi

import (
	"fmt"
	"net/http"

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
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorAPI struct {
	repo    *doctorrepo.DoctorRepository
	appConf *app.ApplicationConfig
}

func (a *DoctorAPI) Initialize(conf *app.ApplicationConfig, repo *doctorrepo.DoctorRepository) {
	a.repo = repo
	a.appConf = conf
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

		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		res, err := a.repo.FindById(&objectId, filterActiveAccounts)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
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

		user, err := a.repo.FindById(&objectId, false)
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

func (a *DoctorAPI) FindAllDoctors() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		pgOpts := api.ParsePaginationOptions(ctx, "doctors")
		sortOpts := api.ParseSortByOptions(ctx)
		fltrOpts := api.ParseFilterByOptions(ctx)
		ketSortBy := "$gt"

		if len(*sortOpts) == 0 {
			sortOpts = &map[string]int8{
				"_id": -1,
			}
		}

		// If sort option is given for latest with paginate ID.
		if pgOpts.PaginateId != nil {
			for key, value := range *sortOpts {
				if key == "_id" && value == -1 {
					ketSortBy = "$lt"
				}
			}
		}

		res, err := a.repo.FindAll(pgOpts, fltrOpts, sortOpts, ketSortBy)
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
			docCount, err = a.repo.GetDocumentsCount(fltrOpts)
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
