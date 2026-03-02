package infrastructure

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	// Data Source Name (DSN) format:
	// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	// It's recommended to use environment variables for credentials in production.
	dsn := "host=localhost user=root password=1212312121 dbname=hex port=5432 sslmode=disable timezone=Asia/Bangkok"

	// Open the connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
		return nil, err
	}

	// Optional: Auto-migrate your schema. Define your structs (models) first.
	// db.AutoMigrate(&Product{}) // Example for a Product model.

	log.Println("Connected to database successfully!")
	return db, nil
}
