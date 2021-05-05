package routes

import (
	"github.com/gin-gonic/gin"
	"go-web-template/controllers"
	"go-web-template/logger"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", controllers.Homepage)
	return r
}
