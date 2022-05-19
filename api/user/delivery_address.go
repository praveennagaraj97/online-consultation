package userapi

import "github.com/gin-gonic/gin"

func (a *UserAPI) AddNewAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//

	}
}

func (a *UserAPI) GetAllAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *UserAPI) GetAddressById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *UserAPI) UpdateAddressById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *UserAPI) DeleteAddressById() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (a *UserAPI) MarkAddressAsDefault() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
