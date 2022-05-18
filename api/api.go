package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Get user if parsed from req context which will be set by middleware.
func GetUserIdFromContext(c *gin.Context) (*primitive.ObjectID, error) {
	value, _ := c.Get("id")

	id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", value))
	if err != nil {
		return nil, errors.New("User Context is not available")
	}
	return &id, nil
}
