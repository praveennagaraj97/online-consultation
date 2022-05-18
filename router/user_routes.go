package router

import userapi "github.com/praveennagaraj97/online-consultation/api/user"

func (r *Router) userRoutes() {

	// Initialize User API
	userAPI := userapi.UserAPI{}
	userAPI.Initialize(r.app, r.repos.GetUserRepository(), r.repos.GetOneTimePasswordRepository())

	authRoutes := r.engine.Group("/api/v1/auth")

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

}
