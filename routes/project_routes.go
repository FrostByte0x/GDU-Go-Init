package routes

import (
	"partage-projets/controllers"

	"github.com/gin-gonic/gin"
)

func ProjectsRoutes(router *gin.Engine) {
	routesGroup := router.Group("/projects")

	{
		routesGroup.GET("/", controllers.GetProjects)
		routesGroup.GET("/:id", controllers.GetProject)
		routesGroup.POST("/", controllers.PostProject)
		routesGroup.PUT("/:id", controllers.PutProject)
		routesGroup.DELETE("/:id", controllers.DeleteProject)
	}
}
