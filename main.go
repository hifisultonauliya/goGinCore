package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hifisultonauliya/goGinCore/src/controllers"
	"github.com/hifisultonauliya/goGinCore/src/helper"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db *gorm.DB
)

func main() {
	var err error
	r := gin.Default()

	db, err = helper.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
		return
	}
	defer db.Close()

	// Middleware to inject the database connection into the Gin context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	userController := controllers.NewUserController()
	taskController := controllers.NewTaskController()
	authController := controllers.NewAuthController()

	// Routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)

			auth.GET("/google/login", authController.GoogleLogin)
			auth.GET("/google/callback", authController.GoogleCallback)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.POST("/", userController.CreateUser)
			users.GET("/", userController.GetUsers)
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}

		// Task routes
		task := v1.Group("/tasks").Use(AuthMiddleware())
		{
			task.POST("/", taskController.CreateTask)
			task.GET("/", taskController.GetTasks)
			task.GET("/:id", taskController.GetTask)
			task.PUT("/:id", taskController.UpdateTask)
			task.DELETE("/:id", taskController.DeleteTask)

		}
	}

	r.Run(":8080")

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		claims, err := helper.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		log.Println("User ID:", claims.UserID)

		c.Next()
	}
}
