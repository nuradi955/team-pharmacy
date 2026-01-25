package transport

import (
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine,
	userService services.UserService,
	cartService services.CartService,
	orderService services.OrderService) {

	userHandler := NewUserHandler(userService)
	categoryHandler := NewCategoryHandler(categoryService)
	cartHandler := NewCartHandler(cartService)
	orderHandler := NewOrderHandler(orderService, userService, cartService)

	userHandler.RegisterRoutes(router)
	categoryHandler.RegisterRoutes(router)
	cartHandler.RegisterRoutes(router)
	orderHandler.RegisterRoutes(router)

}
