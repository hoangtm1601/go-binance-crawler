package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(server *gin.Engine, db *gorm.DB) {
	router := server.Group("/api")
	_ = router
}
