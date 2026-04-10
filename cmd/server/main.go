package main

import (
	"context"
	"log"
	"net/http"
	"os"

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


	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
