package router

import userapi "github.com/praveennagaraj97/online-consultation/api/user"

func (r *Router) userRoutes() {

	// Initialize User API
	userAPI := userapi.UserAPI{}
	userAPI.Initialize(r.app,
		r.repos.GetUserRepository(),
		r.repos.GetOneTimePasswordRepository(),
		r.repos.GetUserRelativeRepository(),
		r.repos.GetUserDeliveryAddressRepository(),
	)

	authRoutes := r.engine.Group("/api/v1/auth")
	userRoutes := r.engine.Group("/api/v1/user")

	// Auth Routes
	authRoutes.POST("/register", userAPI.Register())
	authRoutes.POST("/send_verification_code", userAPI.SendVerificationCode())
	authRoutes.POST("/verify_code/:verification_id", userAPI.VerifyCode())
	authRoutes.POST("/signin_with_phonenumber", userAPI.SignInWithPhoneNumber())
	authRoutes.POST("/signin_with_emaillink", userAPI.SignInWithEmailLink())
	authRoutes.GET("/login_with_token/:token", userAPI.SendLoginCredentialsForEmailLink())
	authRoutes.POST("/request_email_verify", userAPI.RequestEmailVerifyLink())
	authRoutes.GET("/verify_email/:token", userAPI.ConfirmEmail())
	authRoutes.POST("/check_email_taken", userAPI.CheckUserExistsByEmail())
	authRoutes.POST("/check_phone_taken", userAPI.CheckUserExistsByPhoneNumber())

	authRoutes.GET("/refresh_token", r.middlewares.IsAuthorized(), r.middlewares.UserRole([]string{"user"}), userAPI.RefreshToken())

	userRoutes.Use(r.middlewares.IsAuthorized(), r.middlewares.UserRole([]string{"user"}))
	// user Routes
	userRoutes.GET("", userAPI.GetUserDetails())
	userRoutes.PATCH("", userAPI.UpdateUserDetails())
	// Relative
	userRoutes.POST("/relative", userAPI.AddRelative())
	userRoutes.GET("/relative", userAPI.GetListOfRelatives())
	userRoutes.GET("/relative/:id", userAPI.GetRelativeProfileById())
	userRoutes.PATCH("/relative/:id", userAPI.UpdateRelativeProfileById())
	userRoutes.DELETE("/relative/:id", userAPI.DeleteRelativeProfileById())
	// Delivery Address
	userRoutes.POST("/delivery_address", userAPI.AddNewAddress())
	userRoutes.GET("/delivery_address", userAPI.GetAllAddress())
	userRoutes.GET("/delivery_address/:id", userAPI.GetAddressById())
	userRoutes.PATCH("/delivery_address/:id", userAPI.UpdateAddressById())
	userRoutes.DELETE("/delivery_address/:id", userAPI.DeleteAddressById())
	userRoutes.PATCH("/delivery_address/:id/mark_default", userAPI.ToggleDefaultAddress(true))
	userRoutes.PATCH("/delivery_address/:id/unmark_default", userAPI.ToggleDefaultAddress(false))

}
