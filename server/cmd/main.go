package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chrisS41/gobike-server/internal/config"
	"github.com/chrisS41/gobike-server/internal/database"
	"github.com/chrisS41/gobike-server/internal/errors"
	"github.com/chrisS41/gobike-server/internal/handlers"
	"github.com/chrisS41/gobike-server/internal/logger"
	"github.com/chrisS41/gobike-server/internal/models"
	"github.com/chrisS41/gobike-server/internal/version"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 설정 로드
	cfg := config.Load()

	// 로거 초기화
	log := logger.GetInstance(cfg.LogDir, logger.LevelType(cfg.LogLevel))
	log.Info("====================================")
	log.Info("Server initialization started")
	log.Info("====================================")
	log.Trace("Configuration loaded - Log level: %s, Directory: %s", cfg.LogLevel, cfg.LogDir)

	// 버전 정보 출력
	log.Info("Version: %s", version.Version)
	log.Info("Git Revision: %s", version.Revision)
	log.Info("Build Date: %s", version.BuildDate)
	log.Info("Go Version: %s", version.GoVersion)
	log.Info("Compiler: %s", version.Compiler)
	log.Info("Platform: %s", version.Platform)

	// MongoDB 연결
	db, err := database.NewMongoDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()
	log.Info("Database connection established")

	// 핸들러 초기화
	handlers := initializeHandlers(db, log)

	// Gin 설정
	gin.SetMode(cfg.GinMode) //debug, test, release
	router := setupRouter(handlers, log)

	// 그레이스풀 셧다운을 위한 설정
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// 서버 시작을 비동기로 실행
	go func() {
		log.Info("Starting server on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start: %v", err)
		}
	}()

	// 그레이스풀 셧다운 처리
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Info("Server exited properly")
	return nil
}

func initializeHandlers(db *database.MongoDB, log *logger.Log) *handlers.Handlers {
	h := &handlers.Handlers{
		Users:  handlers.NewUserHandler(db.Users, log),
		Routes: handlers.NewRouteHandler(db.Routes, log),
		Rides:  handlers.NewRideHandler(db.Rides, log),
	}
	log.Info("All handlers initialized")
	return h
}

func setupRouter(h *handlers.Handlers, log *logger.Log) *gin.Engine {
	r := gin.New()

	// 미들웨어 설정
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// API 라우트 설정
	api := r.Group("/api")
	{
		setupUserRoutes(api, h.Users)
		setupRouteRoutes(api, h.Routes)
		setupRideRoutes(api, h.Rides)
	}

	// 허용되지 않은 HTTP 메서드 처리
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		// 해당 경로에 등록된 핸들러들의 메서드 목록 가져오기
		allowedMethods := c.Request.Header.Get("Allow")
		c.JSON(
			http.StatusMethodNotAllowed,
			models.NewErrorResponseWithMessage(
				errors.ErrInvalidMethod,
				fmt.Sprintf("메서드 %s는 지원되지 않습니다. 허용된 메서드: %s",
					c.Request.Method,
					allowedMethods,
				),
			),
		)
	})

	// 존재하지 않는 경로 처리
	r.NoRoute(func(c *gin.Context) {
		c.JSON(
			http.StatusNotFound,
			models.NewErrorResponseWithMessage(
				errors.ErrPathNotFound,
				fmt.Sprintf("경로 %s를 찾을 수 없습니다", c.Request.URL.Path),
			),
		)
	})

	log.Info("All routes initialized")
	return r
}

// 라우트 설정 함수들
func setupUserRoutes(api *gin.RouterGroup, h *handlers.UserHandler) {
	users := api.Group("/users")
	{
		// User 관련 엔드포인트
		users.POST("/register", h.Register)
		users.POST("/login", h.Login)
		users.GET("/get/:id", h.GetUser)
		users.PUT("/update/:id", h.UpdateUser)

		// Friend 관련 엔드포인트
		users.POST("/friends/add/:id", h.AddFriend)
		users.GET("/friends/list/:id", h.GetFriends)
	}
}

func setupRouteRoutes(api *gin.RouterGroup, h *handlers.RouteHandler) {
	routes := api.Group("/routes")
	{
		routes.POST("/create", h.CreateRoute)
		routes.GET("/get/:id", h.GetRoute)
		routes.PUT("/update/:id", h.UpdateRoute)
		routes.DELETE("/delete/:id", h.DeleteRoute)
		routes.GET("/list/user/:userId", h.GetUserRoutes)
	}
}

func setupRideRoutes(api *gin.RouterGroup, h *handlers.RideHandler) {
	rides := api.Group("/rides")
	{
		rides.POST("/create", h.CreateRide)
		rides.GET("/get/:id", h.GetRide)
		rides.PUT("/update/:id", h.UpdateRide)
		rides.DELETE("/delete/:id", h.DeleteRide)
		rides.GET("/list/user/:userId", h.GetUserRides)
		rides.GET("/stats/:id", h.GetRideStats)
	}
}
