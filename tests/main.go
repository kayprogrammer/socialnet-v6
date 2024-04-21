package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/socialnet-v6/config"
	"github.com/kayprogrammer/socialnet-v6/models"
	"github.com/kayprogrammer/socialnet-v6/routes"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateTables(db *gorm.DB) {
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

		// chat
		&models.Chat{},
		&models.Message{},
	)
	db.Exec("CREATE UNIQUE INDEX unique_requester_requestee ON friends(LEAST(requester_id, requestee_id), GREATEST(requester_id, requestee_id))")
}

func DropTables(db *gorm.DB) {
	// Drop Tables
	db.Migrator().DropTable(
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

		// chat
		&models.Chat{},
		&models.Message{},
	)
}

func CreateSingleTable(db *gorm.DB, model interface{}) {
	db.AutoMigrate(&model)
}

func DropAndCreateSingleTable(db *gorm.DB, model interface{}) {
	db.Migrator().DropTable(&model)
	db.AutoMigrate(&model)
}

func waitForDBConnection(t *testing.T, dsn string) *gorm.DB {
	maxRetries := 3 // Number of retries to wait for the database to be ready
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		t.Logf("Waiting for the database to be ready... Attempt %d", i+1)
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// for i := 0; i < maxRetries; i++ {
	// 	t.Logf("Waiting for the database to be ready... Attempt %d", i+1)
	// 	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 		Logger: logger.Default.LogMode(logger.Silent),
	// 	})
	// 	if err == nil {
	// 		break
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }

	if err != nil {
		t.Fatalf("Failed to connect to the test database: %v", err)
	}

	// Add UUID extension
	result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if result.Error != nil {
		log.Fatal("failed to create extension: " + result.Error.Error())
	}
	return db
}

func SetupTestDatabase(t *testing.T) *gorm.DB {
	cfg := config.GetConfig(true)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.PostgresServer,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.TestPostgresDB,
		cfg.PostgresPort,
		"disable",
		"UTC",
	)
	return waitForDBConnection(t, dsn)
}

func CloseTestDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection: " + err.Error())
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal("Failed to close database connection: " + err.Error())
	}
}

func Setup(t *testing.T, app *fiber.App) *gorm.DB {
	os.Setenv("ENVIRONMENT", "TESTING")

	// Set up the test database
	db := SetupTestDatabase(t)

	routes.SetupRoutes(app, db)
	DropTables(db)
	CreateTables(db)
	return db
}

func ParseResponseBody(t *testing.T, b io.ReadCloser) interface{} {
	body, _ := io.ReadAll(b)
	// Parse the response body as JSON
	responseBody := make(map[string]interface{})
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Errorf("error parsing response body as JSON: %s", err)
	}
	return responseBody
}

func ProcessTestBody(t *testing.T, app *fiber.App, url string, method string, body interface{}, access ...string) *http.Response {
	// Marshal the test data to JSON
	requestBytes, err := json.Marshal(body)
	requestBody := bytes.NewReader(requestBytes)
	assert.Nil(t, err)
	req := httptest.NewRequest(method, url, requestBody)
	req.Header.Set("Content-Type", "application/json")
	if access != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access[0]))
	}
	res, err := app.Test(req)
	if err != nil {
		log.Println(err)
	}
	return res
}
