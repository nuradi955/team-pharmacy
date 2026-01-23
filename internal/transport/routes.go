package transport

import (
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, userService services.UserService, cartService services.CartService) {
	userHandler := NewUserHandler(userService)
	cartHandler := NewCartHandler(cartService)

	userHandler.RegisterRoutes(router)
	cartHandler.RegisterRoutes(router)

}


