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

	if err := db.AutoMigrate(&models.User{}, &models.Cart{}, &models.Medicine{}, &models.CartItem{}, &models.Order{}, &models.OrderItem{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}
	userRepo := repository.NewUserRepository(db)
	cartRepo := repository.NewCartRepository(db)
	medicRepo := repository.NewMedicineRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, userRepo, medicRepo)
	orderServer := services.NewOrderService(orderRepo, userRepo, cartRepo, medicRepo)

router := gin.Default()

	transport.RegisterRoutes(router, userService, cartService, orderServer)

if err := router.Run(); err != nil {
	log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
}
}
