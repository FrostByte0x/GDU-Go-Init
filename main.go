package main

import (
	"log"
	"log/slog"
	"partage-projets/config"
	"partage-projets/models"
	"partage-projets/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "partage-projets/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Share projects over http
// @version 1.0
// @description Share files using a web backend and a database
// @securityDefinitions.ApiKey BearerAuth
// @in header
// @name Authorization
func main() {
	//Create the default router
	router := gin.Default()

	router.SetTrustedProxies(nil)
	router.Use(config.SecurityMiddleware())
	router.Use(config.CORSMiddleware())
	router.Use(config.RateLimit(100))

	err := godotenv.Load()
	if err != nil {
		slog.Warn(err.Error())
	}
	// Register the routes
	routes.ProjectsRoutes(router)
	routes.UserRoutes(router)
	routes.CommentRoutes(router)

	// publish swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Connect to the Database server
	slog.Info("Server starting. Connecting to database..")
	err = config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Database connection succesful!")
	// Auto migrate will create tables and columns as needed by the models.
	config.DB.AutoMigrate(&models.Project{})
	config.DB.AutoMigrate(&models.User{})
	config.DB.AutoMigrate(&models.Comment{})
	slog.Info("Starting web server")
	router.Run(":8080")
}
