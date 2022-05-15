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

}
