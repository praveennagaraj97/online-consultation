package serialize

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type ErrorResponse struct {
	Errors *map[string]string `json:"errors,omitempty"`
	Response
}

type DataResponse[T interface{}] struct {
	Data T `json:"result"`
	Response
}

type PaginatedDataResponse[T interface{}] struct {
	Count            *uint64 `json:"count"`
	Next             *bool   `json:"next"`
	Prev             *bool   `json:"prev"`
	PaginateKeySetID *string `json:"paginate_id"`
	DataResponse[T]
}

type InvalidVerificationCodeErrorResponse struct {
	RemainingAttempts uint8 `json:"remaining_attempts"`
	Response
}

type AuthResponse[T interface{}] struct {
	DataResponse[T]
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Response
}
