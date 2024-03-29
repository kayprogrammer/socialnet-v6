package database

import (
	"fmt"
	"log"
	"os"

	"github.com/kayprogrammer/socialnet-v6/config"
	"github.com/kayprogrammer/socialnet-v6/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDb(cfg config.Config) *gorm.DB {
	dsnTemplate := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
	dsn := fmt.Sprintf(
		dsnTemplate,
		cfg.PostgresServer,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
		"disable",
		"UTC",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(2)
	}
	log.Println("Connected to the database successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	// Add UUID extension
	result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if result.Error != nil {
		log.Fatal("failed to create extension: " + result.Error.Error())
	}

	// Add Migrations
	db.AutoMigrate(
		// general
		&models.SiteDetail{},

		// accounts
		&models.Country{},
		&models.Region{},
		&models.City{},
		&models.User{},
		&models.Otp{},

		// feed
		&models.Post{},
		&models.Comment{},
		&models.Reply{},
		&models.Reaction{},

		// profiles
		&models.Friend{},
		&models.Notification{},
	)
	db.Exec("CREATE UNIQUE INDEX unique_requester_requestee ON friends(LEAST(requester_id, requestee_id), GREATEST(requester_id, requestee_id))")
	return db
}
