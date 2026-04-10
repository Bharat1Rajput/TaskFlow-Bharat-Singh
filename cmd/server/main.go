package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/Bharat1Rajput/taskflow-backend/internal/handler"
	"github.com/Bharat1Rajput/taskflow-backend/internal/middleware"
	"github.com/Bharat1Rajput/taskflow-backend/internal/repository"
	"github.com/Bharat1Rajput/taskflow-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("API_PORT")
	dbURL := os.Getenv("DATABASE_URL")

	if port == "" {
		port = "8080"
	}

	// DB
	db, err := repository.NewDB(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	// Run migrations
	if err := repository.RunMigrations(dbURL); err != nil {
		log.Fatal("Migration failed:", err)
	}
	r := gin.Default()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/projects", projectHandler.List)
	auth.POST("/projects", projectHandler.Create)
	auth.PATCH("/projects/:id", projectHandler.Update)
	auth.DELETE("/projects/:id", projectHandler.Delete)

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, projectRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	auth.GET("/projects/:id/tasks", taskHandler.List)
	auth.POST("/projects/:id/tasks", taskHandler.Create)
	auth.PATCH("/tasks/:id", taskHandler.Update)
	auth.DELETE("/tasks/:id", taskHandler.Delete)

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	log.Println("Server running on port:", port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("server error", "err", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")
}
