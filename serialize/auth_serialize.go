package serialize

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
