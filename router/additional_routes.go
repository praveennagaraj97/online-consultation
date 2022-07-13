package router

import (
	additionalapi "github.com/praveennagaraj97/online-consultation/api/additional"
)

func (r *Router) additionalRoutes() {

	api := additionalapi.AdditionalAPI{}
	api.Initailize(r.app)

	router := r.engine.Group("/api/v1/additional")

	router.GET("/jwt/status", api.CheckJWTStatus())
	router.StaticFile("/countries.json", "./static/countries.json")
	router.StaticFile("/states.json", "./static/states.json")
}
