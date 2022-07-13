package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/serialize"
	"go.mongodb.org/mongo-driver/bson"
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
		return nil, errors.New("user Context is not available")
	}
	return &id, nil
}

// Parse filter Options | available options any model with ([eq], [lte], [gte], [in], [gt], [lt], [ne]).
func ParseFilterByOptions(c *gin.Context) *map[string]bson.M {
	var opts map[string]bson.M = make(map[string]bson.M)

	filterKeys := c.Request.URL.Query()

	for key := range filterKeys {

		// ignore sort and pagination keys.
		if key == "sort" || key == "per_page" || key == "page_num" || key == "paginate_id" {
			continue
		}

		filterBy := filterParamsBinder(key)
		if len(filterBy) > 0 && contains(filterBy[1]) {
			var filterValue interface{} = filterKeys.Get(key)

			// Object Id Parse
			objectId, err := primitive.ObjectIDFromHex(filterKeys.Get(key))
			if err == nil {
				filterValue = objectId
			} else {
				// Array In Operator Parse

				// Number Parse
				val, err := strconv.Atoi(filterKeys.Get(key))
				if err == nil {
					filterValue = val
				} else {
					// boolean parser
					value, err := strconv.ParseBool(filterKeys.Get(key))
					if err == nil {
						filterValue = value
					}
				}
			}

			operator := filterBy[1]

			if operator == "in" {
				filterValue = []interface{}{filterValue}
			}

			if operator == "search" {
				operator = "regex"

				filterValue = primitive.Regex{Pattern: fmt.Sprintf("%s", filterValue), Options: "i"}
			}

			opts[filterBy[0]] = bson.M{fmt.Sprintf("$%s", operator): filterValue}
		}
	}

	return &opts
}

// helpers for filtering
func filterParamsBinder(param string) []string {
	k := strings.FieldsFunc(param, func(r rune) bool {
		return r == '[' || r == ']'
	})
	if len(k) == 2 {
		return k
	}
	return nil
}

var acceptedFilterKeys []string = []string{"eq", "lte", "gte", "in", "gt", "lt", "ne", "search"}

func contains(searchterm string) bool {
	for i := 0; i < len(acceptedFilterKeys); i++ {
		if acceptedFilterKeys[i] == searchterm {
			return true
		}
	}

	return false
}
