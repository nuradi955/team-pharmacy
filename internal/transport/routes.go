package transport

import (
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	userService services.UserService,
	categoryService services.CategoryService,
	cartService services.CartService,
) {
	userHandler := NewUserHandler(userService)
	categoryHandler := NewCategoryHandler(categoryService)
	cartHandler := NewCartHandler(cartService)

	userHandler.RegisterRoutes(router)
	categoryHandler.RegisterRoutes(router)
	cartHandler.RegisterRoutes(router)
}
