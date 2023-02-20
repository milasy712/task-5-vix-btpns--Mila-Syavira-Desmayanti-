package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"task-5-vix-fullstack/controllers"
	"task-5-vix-fullstack/middlewares"
)

// Set up routes end point
func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// User routes
	r.POST("/users/login", controllers.Login)
	r.POST("/users/register", controllers.CreateUser)
	r.PUT("/users/:userId", controllers.UpdateUser)
	r.DELETE("/users/:userId", controllers.DeleteUser)

	// PhotoUrl routes
	r.GET("/photos", controllers.GetPhoto)

	// Photo Url yang membutuhkan token jwt
	authorized := r.Group("/").Use(middlewares.Auth())
	{
		authorized.POST("/photos", controllers.CreatePhoto)
		authorized.PUT("/photos/:photoId", controllers.UpdatePhoto)
		authorized.DELETE("/photos/:photoId", controllers.DeletePhoto)
	}

	return r
}


