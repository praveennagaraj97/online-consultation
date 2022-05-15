package api

import (
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/serialize"
)

/*
SendErrorResponse is a helper function for sending error response as json with transaltions.
*/
func SendErrorResponse(ctx *gin.Context, msg string, statusCode int, errors *map[string]string) {

	ctx.JSON(statusCode, &serialize.ErrorResponse{
		Response: serialize.Response{
			StatusCode: statusCode,
			Message:    msg,
		},
		Errors: errors,
	})

}
