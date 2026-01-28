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
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(
		"logs/app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		defer logFile.Close()
		panic(err)
	}

	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
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
	if env == "" {
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
	subCategory := repository.NewSubcategoryRepository(db)

	userService := services.NewUserService(userRepo)
	cartService := services.NewCartService(cartRepo, userRepo, medicRepo, logger)
	orderService := services.NewOrderService(orderRepo, userRepo, cartRepo, medicRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	subCategoryService := services.NewSubcategoryService(subCategory, categoryRepo)

	router := gin.Default()

	transport.RegisterRoutes(router, userService, cartService, orderService, categoryService, subCategoryService, logger)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
