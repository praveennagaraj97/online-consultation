package serialize

import (
	adminmodel "github.com/praveennagaraj97/online-consultation/models/admin"
	usermodel "github.com/praveennagaraj97/online-consultation/models/user"
)

type InvalidVerificationCodeErrorResponse struct {
	RemainingAttempts uint8 `json:"remaining_attempts"`
	Response
}

type AuthResponse struct {
	DataResponse[*usermodel.UserEntity]
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Response
}

type AdminAuthResponse struct {
	DataResponse[*adminmodel.AdminEntity]
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
