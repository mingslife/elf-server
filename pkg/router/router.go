package router

import (
	"time"

	"github.com/gin-gonic/gin"

	"elf-server/pkg/conf"
	"elf-server/pkg/controllers"
	"elf-server/pkg/middleware"
	"elf-server/pkg/render"
)

func NewRouter(cfg *conf.Config) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Static("/elf", "./web")
	router.Static("/assets", "./theme/assets")
	router.Static("/upload", "./upload")
	router.HTMLRender = render.New(&render.RenderOptions{
		TemplateDir:     "./theme/templates",
		ContentType:     "text/html; charset=utf-8",
		DevelopmentMode: cfg.Debug,
	})
	controllers.NewPortalController(router)

	apiRouter := router.Group("/api/v1")
	apiRouter.Use(middleware.NewCorsMiddleware())
	apiRouter.Use(middleware.NewJwtMiddleware(&middleware.JwtMiddlewareConfig{
		Realm:      cfg.Name,
		Key:        []byte(cfg.JwtKey),
		Timeout:    time.Duration(cfg.JwtTimeout) * time.Hour,
		MaxRefresh: time.Duration(cfg.JwtMaxRefresh) * time.Hour,
		ExcludePaths: []string{
			"/api/v1/auth/login",
			"/api/v1/auth/register",
			"/api/v1/auth/settings",
			// "/api/v1/upload/file",
			// "/api/v1/upload/image",
		},
	}))
	controllers.NewAuthController(apiRouter)
	controllers.NewUploadController(apiRouter)
	controllers.NewSettingController(apiRouter)
	controllers.NewNavigationController(apiRouter)
	controllers.NewUserController(apiRouter)
	controllers.NewPostController(apiRouter)
	controllers.NewCategoryController(apiRouter)
	controllers.NewCommentController(apiRouter)

	return router
}
