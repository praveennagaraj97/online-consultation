package constants

const VerifyCodeAcceptedAttempts uint8 = 3

type Constants string

const (
	AUTH_TOKEN              Constants = "AUTH_TOKEN"
	REFRESH_TOKEN           Constants = "REFRESH_TOKEN"
	CUSTOME_HEADER_LANG_KEY Constants = "LOCALE"
)

// Cookie - 30 min
const CookieAccessExpiryTime int = 60 * 30

// Cookie - 1 month
const CookieRefreshExpiryTime int = 2.628e+6

// Token - 30 min
const JWT_AccessTokenExpiry = 30

// Pagination options - PerPage default
const DefaultPerPageResults = 10
