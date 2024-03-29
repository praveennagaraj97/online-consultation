package consultationapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	consultationdto "github.com/praveennagaraj97/online-consultation/dto/consultation"
	consultationmodel "github.com/praveennagaraj97/online-consultation/models/consultation"
	awspkg "github.com/praveennagaraj97/online-consultation/pkg/aws"
	consultationrepository "github.com/praveennagaraj97/online-consultation/repository/consultation"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsultationAPI struct {
	consultRepo *consultationrepository.ConsultationRepository
	appConf     *app.ApplicationConfig
}

func (a *ConsultationAPI) Initialize(appConf *app.ApplicationConfig, consultRepo *consultationrepository.ConsultationRepository) {
	a.appConf = appConf
	a.consultRepo = consultRepo
}

func (a *ConsultationAPI) AddNewConsultationType() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload consultationdto.AddConsultationDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		if errors := payload.ValidateAddConsultationDTO(); errors != nil {
			api.SendErrorResponse(ctx, errors.Message, errors.StatusCode, errors.Errors)
			return
		}

		if exist := a.consultRepo.CheckIfConsultationTypeExists(payload.Type); exist {
			api.SendErrorResponse(ctx, "Consultation with given type already exist", http.StatusUnprocessableEntity, nil)
			return
		}

		doc := &consultationmodel.ConsultationEntity{
			ID:          primitive.NewObjectID(),
			Title:       payload.Title,
			Description: payload.Description,
			Price:       payload.Price,
			ActionName:  payload.ActionName,
			Type:        payload.Type,
			Discount:    payload.Discount,
		}

		var ch chan *awspkg.S3UploadChannelResponse = make(chan *awspkg.S3UploadChannelResponse, 1)
		defer close(ch)

		a.appConf.AwsUtils.UploadImageToS3(ctx, string(constants.ConsultationIcon), doc.ID.Hex(), "icon", payload.IconWidth, payload.IconHeight, ch)

		select {
		case value, ok := <-ch:
			if ok {
				if value.Err != nil {
					api.SendErrorResponse(ctx, value.Err.Error(), http.StatusInternalServerError, nil)
					return
				} else {
					doc.Icon = value.Result
				}
			}
		default:
		}

		if err := a.consultRepo.CreateOne(doc); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})

			a.appConf.AwsUtils.DeleteAsset(&doc.Icon.OriginalImagePath)
			a.appConf.AwsUtils.DeleteAsset(&doc.Icon.BlurImagePath)
			return
		}

		doc.Icon.OriginalSrc = fmt.Sprintf("%s/%s", a.appConf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL, doc.Icon.OriginalImagePath)
		doc.Icon.BlurDataURL = fmt.Sprintf("%s/%s", a.appConf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL, doc.Icon.BlurImagePath)

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*consultationmodel.ConsultationEntity]{
			Data: doc,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Consultation type added successfully",
			},
		})

	}
}

func (a *ConsultationAPI) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get pagination/sort/filter options.
		pgOpts := api.ParsePaginationOptions(ctx, "consultation_type")
		srtOpts := bson.D{{Key: "_id", Value: -1}}
		filterOpts := api.ParseFilterByOptions(ctx)
		keySetSortby := "$lt"

		res, err := a.consultRepo.FindAll(pgOpts, &srtOpts, filterOpts, keySetSortby)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		resLen := len(res)

		// Paginate Options
		var docCount int64
		var lastResId *primitive.ObjectID

		if pgOpts.PaginateId == nil {
			docCount, err = a.consultRepo.GetDocumentsCount(filterOpts)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		if resLen > 0 {
			lastResId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, int64(resLen), lastResId, "consultation_type")

		ctx.JSON(http.StatusOK, serialize.PaginatedDataResponse[[]consultationmodel.ConsultationEntity]{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			PaginatedData: serialize.PaginatedData[[]consultationmodel.ConsultationEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "List of consultation types retrieved successfully",
				},
			},
		})

	}
}

func (a *ConsultationAPI) FindByType(typ consultationmodel.ConsultationType) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		res, err := a.consultRepo.FindByType(typ)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		res.Icon.OriginalSrc = a.appConf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + res.Icon.OriginalImagePath
		res.Icon.BlurDataURL = a.appConf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + res.Icon.BlurImagePath

		ctx.JSON(http.StatusOK, serialize.DataResponse[*consultationmodel.ConsultationEntity]{
			Data: res,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Consultation info retrieved",
			},
		})
	}
}

func (a *ConsultationAPI) UpdateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload consultationdto.EditConsultationDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		objectId := ctx.Param(("id"))

		docId, err := primitive.ObjectIDFromHex(objectId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		file, _ := ctx.FormFile("icon")

		if file != nil {
			res, err := a.consultRepo.FindById(&docId)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
				return
			}

			if res.Icon != nil {
				a.appConf.AwsUtils.DeleteAsset(&res.Icon.OriginalImagePath)
				a.appConf.AwsUtils.DeleteAsset(&res.Icon.BlurImagePath)
			}

			if payload.IconWidth == 0 || payload.IconHeight == 0 {
				payload.IconWidth = res.Icon.Width
				payload.IconHeight = res.Icon.Height
			}

			var ch chan *awspkg.S3UploadChannelResponse = make(chan *awspkg.S3UploadChannelResponse, 1)
			a.appConf.AwsUtils.UploadImageToS3(ctx, string(constants.ConsultationIcon), docId.Hex(), "icon", payload.IconWidth, payload.IconHeight, ch)

			select {
			case value, ok := <-ch:
				if ok {
					if value.Err != nil {
						api.SendErrorResponse(ctx, value.Err.Error(), http.StatusInternalServerError, nil)
						return
					} else {
						payload.Icon = value.Result
					}
				}
			default:
			}
		}

		if err := a.consultRepo.UpdateById(&docId, &payload); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}

func (a *ConsultationAPI) DeleteConsultationType() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		objectId := ctx.Param("id")

		docId, err := primitive.ObjectIDFromHex(objectId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		doc, err := a.consultRepo.FindById(&docId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		if doc.Icon != nil {
			a.appConf.AwsUtils.DeleteAsset(&doc.Icon.OriginalImagePath)
			a.appConf.AwsUtils.DeleteAsset(&doc.Icon.BlurImagePath)
		}

		if err = a.consultRepo.DeleteById(&docId); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}
