package specialityapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/constants"
	specialitydto "github.com/praveennagaraj97/online-consultation/dto/speciality"
	specialitymodel "github.com/praveennagaraj97/online-consultation/models/speciality"
	awspkg "github.com/praveennagaraj97/online-consultation/pkg/aws"
	specialityrepository "github.com/praveennagaraj97/online-consultation/repository/specialities"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpecialityAPI struct {
	conf    *app.ApplicationConfig
	splRepo *specialityrepository.SpecialitysRepository
}

func (a *SpecialityAPI) Initialize(conf *app.ApplicationConfig, splRepo *specialityrepository.SpecialitysRepository) {
	a.conf = conf
	a.splRepo = splRepo
}

func (a *SpecialityAPI) AddNewSpeciality() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload specialitydto.AddSpecialityDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()
		payload.GenerateSlug()

		if err := payload.ValidateAddSpecialityDTO(); err != nil {
			api.SendErrorResponse(ctx, err.Message, err.StatusCode, err.Errors)
			return
		}

		// Check if exist by name or slug
		if exits := a.splRepo.CheckIfExists(payload.Title, payload.Slug); exits {
			api.SendErrorResponse(ctx, "Speciality with given title or slug already exist", http.StatusUnprocessableEntity, nil)
			return
		}

		doc := &specialitymodel.SpecialityEntity{
			ID:          primitive.NewObjectID(),
			Title:       payload.Title,
			Description: payload.Description,
			Slug:        payload.Slug,
		}

		var ch chan *awspkg.S3UploadChannelResponse = make(chan *awspkg.S3UploadChannelResponse, 1)
		defer close(ch)

		a.conf.AwsUtils.UploadImageToS3(ctx, string(constants.SpecialityThumbnail), doc.ID.Hex(),
			"thumbnail", payload.ThumbnailWidth, payload.ThumbnailHeight, ch)

		select {
		case value, ok := <-ch:
			if ok {
				if value.Err != nil {
					api.SendErrorResponse(ctx, value.Err.Error(), http.StatusInternalServerError, nil)
					return
				} else {
					doc.Thumbnail = value.Result
				}
			}
		default:
		}

		if err := a.splRepo.CreateOne(doc); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, &map[string]string{
				"reason": err.Error(),
			})

			a.conf.AwsUtils.DeleteAsset(&doc.Thumbnail.OriginalImagePath)
			a.conf.AwsUtils.DeleteAsset(&doc.Thumbnail.BlurImagePath)
			return
		}

		doc.Thumbnail.OriginalSrc = fmt.Sprintf("%s/%s", a.conf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL, doc.Thumbnail.OriginalImagePath)
		doc.Thumbnail.BlurDataURL = fmt.Sprintf("%s/%s", a.conf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL, doc.Thumbnail.BlurImagePath)

		ctx.JSON(http.StatusCreated, serialize.DataResponse[*specialitymodel.SpecialityEntity]{
			Data: doc,
			Response: serialize.Response{
				StatusCode: http.StatusCreated,
				Message:    "Speciality added successfully",
			},
		})

	}
}

func (a *SpecialityAPI) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		result, err := a.splRepo.FindById(&objectId)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		result.Thumbnail.OriginalSrc = a.conf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + result.Thumbnail.OriginalImagePath
		result.Thumbnail.BlurDataURL = a.conf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + result.Thumbnail.BlurImagePath

		ctx.JSON(http.StatusOK, serialize.DataResponse[*specialitymodel.SpecialityEntity]{
			Data: result,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Speciality details retrieved successfully",
			},
		})
	}
}

func (a *SpecialityAPI) GetBySlug() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slug := ctx.Param("slug")

		result, err := a.splRepo.FindBySlug(slug)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		result.Thumbnail.OriginalSrc = a.conf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + result.Thumbnail.OriginalImagePath
		result.Thumbnail.BlurDataURL = a.conf.AwsUtils.S3_PUBLIC_ACCESS_BASEURL + "/" + result.Thumbnail.BlurImagePath

		ctx.JSON(http.StatusOK, serialize.DataResponse[*specialitymodel.SpecialityEntity]{
			Data: result,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "Speciality details retrieved successfully",
			},
		})

	}
}

func (a *SpecialityAPI) UpdateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		var payload specialitydto.EditSpecialityDTO

		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}
		defer ctx.Request.Body.Close()

		var doc *specialitymodel.SpecialityEntity

		//  Check for uniuqe
		if payload.Title != "" {
			doc, err = a.splRepo.FindById(&objectId)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
				return
			}
			if payload.Title != doc.Title {
				if exists := a.splRepo.CheckIfExistsByTitle(payload.Title); exists {
					api.SendErrorResponse(ctx, "Speciality with given title already exist", http.StatusUnprocessableEntity, nil)
					return
				}
			}

		}

		if payload.Slug != "" {
			doc, err = a.splRepo.FindById(&objectId)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
				return
			}

			if payload.Slug != doc.Slug {
				if exists := a.splRepo.CheckIfExistsBySlug(payload.Slug); exists {
					api.SendErrorResponse(ctx, "Speciality with given slug already exist", http.StatusUnprocessableEntity, nil)
					return
				}
			}

		}

		// Replace image if file exists
		file, _ := ctx.FormFile("thumbnail")
		if file != nil {
			if doc == nil {
				doc, err = a.splRepo.FindById(&objectId)
				if err != nil {
					api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
					return
				}
			}

			if doc.Thumbnail != nil {
				a.conf.AwsUtils.DeleteAsset(&doc.Thumbnail.OriginalImagePath)
				a.conf.AwsUtils.DeleteAsset(&doc.Thumbnail.BlurImagePath)
			}

			var ch chan *awspkg.S3UploadChannelResponse = make(chan *awspkg.S3UploadChannelResponse, 1)
			defer close(ch)

			if payload.ThumbnailWidth == 0 || payload.ThumbnailHeight == 0 {
				payload.ThumbnailWidth = doc.Thumbnail.Width
				payload.ThumbnailHeight = doc.Thumbnail.Height
			}

			a.conf.AwsUtils.UploadImageToS3(ctx, string(constants.SpecialityThumbnail), doc.ID.Hex(),
				"thumbnail", payload.ThumbnailWidth, payload.ThumbnailHeight, ch)

			select {
			case value, ok := <-ch:
				if ok {
					if value.Err != nil {
						api.SendErrorResponse(ctx, value.Err.Error(), http.StatusInternalServerError, nil)
						return
					} else {
						payload.Thumbnail = value.Result
					}
				}
			default:
			}

		}

		if err = a.splRepo.UpdateById(&objectId, &payload); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}

func (a *SpecialityAPI) DeleteById() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		doc, err := a.splRepo.FindById(&objectId)

		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusNotFound, nil)
			return
		}

		if doc.Thumbnail != nil {
			a.conf.AwsUtils.DeleteAsset(&doc.Thumbnail.OriginalImagePath)
			a.conf.AwsUtils.DeleteAsset(&doc.Thumbnail.BlurImagePath)
		}

		if err := a.splRepo.DeleteById(&objectId); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusBadRequest, &map[string]string{
				"reason": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)

	}
}

func (a *SpecialityAPI) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		pgOpts := api.ParsePaginationOptions(ctx, "speciality")
		sortOpts := map[string]int8{"_id": -1}
		filterOptions := api.ParseFilterByOptions(ctx)
		keySetSortBy := "$lt"

		res, err := a.splRepo.FindAll(pgOpts, &sortOpts, filterOptions, keySetSortBy)

		resLen := len(res)

		// Paginate Options
		var docCount int64
		var lastResId *primitive.ObjectID

		if pgOpts.PaginateId == nil {
			docCount, err = a.splRepo.GetDocumentsCount(filterOptions)
			if err != nil {
				api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
				return
			}
		}

		if resLen > 0 {
			lastResId = &res[resLen-1].ID
		}

		count, next, prev, paginateKeySetID := api.GetPaginateOptions(docCount, pgOpts, int64(resLen), lastResId, "speciality")

		ctx.JSON(http.StatusOK, serialize.PaginatedDataResponse[[]specialitymodel.SpecialityEntity]{
			Count:            count,
			Next:             next,
			Prev:             prev,
			PaginateKeySetID: paginateKeySetID,
			PaginatedData: serialize.PaginatedData[[]specialitymodel.SpecialityEntity]{
				Data: res,
				Response: serialize.Response{
					StatusCode: http.StatusOK,
					Message:    "List of specialities retrieved successfully",
				},
			},
		})

	}
}
