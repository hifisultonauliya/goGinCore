package helper

import (
	"fmt"
	"log"
	"os"

	"github.com/hifisultonauliya/goGinCore/src/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectDB() (*gorm.DB, error) {
	// Construct the connection string
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open a connection to the database
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Check if the connection is successful
	err = db.DB().Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Successfully connected to the database!")

	// Now you can use db for your database operations with GORM

	db.AutoMigrate(&models.User{}, &models.Task{})

	// Add foreign key constraint
	db.Model(&models.Task{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	return db, err
}
