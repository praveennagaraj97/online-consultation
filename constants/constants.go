package constants

const VerifyCodeAcceptedAttempts uint8 = 3

type CookieNames string

const (
	AUTH_TOKEN              CookieNames = "AUTH_TOKEN"
	REFRESH_TOKEN           CookieNames = "REFRESH_TOKEN"
	CUSTOME_HEADER_LANG_KEY CookieNames = "LOCALE"
)

type UserType string

const (
	SUPER_ADMIN UserType = "super_admin"
	ADMIN       UserType = "admin"
	EDITOR      UserType = "editor"
	USER        UserType = "user"
	DOCTOR      UserType = "doctor"
)

// Cookie - 30 min
const CookieAccessExpiryTime int = 60 * 30

// Cookie - 1 month
const CookieRefreshExpiryTime int = 2.628e+6

// Token - 30 min
const JWT_AccessTokenExpiry = 30

// Pagination options - PerPage default
const DefaultPerPageResults = 10

type PaymentFor string

const (
	ScheduledAppointment PaymentFor = "scheduled_appointment"
)

const TimeZoneHeaderKey string = "Time-Zone"
