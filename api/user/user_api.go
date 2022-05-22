package userapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/api"
	"github.com/praveennagaraj97/online-consultation/app"
	userdto "github.com/praveennagaraj97/online-consultation/dto/user"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
	"github.com/praveennagaraj97/online-consultation/serialize"

	onetimepasswordrepository "github.com/praveennagaraj97/online-consultation/repository/onetimepassword"
	userrepository "github.com/praveennagaraj97/online-consultation/repository/user"
)

type UserAPI struct {
	appConf       *app.ApplicationConfig
	userRepo      *userrepository.UserRepository
	otpRepo       *onetimepasswordrepository.OneTimePasswordRepository
	relativeRepo  *userrepository.UserRelativesRepository
	delvrAddrRepo *userrepository.UserDeliveryAddressRepository
}

func (a *UserAPI) Initialize(appConf *app.ApplicationConfig,
	repo *userrepository.UserRepository,
	otpRepo *onetimepasswordrepository.OneTimePasswordRepository,
	relativeRepo *userrepository.UserRelativesRepository,
	dlvrAddrRepo *userrepository.UserDeliveryAddressRepository,
) {
	a.appConf = appConf
	a.userRepo = repo
	a.otpRepo = otpRepo
	a.relativeRepo = relativeRepo
	a.delvrAddrRepo = dlvrAddrRepo
}

// Get Users Details
func (a *UserAPI) GetUserDetails() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest, nil)
			return
		}

		user, err := a.userRepo.FindById(userId)
		if err != nil {
			api.SendErrorResponse(ctx, "User not found", http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusOK, serialize.DataResponse[*usermodel.UserEntity]{
			Data: user,
			Response: serialize.Response{
				StatusCode: http.StatusOK,
				Message:    "User Details retrieved successfully",
			},
		})

	}
}

// Update user details by ID
func (a *UserAPI) UpdateUserDetails() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload userdto.UpdateUserDTO
		if err := ctx.ShouldBind(&payload); err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusUnprocessableEntity, nil)
			return
		}

		userId, err := api.GetUserIdFromContext(ctx)
		if err != nil {
			api.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		if err = a.userRepo.UpdateById(userId, &payload); err != nil {
			api.SendErrorResponse(ctx, "Something went wrong", http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}
