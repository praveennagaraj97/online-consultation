package interfaces

type PhoneType struct {
	Code   string `json:"code" bson:"code" form:"code"`
	Number string `json:"number" bson:"number" form:"number"`
}

type SMSType struct {
	Message string
	To      string
}
