package constants

const VerifyCodeAcceptedAttempts uint8 = 3

type CookieNames string

// Browser cookie names
const (
	// User Auth Cookie
	AUTH_TOKEN    CookieNames = "AUTH_TOKEN"
	REFRESH_TOKEN CookieNames = "REFRESH_TOKEN"
	// Admin Auth Cookie
	ADMIN_AUTH_TOKEN    CookieNames = "ADMIN_AUTH_TOKEN"
	ADMIN_REFRESH_TOKEN CookieNames = "ADMIN_REFRESH_TOKEN"
	// Doctor Auth Cookie
	DOCTOR_AUTH_TOKEN    CookieNames = "DOCTOR_AUTH_TOKEN"
	DOCTOR_REFRESH_TOKEN CookieNames = "DOCTOR_REFRESH_TOKEN"

	CUSTOME_HEADER_LANG_KEY CookieNames = "LOCALE"
)

type CookieType string

// User for setting cookie | Each type is used for different browser.
const (
	USER_AUTH_COOKIE   CookieType = "USER_AUTH_COOKIE"
	ADMIN_AUTH_COOKIE  CookieType = "ADMIN_AUTH_COOKIE"
	DOCTOR_AUTH_COOKIE CookieType = "DOCTOR_AUTH_COOKIE"
)

type UserType string

// Type of user, set in JWT.
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
