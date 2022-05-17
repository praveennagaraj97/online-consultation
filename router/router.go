package router

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/online-consultation/app"
	"github.com/praveennagaraj97/online-consultation/middlewares"
	"github.com/praveennagaraj97/online-consultation/repository"
)

type Router struct {
	engine      *gin.Engine
	app         *app.ApplicationConfig
	repos       *repository.Repository
	middlewares *middlewares.Middlewares
}

func (router *Router) ListenAndServe(conf *app.ApplicationConfig) {

	router.app = conf
	router.middlewares = &middlewares.Middlewares{}
	router.initializeRepositories()
	r := gin.New()

	r.SetTrustedProxies(nil)
	router.engine = r

	// Initialize Routes
	r.Use(router.middlewares.CORSMiddleware())

	// User and Auth Routes
	router.userRoutes()
	router.consultationRoutes()

	// 404
	r.Use(func(ctx *gin.Context) {
		ctx.JSON(404, map[string]interface{}{
			"message": "Route not available",
		})
	})

	if conf.Environment == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Start the server
	if err := r.Run(fmt.Sprintf(":%s", conf.Port)); err != nil {
		log.Fatalln("Failed to start server")
	}

}

func (r *Router) initializeRepositories() {
	repos := &repository.Repository{}
	repos.Initialize(r.app.MongoClient, r.app.DB.MONGO_DBNAME)

	r.repos = repos
}
