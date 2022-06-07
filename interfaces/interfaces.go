package interfaces

type PhoneType struct {
	Code   string `json:"code" bson:"code" form:"code"`
	Number string `json:"number" bson:"number" form:"number"`
}

type SMSType struct {
	Message string
	To      string
}

type ImageType struct {
	OriginalImagePath string `json:"-" bson:"original_image_path"`
	BlurImagePath     string `json:"-" bson:"blur_image_path"`
	OriginalSrc       string `json:"image_src" bson:"image_src,omitempty"`
	BlurDataURL       string `json:"blur_data_url" bson:"blur_data_url,omitempty"`
	Width             uint64 `json:"width" form:"width" bson:"width"`
	Height            uint64 `json:"height" form:"height" bson:"height"`
}

type LocationCoordinatesType struct {
	Longitude int `json:"longitude" bson:"x" form:"longitude"`
	Latitute  int `json:"latitude" bson:"y" form:"latitude"`
}

type MongoPointLocationType struct {
	// Point
	Type string `json:"type" bson:"type"`
	// [longitude, latitude]
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}
