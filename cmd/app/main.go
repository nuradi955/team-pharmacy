package main

import (
	"log"
	"team-pharmacy/internal/config"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"
	"team-pharmacy/internal/services"
	"team-pharmacy/internal/transport"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}
	userRepo := repository.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	router := gin.Default()

	transport.RegisterRoutes(router, userService)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
