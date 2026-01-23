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

if err := db.AutoMigrate(
	&models.User{},
	&models.Category{},
	&models.Cart{},
	&models.Medicine{},
	&models.CartItem{},
); err != nil {
	log.Fatalf("не удалось выполнить миграции: %v", err)
}

userRepo := repository.NewUserRepository(db)
categoryRepo := repository.NewCategoryRepository(db)
cartRepo := repository.NewCartRepository(db)
medRepo := repository.NewMedicineRepository(db)

userService := services.NewUserService(userRepo)
categoryService := services.NewCategoryService(categoryRepo)
cartService := services.NewCartService(cartRepo, userRepo, medRepo)

router := gin.Default()

transport.RegisterRoutes(router, userService, categoryService, cartService)

if err := router.Run(); err != nil {
	log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
}
}
