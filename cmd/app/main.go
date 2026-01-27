package main

import (
	"log"
	"log/slog"
	"os"
	"team-pharmacy/internal/config"
	"team-pharmacy/internal/models"
	"team-pharmacy/internal/repository"
	"team-pharmacy/internal/services"
	"team-pharmacy/internal/transport"

	"github.com/gin-gonic/gin"
)

func setupLogger() *slog.Logger {
	level := slog.LevelInfo

	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		switch lvl {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warm":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func main() {

	logger := setupLogger()
	addr := ":8080"
	env := os.Getenv("APP_ENV")
	if env != "" {
		env = "local"
	}

	logger.Info("server started",
		slog.String("addr", addr),
		slog.String("env", env),
	)

	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(&models.User{}, &models.Cart{}, &models.Medicine{}, &models.CartItem{}, &models.Order{}, &models.OrderItem{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}
	userRepo := repository.NewUserRepository(db)
	cartRepo := repository.NewCartRepository(db, logger)
	medicRepo := repository.NewMedicineRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, userRepo, medicRepo, logger)
	orderServer := services.NewOrderService(orderRepo, userRepo, cartRepo, medicRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	router := gin.Default()

	transport.RegisterRoutes(router, userService, cartService, orderServer, categoryService, logger)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
