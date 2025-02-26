package main

import (
	"golangapi/config"
	"golangapi/controllers"
	"golangapi/middleware"
	"golangapi/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {
	err := 	godotenv.Load()

	if err != nil {
		log.Fatal("ENV gagal terbaca")
	}

	// gin init
	r := gin.Default()

	// database connect
	db := config.ConnectDatabase()

	// auto migrate
	db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Post{}, &models.Group{})

	// inisialisasi controllers
	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)
	profileController := controllers.NewRelationController(db)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users", userController.GetUsers)
			protected.POST("/profiles", profileController.CreateProfile)
			protected.POST("/posts", profileController.CreatePost)
			protected.POST("/groups", profileController.CreateGroups)
			protected.POST("/groups/add", profileController.AddUserToGroup)
		}
	}

	r.Run(":8080")
}