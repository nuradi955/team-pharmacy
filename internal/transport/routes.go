package transport

import (
	"team-pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, userService services.UserService, categoryService services.CategoryService) {
	userHandler := NewUserHandler(userService)
	categoryHandler := NewCategoryHandler(categoryService)

	userHandler.RegisterRoutes(router)
	categoryHandler.RegisterRoutes(router)

}


