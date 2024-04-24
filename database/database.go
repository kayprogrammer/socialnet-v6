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

func Models() []interface{} {
	return []interface{}{
		// general
		&models.File{},
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

		// chat
		&models.Chat{},
		&models.Message{},
	}
}

func MakeMigrations(db *gorm.DB) {
	models := Models()
	for _, model := range models {
        db.AutoMigrate(model)
    }
	db.Exec("CREATE UNIQUE INDEX unique_requester_requestee ON friends(LEAST(requester_id, requestee_id), GREATEST(requester_id, requestee_id))")
}

func CreateTables(db *gorm.DB) {
	models := Models()
	for _, model := range models {
        db.Migrator().CreateTable(model)
    }
}

func DropTables(db *gorm.DB) {
	// Drop Tables
	models := Models()
	for _, model := range models {
        db.Migrator().DropTable(model)
    }
}

func ConnectDb(cfg config.Config, logs...bool) *gorm.DB {
	dsnTemplate := "host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s"
	dsn := fmt.Sprintf(
		dsnTemplate,
		cfg.PostgresServer,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
		"UTC",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(2)
	}
	log.Println("Connected to the database successfully")

	if len(logs) == 0 { 
		// When extra parameter is passed, don't do the following (from sockets)
		log.Println("Running Migrations")

		// Add UUID extension
		result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
		if result.Error != nil {
			log.Fatal("failed to create extension: " + result.Error.Error())
		}

		// Add Migrations
		MakeMigrations(db)
	}
	return db
}
