package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() *gorm.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Fail to .env file")
	}
	// Get database credentials from environment variables
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")

	// Construct the DSN using environment variables
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=UTC"

	// Open a connection to the MySQL database
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Migrate the schema
	err = Db.AutoMigrate(&User{}, &Location{}, &Schedule{}, &FriendRequest{}, &Invitation{})
	if err != nil {
		panic("failed to migrate schema: " + err.Error())
	}

	return Db
}
