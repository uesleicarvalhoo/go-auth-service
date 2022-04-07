package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/config"
	server "github.com/uesleicarvalhoo/go-auth-service/internal/infra/delivery/http"
	"github.com/uesleicarvalhoo/go-auth-service/internal/infra/repository"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
	"github.com/uesleicarvalhoo/go-auth-service/internal/services/auth"
	"github.com/uesleicarvalhoo/go-auth-service/internal/services/user"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/broker"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/cache"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/database"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/logger"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/trace"
)

const (
	eventChannelBuffer      = 100
	gracefulShutdownTimeout = time.Second * 30
)

func runServer(env config.AppSettings, authService *auth.Service, userService *user.Service) {
	srv := server.NewServer(env, authService, userService)

	// Run server
	go func() {
		logger.Infof("Server running at port :%d", env.ServerPort)

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Info("Listen:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Info("Api server forced to shutdown:", err)
	}
}

func main() {
	// Initialize
	if err := logger.InitLogger(logger.Config{}); err != nil {
		panic(err)
	}

	ctx := context.Background()
	env := config.LoadAppSettingsFromEnv()
	eventChannel := make(chan schemas.Event, eventChannelBuffer)

	cacheClient, err := cache.NewRedisClient(env.CacheConfig)
	if err != nil {
		logger.Fatal("Error on connect to redis:", err)
	}

	// Database
	db, err := database.NewPostgreSQLConnection(env.DatabaseConfig)
	if err != nil {
		logger.Fatal("Error on database connection:", err)
	}

	err = repository.DBMigrate(db, env.DatabaseConfig.Database)
	if err != nil {
		logger.Fatal("Error on run migrations, ", err)
	}

	// Broker
	eventBroker, err := broker.NewRabbitMqClient(env.BrokerConfig)
	if err != nil {
		logger.Fatal("Error on connect to broker, ", err)
	}

	// Tracer
	provider, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: fmt.Sprintf("%s/api/traces", env.TraceURL),
		ServiceName:    env.TraceServiceName,
		ServiceVersion: config.ServiceVersion,
		Environment:    env.Env,
		Disabled:       false,
	})
	if err != nil {
		logger.Fatal(err)
	}

	// Services
	go eventBroker.Start(eventChannel)
	defer eventBroker.End()
	defer provider.Close(ctx)

	userRepository := repository.NewUserRepository(db)
	userService := user.NewService(userRepository, eventChannel)
	authService := auth.NewService(userService, cacheClient, env.SecretKey, eventChannel)

	// Server
	runServer(env, authService, userService)
}
