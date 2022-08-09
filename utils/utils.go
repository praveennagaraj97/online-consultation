package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/constants"
	"github.com/praveennagaraj97/online-consultation/interfaces"
	logger "github.com/praveennagaraj97/online-consultation/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndex - creates an index for a collection
func CreateIndex(collection *mongo.Collection, keys bson.D, indexName string, unique bool) bool {

	var indexOptions *options.IndexOptions = &options.IndexOptions{}

	indexOptions.Unique = &unique
	indexOptions.Name = options.Index().SetName(indexName).Name

	// 1. Field key
	mod := mongo.IndexModel{
		Keys:    keys,
		Options: indexOptions,
	}

	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {

		logger.ErrorLogFatal(err.Error())
		return false
	}

	return true
}

func GenerateRandomCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var codes []byte = make([]byte, length)

	for i := 0; i < length; i++ {
		codes[i] = uint8(48 + r.Intn(10))
	}

	return string(codes)
}

func PrettyPrint(pipeline interface{}) {
	json, _ := json.MarshalIndent(pipeline, "", "  ")

	fmt.Println(string(json))
}

// Decode Phone Veridication ID
func DecodeVerificationID(verification_query_str string) (*primitive.ObjectID, *interfaces.PhoneType, error) {

	decodedStr, err := base64.StdEncoding.DecodeString(verification_query_str)
	if err != nil {

		return nil, nil, err
	}

	parsedQuery, err := url.ParseQuery(string(decodedStr))
	if err != nil {
		return nil, nil, err
	}

	objectId, err := primitive.ObjectIDFromHex(parsedQuery.Get("_id"))
	if err != nil {
		return nil, nil, err
	}

	phone := interfaces.PhoneType{
		Code:   parsedQuery.Get("phone_code"),
		Number: parsedQuery.Get("phone_number"),
	}

	return &objectId, &phone, nil
}

func GetTimeZone(ctx *gin.Context) string {
	return ctx.Request.Header.Get(constants.TimeZoneHeaderKey)
}

// Sets Authentication Cookie
func SetAuthCookie(ctx *gin.Context, access, refresh string, accessTime int, domain string, env string) {

	var production = env == "production"

	if production {
		ctx.SetSameSite(http.SameSiteNoneMode)
	} else {
		ctx.SetSameSite(http.SameSiteLaxMode)
	}

	// Set Access Token
	ctx.SetCookie(string(constants.AUTH_TOKEN),
		access,
		accessTime, "/", domain, production, true)

	// Set Refresh Token
	ctx.SetCookie(string(constants.REFRESH_TOKEN),
		refresh,
		constants.CookieRefreshExpiryTime, "/", domain, production, true)
}
