package http

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/uesleicarvalhoo/go-auth-service/docs" // Load files for generate docs gin-swagger
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/config"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/delivery/http/handler"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/delivery/http/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func initRoutes(engine *gin.Engine, handlers *handler.Handler) {
	authMiddleware := middleware.AuthenticationMiddlware(handlers.AuthSvc)

	engine.GET("/health-check", handlers.HealthCheck)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := engine.Group("/v1/auth")

	auth.POST("/signup", handlers.SignUp)
	auth.POST("/login", handlers.Login)
	auth.POST("/recovery-password", handlers.SendRecoveryPasswordToken)
	auth.POST("/reset-password", handlers.ResetPassword)

	auth.POST("/authorize", handlers.Authorize)
	auth.POST("/refresh-access-token", handlers.RefreshAccessToken)

	user := engine.Group("/v1/user")
	user.Use(authMiddleware)
	user.GET("/me", handlers.GetMe)
	user.POST("/me", handlers.UpdateMe)
	user.DELETE("/me", handlers.DeleteMe)
}

func NewServer(
	cfg config.AppSettings, authService handler.AuthenticationService, userService handler.UserService,
) *http.Server {
	docs.SwaggerInfo.Title = config.ServiceName
	docs.SwaggerInfo.Version = config.ServiceVersion

	// Gin Handler
	customCors := cors.New(cors.Config{
		AllowOrigins: []string{cfg.CorsAllowOrigins},
		AllowMethods: []string{cfg.CorsAllowMethods},
		AllowHeaders: []string{cfg.CorsAllowHeaders},
	})

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
		customCors,
		gin.Recovery(),
		middleware.LogMiddleware(),
		gzip.Gzip(gzip.DefaultCompression),
		otelgin.Middleware(cfg.TraceServiceName),
	)

	handler := handler.NewHandler(authService, userService)

	initRoutes(engine, handler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: engine,
	}
}
