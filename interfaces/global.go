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
	OriginalSrc string `json:"src" bson:"src" form:"src"`
	BlurDataURL string `json:"blur_data_url" bson:"blur_data_url" form:"blur_data_url"`
	Width       string `json:"width" form:"width" bson:"width"`
	Height      string `json:"height" form:"height" bson:"height"`
	KeyID       string `json:"-" bson:"key_id"`
}
