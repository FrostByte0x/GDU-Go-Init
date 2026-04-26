package routes

import (
	"partage-projets/controllers"
	"partage-projets/middlewares"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.Engine) {
	routesGroup := router.Group("/comments")
	// add the auth middleware
	routesGroup.Use(middlewares.Authentication())

	{
		// If we use a trailing slash here, the 307 redirect will make the client drop the authorization header.
		routesGroup.POST("", controllers.PostComment)
	}
}
