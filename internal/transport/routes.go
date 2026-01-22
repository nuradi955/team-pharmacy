package transport

import (
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, userService services.UserService) {
	userHandler := NewUserHandler(userService)

	userHandler.RegisterRoutes(router)

}


