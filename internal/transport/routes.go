package transport

import (
	"log/slog"
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine,
	userService services.UserService,
	cartService services.CartService,
	orderService services.OrderService,
	categoryService services.CategoryService,
	subcategoryService services.SubcategoryService,
	logger *slog.Logger) {

	userHandler := NewUserHandler(userService)
	categoryHandler := NewCategoryHandler(categoryService)
	subcategoryHandler := NewSubcategoryHandler(subcategoryService)
	cartHandler := NewCartHandler(logger, cartService)
	orderHandler := NewOrderHandler(orderService, userService, cartService)

	userHandler.RegisterRoutes(router)
	categoryHandler.RegisterRoutes(router)
	subcategoryHandler.RegisterRoutes(router)
	cartHandler.RegisterRoutes(router)
	orderHandler.RegisterRoutes(router)

}
