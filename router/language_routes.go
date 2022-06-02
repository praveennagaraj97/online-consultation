package router

import (
	languageapi "github.com/praveennagaraj97/online-consultation/api/language"
	"github.com/praveennagaraj97/online-consultation/constants"
)

func (r *Router) languageRoutes() {

	api := languageapi.LanguageAPI{}
	api.Initialize(r.repos.GetLanguageRepository(), r.app)

	adminRoutes := r.engine.Group("/api/v1/admin/language")
	routes := r.engine.Group("/api/v1/language")

	adminRoutes.Use(r.middlewares.IsAuthorized())

	adminRoutes.POST("", r.middlewares.UserRole([]constants.UserType{constants.ADMIN, constants.EDITOR, constants.SUPER_ADMIN}), api.AddNewLanguage())
	adminRoutes.DELETE("/:id", r.middlewares.UserRole([]constants.UserType{constants.ADMIN, constants.EDITOR, constants.SUPER_ADMIN}), api.DeleteLanguageById())

	routes.GET("", api.GetAllLanguages())
	routes.GET("/:id", api.GetLanguageById())

}
